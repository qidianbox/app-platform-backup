<template>
  <div class="config-page">
    <h3>模块管理</h3>
    <p class="desc">管理当前APP启用的功能模块</p>
    <el-card class="config-section">
      <el-table :data="modules" border>
        <el-table-column prop="name" label="模块名称" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-switch v-model="row.enabled" />
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const modules = ref([
  { name: '用户管理', description: '用户注册、登录、信息管理', enabled: true },
  { name: '消息中心', description: '站内消息和通知', enabled: true },
  { name: '推送服务', description: '移动端推送通知', enabled: true },
  { name: '数据统计', description: '埋点和数据分析', enabled: false }
])
</script>
<style lang="scss" scoped>
.config-page { h3 { margin: 0 0 5px; } .desc { color: #999; margin: 0 0 20px; } .config-section { margin-bottom: 20px; } }
</style>
