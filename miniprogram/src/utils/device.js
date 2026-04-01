/**
 * 小程序原生功能封装
 * 替代 Capacitor 插件
 */

/**
 * 获取位置信息
 */
export function getLocation() {
  return new Promise((resolve, reject) => {
    uni.getLocation({
      type: 'gcj02',
      success: (res) => {
        resolve({
          latitude: res.latitude,
          longitude: res.longitude,
          accuracy: res.accuracy,
          address: res.address // 可能需要逆地理编码
        })
      },
      fail: (error) => {
        reject(new Error(error.errMsg || '获取位置失败'))
      }
    })
  })
}

/**
 * 选择图片
 */
export function chooseImage(count = 1, sourceType = ['album', 'camera']) {
  return new Promise((resolve, reject) => {
    uni.chooseImage({
      count,
      sourceType,
      success: (res) => {
        resolve({
          tempFilePaths: res.tempFilePaths,
          tempFiles: res.tempFiles
        })
      },
      fail: (error) => {
        reject(new Error(error.errMsg || '选择图片失败'))
      }
    })
  })
}

/**
 * 拍照
 */
export function takePhoto() {
  return chooseImage(1, ['camera'])
}

/**
 * 从相册选择图片
 */
export function chooseFromAlbum(count = 1) {
  return chooseImage(count, ['album'])
}

/**
 * 上传图片
 */
export function uploadImage(filePath) {
  const token = uni.getStorageSync('token')
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: 'https://home.mbed.org.cn/upload/image',
      filePath: filePath,
      name: 'file',
      header: {
        'Authorization': `Bearer ${token}`
      },
      success: (response) => {
        if (response.statusCode === 200) {
          try {
            const data = JSON.parse(response.data)
            resolve(data)
          } catch (e) {
            reject(new Error('解析响应失败'))
          }
        } else {
          reject(new Error('上传失败'))
        }
      },
      fail: (error) => {
        reject(new Error(error.errMsg || '上传失败'))
      }
    })
  })
}

/**
 * 预览图片
 */
export function previewImage(urls, current = 0) {
  uni.previewImage({
    urls,
    current
  })
}

/**
 * 选择文件（从聊天记录）
 * 用于 Excel 导入功能
 */
export function chooseFile(extension = ['xlsx', 'xls']) {
  return new Promise((resolve, reject) => {
    // #ifdef MP-WEIXIN
    wx.chooseMessageFile({
      count: 1,
      type: 'file',
      extension,
      success: (res) => {
        resolve({
          path: res.tempFiles[0].path,
          name: res.tempFiles[0].name,
          size: res.tempFiles[0].size
        })
      },
      fail: (error) => {
        reject(new Error(error.errMsg || '选择文件失败'))
      }
    })
    // #endif

    // #ifndef MP-WEIXIN
    reject(new Error('当前平台不支持选择文件'))
    // #endif
  })
}

/**
 * 上传文件
 */
export function uploadFile(filePath, url = '/upload/file') {
  const token = uni.getStorageSync('token')
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: 'https://home.mbed.org.cn' + url,
      filePath: filePath,
      name: 'file',
      header: {
        'Authorization': `Bearer ${token}`
      },
      success: (response) => {
        if (response.statusCode === 200) {
          try {
            const data = JSON.parse(response.data)
            resolve(data)
          } catch (e) {
            reject(new Error('解析响应失败'))
          }
        } else {
          reject(new Error('上传失败'))
        }
      },
      fail: (error) => {
        reject(new Error(error.errMsg || '上传失败'))
      }
    })
  })
}

/**
 * 开始录音
 */
export function startRecord() {
  return new Promise((resolve, reject) => {
    // #ifdef MP-WEIXIN
    const recorderManager = wx.getRecorderManager()

    recorderManager.onStart(() => {
      resolve(recorderManager)
    })

    recorderManager.onError((error) => {
      reject(new Error(error.errMsg || '录音失败'))
    })

    recorderManager.start({
      format: 'mp3',
      duration: 60000 // 最长 60 秒
    })
    // #endif

    // #ifndef MP-WEIXIN
    reject(new Error('当前平台不支持录音'))
    // #endif
  })
}

/**
 * 停止录音
 */
export function stopRecord(recorderManager) {
  return new Promise((resolve, reject) => {
    // #ifdef MP-WEIXIN
    recorderManager.onStop((res) => {
      resolve({
        tempFilePath: res.tempFilePath,
        duration: res.duration,
        fileSize: res.fileSize
      })
    })

    recorderManager.stop()
    // #endif

    // #ifndef MP-WEIXIN
    reject(new Error('当前平台不支持录音'))
    // #endif
  })
}

/**
 * 播放音频
 */
export function playVoice(filePath) {
  return new Promise((resolve, reject) => {
    // #ifdef MP-WEIXIN
    const innerAudioContext = wx.createInnerAudioContext()

    innerAudioContext.src = filePath
    innerAudioContext.onEnded(() => {
      resolve()
    })
    innerAudioContext.onError((error) => {
      reject(new Error(error.errMsg || '播放失败'))
    })
    innerAudioContext.play()
    // #endif

    // #ifndef MP-WEIXIN
    reject(new Error('当前平台不支持播放音频'))
    // #endif
  })
}

/**
 * 扫码
 */
export function scanCode() {
  return new Promise((resolve, reject) => {
    uni.scanCode({
      success: (res) => {
        resolve({
          result: res.result,
          scanType: res.scanType
        })
      },
      fail: (error) => {
        reject(new Error(error.errMsg || '扫码失败'))
      }
    })
  })
}

/**
 * 获取系统信息
 */
export function getSystemInfo() {
  return new Promise((resolve) => {
    uni.getSystemInfo({
      success: (res) => {
        resolve(res)
      }
    })
  })
}

/**
 * 拨打电话
 */
export function makePhoneCall(phoneNumber) {
  return new Promise((resolve, reject) => {
    uni.makePhoneCall({
      phoneNumber,
      success: () => resolve(),
      fail: (error) => reject(new Error(error.errMsg || '拨打失败'))
    })
  })
}

/**
 * 设置剪贴板
 */
export function setClipboardData(data) {
  return new Promise((resolve) => {
    uni.setClipboardData({
      data,
      success: () => resolve()
    })
  })
}

/**
 * 显示地图
 */
export function openLocation(latitude, longitude, name = '', address = '') {
  uni.openLocation({
    latitude,
    longitude,
    name,
    address,
    scale: 18
  })
}

export default {
  getLocation,
  chooseImage,
  takePhoto,
  chooseFromAlbum,
  uploadImage,
  previewImage,
  chooseFile,
  uploadFile,
  startRecord,
  stopRecord,
  playVoice,
  scanCode,
  getSystemInfo,
  makePhoneCall,
  setClipboardData,
  openLocation
}
