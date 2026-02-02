<template>
  <div class="table-toolbar">
    <div class="toolbar-left">
      <slot name="left"></slot>
    </div>
    <div class="toolbar-right">
      <slot name="right"></slot>
    </div>
  </div>
</template>

<script setup>
// 表格工具栏组件
</script>

<style scoped>
.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  gap: 16px;
}

.toolbar-left,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

@media (max-width: 768px) {
  .table-toolbar {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
