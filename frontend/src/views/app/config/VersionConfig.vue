<template>
  <div class="config-page">
    <h3>版本管理</h3>
    <p class="desc">管理APP版本和更新策略</p>
    <el-card class="config-section">
      <template #header>版本列表</template>
      <el-table :data="versions" border>
        <el-table-column prop="version" label="版本号" width="120" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="发布时间" />
      </el-table>
    </el-card>
    <el-button type="primary">发布新版本</el-button>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const versions = ref([
  { version: '1.0.0', platform: 'iOS', status: 'active', createdAt: '2026-01-01' },
  { version: '1.0.0', platform: 'Android', status: 'active', createdAt: '2026-01-01' }
])
</script>
<style lang="scss" scoped>
.config-page { h3 { margin: 0 0 5px; } .desc { color: #999; margin: 0 0 20px; } .config-section { margin-bottom: 20px; } }
</style>
