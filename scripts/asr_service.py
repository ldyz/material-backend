#!/usr/bin/env python3
"""
本地语音转文字服务
使用 FunASR Paraformer 进行语音识别（阿里达摩院开源模型）
中文识别准确率业界领先，支持热词增强
"""

import os
import sys
import json
import tempfile
import logging
import threading
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs
from collections import defaultdict
from pathlib import Path

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# 全局模型
model = None

# 流式识别会话存储
stream_sessions = defaultdict(lambda: {'chunks': [], 'last_text': '', 'lock': threading.Lock()})

# 热词文件路径
HOTWORD_FILE = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'hotwords.txt')


def load_model():
    """加载 FunASR Paraformer 模型"""
    global model
    if model is not None:
        return model

    try:
        from funasr import AutoModel

        # 检查热词文件是否存在
        hotword_path = None
        if os.path.exists(HOTWORD_FILE):
            hotword_path = HOTWORD_FILE
            logger.info(f"使用热词文件: {HOTWORD_FILE}")

        # 初始化 Paraformer 模型
        # paraformer-zh: 中文语音识别模型
        # fsmn-vad: 语音活动检测，自动去除静音
        # ct-punc: 标点恢复模型
        model = AutoModel(
            model="paraformer-zh",           # 中文语音识别模型
            vad_model="fsmn-vad",            # 语音活动检测
            punc_model="ct-punc",            # 标点恢复
            hotword=hotword_path,            # 热词文件
        )
        logger.info("FunASR Paraformer 模型加载成功")
        return model
    except Exception as e:
        logger.error(f"加载模型失败: {e}")
        raise


def transcribe_with_funasr(audio_path):
    """使用 FunASR 进行语音识别"""
    try:
        asr_model = load_model()

        result = asr_model.generate(
            input=audio_path,
            batch_size_s=300,  # 批量处理，单位秒
        )

        if result and len(result) > 0:
            text = result[0].get("text", "")
            return text.strip()
        return ""
    except Exception as e:
        logger.error(f"FunASR 识别失败: {e}")
        raise


class ASRHandler(BaseHTTPRequestHandler):
    """语音识别 HTTP 处理器"""

    def log_message(self, format, *args):
        """自定义日志格式"""
        logger.info("%s - %s", self.client_address[0], format % args)

    def do_GET(self):
        """处理 GET 请求"""
        parsed_path = urlparse(self.path)

        if parsed_path.path == '/health':
            self.send_json_response({'status': 'ok', 'service': 'funasr-paraformer'})
        elif parsed_path.path == '/':
            self.send_json_response({
                'service': 'FunASR Paraformer 本地语音识别服务',
                'version': '3.0.0',
                'model': 'paraformer-zh',
                'features': ['中文高精度识别', 'VAD语音活动检测', '标点自动恢复', '热词增强'],
                'endpoints': {
                    'POST /transcribe': '上传音频文件进行识别',
                    'POST /stream/start': '开始流式识别会话',
                    'POST /stream/chunk': '发送音频片段',
                    'POST /stream/end': '结束流式识别会话',
                    'GET /health': '健康检查'
                }
            })
        else:
            self.send_error(404, 'Not Found')

    def do_POST(self):
        """处理 POST 请求"""
        parsed_path = urlparse(self.path)

        if parsed_path.path == '/transcribe':
            self.handle_transcribe()
        elif parsed_path.path == '/stream/start':
            self.handle_stream_start()
        elif parsed_path.path == '/stream/chunk':
            self.handle_stream_chunk()
        elif parsed_path.path == '/stream/end':
            self.handle_stream_end()
        else:
            self.send_error(404, 'Not Found')

    def handle_transcribe(self):
        """处理普通语音识别请求"""
        try:
            content_type = self.headers.get('Content-Type', '')

            if 'multipart/form-data' in content_type:
                text = self.handle_multipart()
            else:
                content_length = int(self.headers.get('Content-Length', 0))
                audio_data = self.rfile.read(content_length)
                text = self.transcribe_audio_data(audio_data)

            self.send_json_response({
                'success': True,
                'text': text
            })

        except Exception as e:
            logger.error(f"语音识别失败: {e}")
            self.send_json_response({
                'success': False,
                'error': str(e)
            }, status=500)

    def handle_stream_start(self):
        """开始流式识别会话"""
        try:
            content_length = int(self.headers.get('Content-Length', 0))
            body = self.rfile.read(content_length)
            data = json.loads(body) if body else {}

            import uuid
            session_id = data.get('session_id', str(uuid.uuid4()))

            # 初始化会话
            with stream_sessions[session_id]['lock']:
                stream_sessions[session_id]['chunks'] = []
                stream_sessions[session_id]['last_text'] = ''

            logger.info(f"流式识别会话开始: {session_id}")

            self.send_json_response({
                'success': True,
                'session_id': session_id
            })
        except Exception as e:
            logger.error(f"启动流式识别失败: {e}")
            self.send_json_response({
                'success': False,
                'error': str(e)
            }, status=500)

    def handle_stream_chunk(self):
        """处理流式音频片段"""
        try:
            content_type = self.headers.get('Content-Type', '')

            # 获取 session_id
            parsed = urlparse(self.path)
            params = parse_qs(parsed.query)
            session_id = params.get('session_id', [None])[0]

            if not session_id:
                # 从 body 获取
                content_length = int(self.headers.get('Content-Length', 0))
                body = self.rfile.read(content_length)
                self.send_json_response({
                    'success': False,
                    'error': 'session_id is required'
                }, status=400)
                return

            # 读取音频数据
            if 'multipart/form-data' in content_type:
                import cgi
                form = cgi.FieldStorage(
                    fp=self.rfile,
                    headers=self.headers,
                    environ={
                        'REQUEST_METHOD': 'POST',
                        'CONTENT_TYPE': content_type
                    }
                )
                audio_field = None
                for key in form.keys():
                    if key.lower() in ['audio', 'chunk', 'data']:
                        audio_field = form[key]
                        break
                if audio_field:
                    audio_data = audio_field.file.read()
                else:
                    audio_data = form.getvalue('audio', form.getvalue('chunk', b''))
            else:
                content_length = int(self.headers.get('Content-Length', 0))
                audio_data = self.rfile.read(content_length)

            if not audio_data:
                self.send_json_response({
                    'success': False,
                    'error': 'No audio data'
                }, status=400)
                return

            # 累积音频
            with stream_sessions[session_id]['lock']:
                stream_sessions[session_id]['chunks'].append(audio_data)
                chunks = stream_sessions[session_id]['chunks'].copy()

            # 合并所有音频片段进行识别
            combined_audio = b''.join(chunks)

            # 识别当前累积的音频
            text = self.transcribe_audio_data(combined_audio)

            self.send_json_response({
                'success': True,
                'text': text,
                'partial': True,
                'chunks_count': len(chunks)
            })

        except Exception as e:
            logger.error(f"处理音频片段失败: {e}")
            self.send_json_response({
                'success': False,
                'error': str(e)
            }, status=500)

    def handle_stream_end(self):
        """结束流式识别会话"""
        try:
            content_length = int(self.headers.get('Content-Length', 0))
            body = self.rfile.read(content_length)
            data = json.loads(body) if body else {}

            session_id = data.get('session_id')

            if not session_id or session_id not in stream_sessions:
                self.send_json_response({
                    'success': False,
                    'error': 'Invalid session_id'
                }, status=400)
                return

            # 获取最终结果
            with stream_sessions[session_id]['lock']:
                chunks = stream_sessions[session_id]['chunks'].copy()
                # 清理会话
                del stream_sessions[session_id]

            # 最终识别
            if chunks:
                combined_audio = b''.join(chunks)
                final_text = self.transcribe_audio_data(combined_audio)
            else:
                final_text = ''

            logger.info(f"流式识别会话结束: {session_id}, 结果: {final_text}")

            self.send_json_response({
                'success': True,
                'text': final_text,
                'partial': False
            })

        except Exception as e:
            logger.error(f"结束流式识别失败: {e}")
            self.send_json_response({
                'success': False,
                'error': str(e)
            }, status=500)

    def handle_multipart(self):
        """处理 multipart/form-data 格式的音频"""
        import cgi

        content_type = self.headers.get('Content-Type')
        form = cgi.FieldStorage(
            fp=self.rfile,
            headers=self.headers,
            environ={
                'REQUEST_METHOD': 'POST',
                'CONTENT_TYPE': content_type
            }
        )

        audio_field = None
        for key in form.keys():
            if key.lower() in ['audio', 'file', 'voice']:
                audio_field = form[key]
                break

        if audio_field is None or not audio_field.file:
            raise ValueError("未找到音频文件")

        audio_data = audio_field.file.read()
        return self.transcribe_audio_data(audio_data)

    def transcribe_audio_data(self, audio_data):
        """识别音频数据"""
        if not audio_data or len(audio_data) < 1000:
            return ''

        # 保存到临时文件
        # FunASR 支持多种格式：wav, mp3, webm, m4a 等
        suffix = '.webm'  # 默认 webm 格式

        # 尝试检测格式
        if audio_data[:4] == b'RIFF':
            suffix = '.wav'
        elif audio_data[:3] == b'ID3' or audio_data[:2] == b'\xff\xfb':
            suffix = '.mp3'
        elif audio_data[:4] == b'ftyp':
            suffix = '.m4a'

        with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as f:
            f.write(audio_data)
            temp_path = f.name

        try:
            text = transcribe_with_funasr(temp_path)
            return text

        except Exception as e:
            logger.error(f"识别失败: {e}")
            return ''
        finally:
            if os.path.exists(temp_path):
                os.unlink(temp_path)

    def send_json_response(self, data, status=200):
        """发送 JSON 响应"""
        self.send_response(status)
        self.send_header('Content-Type', 'application/json; charset=utf-8')
        self.send_header('Access-Control-Allow-Origin', '*')
        self.end_headers()
        self.wfile.write(json.dumps(data, ensure_ascii=False).encode('utf-8'))


def run_server(port=8089):
    """启动服务器"""
    logger.info("正在加载 FunASR Paraformer 模型，请稍候...")
    logger.info("首次运行会自动下载模型文件（约2GB），请耐心等待...")
    try:
        load_model()
    except Exception as e:
        logger.error(f"模型加载失败: {e}")
        sys.exit(1)

    server_address = ('', port)
    httpd = HTTPServer(server_address, ASRHandler)

    logger.info(f"语音识别服务启动在端口 {port}")
    logger.info("模型: FunASR Paraformer (阿里达摩院)")
    logger.info("特性: 中文高精度识别、VAD、标点恢复、热词增强")
    logger.info("API 端点:")
    logger.info("  POST /transcribe - 上传音频文件进行识别")
    logger.info("  POST /stream/start - 开始流式识别会话")
    logger.info("  POST /stream/chunk - 发送音频片段")
    logger.info("  POST /stream/end - 结束流式识别会话")
    logger.info("  GET /health - 健康检查")

    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        logger.info("服务停止")
        httpd.shutdown()


if __name__ == '__main__':
    port = int(os.environ.get('ASR_PORT', 8089))
    run_server(port)
