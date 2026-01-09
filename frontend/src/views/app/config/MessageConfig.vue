<template>
  <div class="config-page">
    <h3>消息推送配置</h3>
    <p class="desc">配置消息中心和推送通知相关设置</p>
    
    <el-form label-position="top">
      <el-card class="config-section">
        <template #header>推送配置</template>
        <el-form-item label="启用推送">
          <el-switch v-model="config.enablePush" />
        </el-form-item>
        <el-form-item label="推送平台">
          <el-checkbox-group v-model="config.platforms">
            <el-checkbox label="ios">iOS (APNs)</el-checkbox>
            <el-checkbox label="android">Android (FCM)</el-checkbox>
            <el-checkbox label="huawei">华为推送</el-checkbox>
            <el-checkbox label="xiaomi">小米推送</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-card>
      
      <el-card class="config-section">
        <template #header>消息设置</template>
        <el-form-item label="消息保留天数">
          <el-input-number v-model="config.retentionDays" :min="7" :max="365" />
        </el-form-item>
        <el-form-item label="单用户最大未读数">
          <el-input-number v-model="config.maxUnread" :min="10" :max="1000" />
        </el-form-item>
      </el-card>
      
      <el-button type="primary" @click="handleSave">保存配置</el-button>
    </el-form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const config = ref({
  enablePush: true,
  platforms: ['ios', 'android'],
  retentionDays: 30,
  maxUnread: 100
})

const handleSave = () => {
  ElMessage.success('配置已保存')
}
</script>

<style lang="scss" scoped>
.config-page {
  h3 { margin: 0 0 5px; }
  .desc { color: #999; margin: 0 0 20px; }
  .config-section { margin-bottom: 20px; }
}
</style>
