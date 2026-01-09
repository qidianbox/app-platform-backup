<template>
  <div class="config-page">
    <h3>支付配置</h3>
    <p class="desc">配置支付渠道和相关设置</p>
    <el-form label-position="top">
      <el-card class="config-section">
        <template #header>支付渠道</template>
        <el-form-item label="微信支付">
          <el-switch v-model="config.wechat" />
        </el-form-item>
        <el-form-item label="支付宝">
          <el-switch v-model="config.alipay" />
        </el-form-item>
        <el-form-item label="Apple Pay">
          <el-switch v-model="config.applePay" />
        </el-form-item>
      </el-card>
      <el-button type="primary" @click="handleSave">保存配置</el-button>
    </el-form>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
const config = ref({ wechat: true, alipay: true, applePay: false })
const handleSave = () => ElMessage.success('配置已保存')
</script>
<style lang="scss" scoped>
.config-page { h3 { margin: 0 0 5px; } .desc { color: #999; margin: 0 0 20px; } .config-section { margin-bottom: 20px; } }
</style>
