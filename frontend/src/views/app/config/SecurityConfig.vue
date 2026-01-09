<template>
  <div class="config-page">
    <h3>安全配置</h3>
    <p class="desc">配置安全策略和访问控制</p>
    <el-form label-position="top">
      <el-card class="config-section">
        <template #header>安全策略</template>
        <el-form-item label="启用HTTPS强制">
          <el-switch v-model="config.forceHttps" />
        </el-form-item>
        <el-form-item label="API签名验证">
          <el-switch v-model="config.apiSign" />
        </el-form-item>
        <el-form-item label="IP白名单">
          <el-input v-model="config.ipWhitelist" type="textarea" rows="3" placeholder="每行一个IP" />
        </el-form-item>
      </el-card>
      <el-button type="primary" @click="handleSave">保存配置</el-button>
    </el-form>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
const config = ref({ forceHttps: true, apiSign: true, ipWhitelist: '' })
const handleSave = () => ElMessage.success('配置已保存')
</script>
<style lang="scss" scoped>
.config-page { h3 { margin: 0 0 5px; } .desc { color: #999; margin: 0 0 20px; } .config-section { margin-bottom: 20px; } }
</style>
