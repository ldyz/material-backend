<script>
export default {
  onLaunch: function () {
    console.log('App Launch')
    // 检查登录状态
    const token = uni.getStorageSync('token')
    if (!token) {
      uni.reLaunch({
        url: '/pages/login/index'
      })
    }
  },
  onShow: function () {
    console.log('App Show')
  },
  onHide: function () {
    console.log('App Hide')
  }
}
</script>

<style>
/* 全局样式 */
page {
  background-color: #f7f8fa;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Helvetica Neue', sans-serif;
  font-size: 14px;
  color: #323233;
}

/* 通用工具类 */
.text-center {
  text-align: center;
}

.text-right {
  text-align: right;
}

.flex {
  display: flex;
}

.flex-center {
  display: flex;
  align-items: center;
  justify-content: center;
}

.flex-between {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.flex-column {
  flex-direction: column;
}

.flex-1 {
  flex: 1;
}

/* 安全区域适配 */
.safe-area-bottom {
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}

/* 颜色变量 */
:root {
  --primary-color: #1989fa;
  --success-color: #07c160;
  --warning-color: #ff976a;
  --danger-color: #ee0a24;
  --text-color: #323233;
  --text-color-light: #969799;
  --border-color: #ebedf0;
  --background-color: #f7f8fa;
}
</style>
