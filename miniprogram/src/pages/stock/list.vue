<template>
  <view class="stock-list-page">
    <!-- 搜索栏 -->
    <view class="search-bar">
      <input
        class="search-input"
        type="text"
        v-model="searchKeyword"
        placeholder="搜索物料名称或规格"
        @confirm="onSearch"
      />
      <view class="search-btn" @click="onSearch">搜索</view>
    </view>

    <!-- 列表 -->
    <ListContainer
      :list="list"
      :loading="loading"
      :has-more="hasMore"
      @refresh="onRefresh"
      @load-more="onLoadMore"
    >
      <ListItemCard
        v-for="item in list"
        :key="item.id"
        :title="item.material_name"
        :items="[
          { label: '规格型号', value: item.specification },
          { label: '单位', value: item.unit },
          { label: '库存数量', value: item.quantity },
          { label: '项目', value: item.project_name }
        ]"
        clickable
        @click="goToDetail(item.id)"
      />
    </ListContainer>
  </view>
</template>

<script>
import { getStockList } from '@/api/stock'
import ListContainer from '@/components/ListContainer.vue'
import ListItemCard from '@/components/ListItemCard.vue'

export default {
  components: {
    ListContainer,
    ListItemCard
  },
  data() {
    return {
      list: [],
      loading: false,
      hasMore: false,
      page: 1,
      pageSize: 20,
      searchKeyword: ''
    }
  },
  onShow() {
    this.onRefresh()
  },
  methods: {
    async fetchData(callback) {
      if (this.loading) return
      this.loading = true

      try {
        const res = await getStockList({
          page: this.page,
          page_size: this.pageSize,
          search: this.searchKeyword || undefined
        })

        const items = res.data?.items || res.data || []

        if (this.page === 1) {
          this.list = items
        } else {
          this.list = [...this.list, ...items]
        }

        this.hasMore = items.length >= this.pageSize
      } catch (error) {
        uni.showToast({
          title: error.message || '加载失败',
          icon: 'none'
        })
      } finally {
        this.loading = false
        if (callback) callback()
      }
    },
    onRefresh(callback) {
      this.page = 1
      this.fetchData(callback)
    },
    onLoadMore() {
      this.page++
      this.fetchData()
    },
    onSearch() {
      this.onRefresh()
    },
    goToDetail(id) {
      uni.navigateTo({ url: `/pages/stock/detail?id=${id}` })
    }
  }
}
</script>

<style scoped>
.stock-list-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.search-bar {
  display: flex;
  padding: 24rpx;
  background-color: #ffffff;
  border-bottom: 1rpx solid #ebedf0;
}

.search-input {
  flex: 1;
  height: 72rpx;
  background-color: #f7f8fa;
  border-radius: 8rpx;
  padding: 0 24rpx;
  font-size: 28rpx;
}

.search-btn {
  margin-left: 16rpx;
  padding: 0 24rpx;
  height: 72rpx;
  line-height: 72rpx;
  background-color: #1989fa;
  color: #ffffff;
  border-radius: 8rpx;
  font-size: 28rpx;
}
</style>
