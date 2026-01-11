<template>
  <div class="module-config-container">
    <div class="config-header">
      <h2>用户中心服务配置</h2>
      <p class="config-description">配置用户注册、登录、第三方登录、实名认证、隐私设置等功能</p>
    </div>

    <el-card v-loading="loading">
      <el-form :model="formData" ref="formRef" label-width="160px">
        
        <!-- 注册配置 -->
        <el-divider content-position="left">
          <el-icon><User /></el-icon>
          <span style="margin-left: 8px;">注册配置</span>
        </el-divider>
        
        <el-form-item label="允许注册">
          <el-switch v-model="formData.allow_register" />
          <span class="form-tip">关闭后新用户无法注册</span>
        </el-form-item>

        <el-form-item label="注册方式">
          <el-checkbox-group v-model="formData.register_methods">
            <el-checkbox label="phone">手机号</el-checkbox>
            <el-checkbox label="email">邮箱</el-checkbox>
            <el-checkbox label="username">用户名</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="需要实名认证">
          <el-switch v-model="formData.require_real_name" />
          <span class="form-tip">开启后用户注册时需要完成实名认证</span>
        </el-form-item>

        <el-form-item label="注册验证码">
          <el-switch v-model="formData.register_captcha" />
          <span class="form-tip">开启后注册时需要验证码</span>
        </el-form-item>

        <!-- 第三方登录配置 -->
        <el-divider content-position="left">
          <el-icon><Connection /></el-icon>
          <span style="margin-left: 8px;">第三方登录配置</span>
        </el-divider>

        <!-- 微信登录 -->
        <el-form-item label="微信登录">
          <el-switch v-model="formData.wechat_login.enabled" />
        </el-form-item>
        <template v-if="formData.wechat_login.enabled">
          <el-form-item label="微信AppID">
            <el-input v-model="formData.wechat_login.app_id" placeholder="请输入微信AppID" />
          </el-form-item>
          <el-form-item label="微信AppSecret">
            <el-input v-model="formData.wechat_login.app_secret" type="password" placeholder="请输入微信AppSecret" show-password />
          </el-form-item>
        </template>

        <!-- QQ登录 -->
        <el-form-item label="QQ登录">
          <el-switch v-model="formData.qq_login.enabled" />
        </el-form-item>
        <template v-if="formData.qq_login.enabled">
          <el-form-item label="QQ AppID">
            <el-input v-model="formData.qq_login.app_id" placeholder="请输入QQ AppID" />
          </el-form-item>
          <el-form-item label="QQ AppKey">
            <el-input v-model="formData.qq_login.app_key" type="password" placeholder="请输入QQ AppKey" show-password />
          </el-form-item>
        </template>

        <!-- 微博登录 -->
        <el-form-item label="微博登录">
          <el-switch v-model="formData.weibo_login.enabled" />
        </el-form-item>
        <template v-if="formData.weibo_login.enabled">
          <el-form-item label="微博AppKey">
            <el-input v-model="formData.weibo_login.app_key" placeholder="请输入微博AppKey" />
          </el-form-item>
          <el-form-item label="微博AppSecret">
            <el-input v-model="formData.weibo_login.app_secret" type="password" placeholder="请输入微博AppSecret" show-password />
          </el-form-item>
        </template>

        <!-- Apple ID登录 -->
        <el-form-item label="Apple ID登录">
          <el-switch v-model="formData.apple_login.enabled" />
        </el-form-item>
        <template v-if="formData.apple_login.enabled">
          <el-form-item label="Service ID">
            <el-input v-model="formData.apple_login.service_id" placeholder="请输入Service ID" />
          </el-form-item>
          <el-form-item label="Team ID">
            <el-input v-model="formData.apple_login.team_id" placeholder="请输入Team ID" />
          </el-form-item>
          <el-form-item label="Key ID">
            <el-input v-model="formData.apple_login.key_id" placeholder="请输入Key ID" />
          </el-form-item>
        </template>

        <!-- 抖音登录 -->
        <el-form-item label="抖音登录">
          <el-switch v-model="formData.douyin_login.enabled" />
        </el-form-item>
        <template v-if="formData.douyin_login.enabled">
          <el-form-item label="抖音ClientKey">
            <el-input v-model="formData.douyin_login.client_key" placeholder="请输入抖音ClientKey" />
          </el-form-item>
          <el-form-item label="抖音ClientSecret">
            <el-input v-model="formData.douyin_login.client_secret" type="password" placeholder="请输入抖音ClientSecret" show-password />
          </el-form-item>
        </template>

        <!-- 登录配置 -->
        <el-divider content-position="left">
          <el-icon><Lock /></el-icon>
          <span style="margin-left: 8px;">登录配置</span>
        </el-divider>

        <el-form-item label="密码最小长度">
          <el-input-number v-model="formData.password_min_length" :min="6" :max="20" />
          <span class="form-tip">建议8位以上</span>
        </el-form-item>

        <el-form-item label="密码复杂度要求">
          <el-checkbox-group v-model="formData.password_complexity">
            <el-checkbox label="number">必须包含数字</el-checkbox>
            <el-checkbox label="letter">必须包含字母</el-checkbox>
            <el-checkbox label="special">必须包含特殊字符</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="登录失败锁定">
          <el-switch v-model="formData.login_lock_enabled" />
          <span class="form-tip">防止暴力破解</span>
        </el-form-item>

        <el-form-item label="失败次数阈值" v-if="formData.login_lock_enabled">
          <el-input-number v-model="formData.login_lock_threshold" :min="3" :max="10" />
          <span class="form-tip">次</span>
        </el-form-item>

        <el-form-item label="锁定时长" v-if="formData.login_lock_enabled">
          <el-input-number v-model="formData.login_lock_duration" :min="5" :max="1440" />
          <span class="form-tip">分钟</span>
        </el-form-item>

        <!-- 用户信息管理 -->
        <el-divider content-position="left">
          <el-icon><UserFilled /></el-icon>
          <span style="margin-left: 8px;">用户信息管理</span>
        </el-divider>

        <el-form-item label="必填字段">
          <el-checkbox-group v-model="formData.required_fields">
            <el-checkbox label="nickname">昵称</el-checkbox>
            <el-checkbox label="avatar">头像</el-checkbox>
            <el-checkbox label="gender">性别</el-checkbox>
            <el-checkbox label="birthday">生日</el-checkbox>
            <el-checkbox label="region">地区</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="允许修改用户名">
          <el-switch v-model="formData.allow_change_username" />
          <span class="form-tip">关闭后用户名不可修改</span>
        </el-form-item>

        <el-form-item label="昵称敏感词过滤">
          <el-switch v-model="formData.nickname_filter" />
        </el-form-item>

        <el-form-item label="头像审核">
          <el-switch v-model="formData.avatar_review" />
          <span class="form-tip">开启后头像需要审核</span>
        </el-form-item>

        <!-- 实名认证配置 -->
        <el-divider content-position="left">
          <el-icon><CreditCard /></el-icon>
          <span style="margin-left: 8px;">实名认证配置</span>
        </el-divider>

        <el-form-item label="启用实名认证">
          <el-switch v-model="formData.real_name_auth.enabled" />
        </el-form-item>

        <template v-if="formData.real_name_auth.enabled">
          <el-form-item label="认证方式">
            <el-checkbox-group v-model="formData.real_name_auth.methods">
              <el-checkbox label="idcard">身份证认证</el-checkbox>
              <el-checkbox label="face">人脸识别</el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <el-form-item label="认证服务商">
            <el-select v-model="formData.real_name_auth.provider" placeholder="请选择">
              <el-option label="阿里云" value="aliyun" />
              <el-option label="腾讯云" value="tencent" />
              <el-option label="百度云" value="baidu" />
            </el-select>
          </el-form-item>

          <el-form-item label="服务商AppID">
            <el-input v-model="formData.real_name_auth.app_id" placeholder="请输入服务商AppID" />
          </el-form-item>

          <el-form-item label="服务商AppSecret">
            <el-input v-model="formData.real_name_auth.app_secret" type="password" placeholder="请输入服务商AppSecret" show-password />
          </el-form-item>

          <el-form-item label="强制实名认证">
            <el-switch v-model="formData.real_name_auth.required" />
            <span class="form-tip">开启后用户必须完成实名认证才能使用</span>
          </el-form-item>
        </template>

        <!-- 账号注销配置 -->
        <el-divider content-position="left">
          <el-icon><Delete /></el-icon>
          <span style="margin-left: 8px;">账号注销配置</span>
        </el-divider>

        <el-form-item label="允许账号注销">
          <el-switch v-model="formData.account_deletion.enabled" />
        </el-form-item>

        <template v-if="formData.account_deletion.enabled">
          <el-form-item label="注销冷静期">
            <el-input-number v-model="formData.account_deletion.cooling_period" :min="0" :max="30" />
            <span class="form-tip">天，0表示立即注销</span>
          </el-form-item>

          <el-form-item label="注销前置条件">
            <el-checkbox-group v-model="formData.account_deletion.preconditions">
              <el-checkbox label="clear_data">清空个人数据</el-checkbox>
              <el-checkbox label="unbind_third_party">解绑第三方账号</el-checkbox>
              <el-checkbox label="settle_balance">结清余额</el-checkbox>
              <el-checkbox label="cancel_orders">取消所有订单</el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <el-form-item label="注销确认方式">
            <el-select v-model="formData.account_deletion.confirm_method" placeholder="请选择">
              <el-option label="短信验证码" value="sms" />
              <el-option label="邮箱验证码" value="email" />
              <el-option label="密码确认" value="password" />
            </el-select>
          </el-form-item>
        </template>

        <!-- 用户隐私设置 -->
        <el-divider content-position="left">
          <el-icon><Lock /></el-icon>
          <span style="margin-left: 8px;">用户隐私设置</span>
        </el-divider>

        <el-form-item label="隐私协议">
          <el-input v-model="formData.privacy_policy.url" placeholder="请输入隐私协议URL" />
          <el-button style="margin-left: 10px;" @click="handleEditPrivacyPolicy">编辑协议</el-button>
        </el-form-item>

        <el-form-item label="用户协议">
          <el-input v-model="formData.user_agreement.url" placeholder="请输入用户协议URL" />
          <el-button style="margin-left: 10px;" @click="handleEditUserAgreement">编辑协议</el-button>
        </el-form-item>

        <el-form-item label="数据导出功能">
          <el-switch v-model="formData.data_export.enabled" />
          <span class="form-tip">允许用户导出个人数据</span>
        </el-form-item>

        <el-form-item label="数据导出格式" v-if="formData.data_export.enabled">
          <el-checkbox-group v-model="formData.data_export.formats">
            <el-checkbox label="json">JSON</el-checkbox>
            <el-checkbox label="csv">CSV</el-checkbox>
            <el-checkbox label="excel">Excel</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="隐私设置选项">
          <el-checkbox-group v-model="formData.privacy_settings">
            <el-checkbox label="profile_visible">个人资料可见性</el-checkbox>
            <el-checkbox label="activity_visible">活动记录可见性</el-checkbox>
            <el-checkbox label="search_visible">搜索可见性</el-checkbox>
            <el-checkbox label="recommend_visible">推荐可见性</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <!-- 保存按钮 -->
        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="saving" size="large">
            <el-icon><Check /></el-icon>
            <span style="margin-left: 5px;">保存配置</span>
          </el-button>
          <el-button @click="fetchConfig" size="large">
            <el-icon><RefreshRight /></el-icon>
            <span style="margin-left: 5px;">重置</span>
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { User, Connection, Lock, UserFilled, CreditCard, Delete, Check, RefreshRight } from '@element-plus/icons-vue'
import { getModuleConfig, saveModuleConfig, resetModuleConfig } from '@/api/module'

const props = defineProps({
  appId: {
    type: Number,
    required: true
  }
})

const loading = ref(false)
const saving = ref(false)
const formRef = ref(null)

const formData = reactive({
  // 注册配置
  allow_register: true,
  register_methods: ['phone', 'email'],
  require_real_name: false,
  register_captcha: true,

  // 第三方登录配置
  wechat_login: {
    enabled: false,
    app_id: '',
    app_secret: ''
  },
  qq_login: {
    enabled: false,
    app_id: '',
    app_key: ''
  },
  weibo_login: {
    enabled: false,
    app_key: '',
    app_secret: ''
  },
  apple_login: {
    enabled: false,
    service_id: '',
    team_id: '',
    key_id: ''
  },
  douyin_login: {
    enabled: false,
    client_key: '',
    client_secret: ''
  },

  // 登录配置
  password_min_length: 8,
  password_complexity: ['number', 'letter'],
  login_lock_enabled: true,
  login_lock_threshold: 5,
  login_lock_duration: 30,

  // 用户信息管理
  required_fields: ['nickname'],
  allow_change_username: false,
  nickname_filter: true,
  avatar_review: false,

  // 实名认证配置
  real_name_auth: {
    enabled: false,
    methods: ['idcard'],
    provider: 'aliyun',
    app_id: '',
    app_secret: '',
    required: false
  },

  // 账号注销配置
  account_deletion: {
    enabled: true,
    cooling_period: 7,
    preconditions: ['clear_data', 'unbind_third_party'],
    confirm_method: 'sms'
  },

  // 用户隐私设置
  privacy_policy: {
    url: ''
  },
  user_agreement: {
    url: ''
  },
  data_export: {
    enabled: true,
    formats: ['json', 'csv']
  },
  privacy_settings: ['profile_visible', 'activity_visible']
})

const fetchConfig = async () => {
  loading.value = true
  try {
    const res = await getModuleConfig(props.appId, 'user')
    // request.js已解包，res直接是数据对象
    if (res && res.config) {
      Object.assign(formData, res.config)
    }
  } catch (error) {
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await saveModuleConfig(props.appId, 'user', {
      config: formData
    })
    ElMessage.success('保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const handleEditPrivacyPolicy = () => {
  ElMessage.info('请在新窗口编辑隐私协议')
  // TODO: 打开协议编辑器
}

const handleEditUserAgreement = () => {
  ElMessage.info('请在新窗口编辑用户协议')
  // TODO: 打开协议编辑器
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped lang="scss">
.module-config-container {
  padding: 20px;

  .config-header {
    margin-bottom: 24px;
    
    h2 {
      font-size: 20px;
      font-weight: 600;
      margin-bottom: 8px;
    }

    .config-description {
      color: #909399;
      font-size: 14px;
    }
  }

  .form-tip {
    margin-left: 10px;
    color: #909399;
    font-size: 12px;
  }

  :deep(.el-divider) {
    margin: 30px 0 20px 0;
    
    .el-divider__text {
      display: flex;
      align-items: center;
      font-size: 16px;
      font-weight: 600;
      color: #303133;
    }
  }

  :deep(.el-form-item) {
    margin-bottom: 20px;
  }

  :deep(.el-checkbox-group) {
    .el-checkbox {
      margin-right: 20px;
      margin-bottom: 10px;
    }
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .module-config-container {
    padding: 12px;

    .config-header {
      h2 {
        font-size: 18px;
      }

      .config-description {
        font-size: 12px;
        line-height: 1.5;
      }
    }

    :deep(.el-form) {
      --el-form-label-font-size: 13px;
    }

    :deep(.el-form-item) {
      flex-direction: column !important;
      align-items: flex-start !important;
      margin-bottom: 16px !important;
    }

    :deep(.el-form-item__label) {
      width: 100% !important;
      text-align: left !important;
      margin-bottom: 8px !important;
      padding-right: 0 !important;
      line-height: 1.4 !important;
      white-space: normal !important;
    }

    :deep(.el-form-item__content) {
      width: 100% !important;
      margin-left: 0 !important;
      flex-wrap: wrap !important;
    }

    :deep(.el-checkbox-group) {
      display: flex !important;
      flex-direction: column !important;
      gap: 8px !important;
      
      .el-checkbox {
        margin-right: 0 !important;
      }
    }

    .form-tip {
      display: block !important;
      width: 100% !important;
      margin-top: 4px !important;
      margin-left: 0 !important;
      font-size: 11px;
      line-height: 1.4;
    }

    :deep(.el-input),
    :deep(.el-select) {
      width: 100% !important;
    }

    :deep(.el-input-number) {
      width: 100% !important;
      max-width: 150px !important;
    }

    :deep(.el-divider__text) {
      font-size: 13px;
    }

    :deep(.el-button) {
      margin-left: 0 !important;
      margin-bottom: 8px !important;
    }
  }
}

@media (max-width: 480px) {
  .module-config-container {
    padding: 8px;

    .config-header {
      h2 {
        font-size: 16px;
      }
    }

    :deep(.el-form-item__label) {
      font-size: 13px !important;
    }
  }
}
</style>
