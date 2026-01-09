<template>
  <div class="config-page">
    <h3>数据统计配置</h3>
    <p class="desc">配置埋点和数据分析相关设置</p>
    <el-form label-position="top">
      <el-card class="config-section">
        <template #header>埋点设置</template>
        <el-form-item label="启用自动埋点">
          <el-switch v-model="config.autoTrack" />
        </el-form-item>
        <el-form-item label="数据保留天数">
          <el-input-number v-model="config.retentionDays" :min="30" :max="365" />
        </el-form-item>
      </el-card>
      <el-button type="primary" @click="handleSave">保存配置</el-button>
    </el-form>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
const config = ref({ autoTrack: true, retentionDays: 90 })
const handleSave = () => ElMessage.success('配置已保存')
</script>
<style lang="scss" scoped>
.config-page { h3 { margin: 0 0 5px; } .desc { color: #999; margin: 0 0 20px; } .config-section { margin-bottom: 20px; } }
</style>
