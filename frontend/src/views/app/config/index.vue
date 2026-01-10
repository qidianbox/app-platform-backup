<template>
  <div class="app-detail">
    <!-- 顶部导航栏 -->
    <div class="top-header">
      <!-- 左侧：Logo + APP信息 -->
      <div class="header-logo">
        <div class="app-icon">
          <span>拓</span>
        </div>
        <span class="app-name">拓客APP中台</span>
      </div>
      
      <!-- 右侧：工作台 | 配置中心 Tab -->
      <div class="header-nav">
        <div 
          class="nav-item" 
          :class="{ active: activeTab === 'workspace' }"
          @click="activeTab = 'workspace'"
        >
          工作台
        </div>
        <div 
          class="nav-item" 
          :class="{ active: activeTab === 'config' }"
          @click="activeTab = 'config'"
        >
          配置中心
        </div>
      </div>
    </div>

    <div class="main-container">
      <!-- 左侧边栏 - 仅在配置中心模式显示 -->
      <div class="sidebar" v-show="activeTab === 'config'">
        <div class="sidebar-menu">
          <!-- 概览 -->
          <div 
            class="sidebar-item"
            :class="{ active: currentPage === 'overview' }"
            @click="switchPage('overview')"
          >
            <el-icon><House /></el-icon>
            <span>概览</span>
          </div>
          
          <!-- 基础配置 -->
          <div 
            class="sidebar-item"
            :class="{ active: currentPage === 'basic' }"
            @click="switchPage('basic')"
        >
          <el-icon><Setting /></el-icon>
          <span>基础配置</span>
        </div>

        <!-- 模块分组 -->
        <template v-for="group in moduleGroups" :key="group.key">
          <div 
            v-if="hasModulesInGroup(group.key)"
            class="sidebar-group"
          >
            <div 
              class="group-header"
              @click="toggleGroup(group.key)"
            >
              <div class="group-title">
                <el-icon><component :is="group.icon" /></el-icon>
                <span>{{ group.name }}</span>
              </div>
              <el-icon class="expand-icon" :class="{ expanded: expandedGroups.includes(group.key) }">
                <ArrowRight />
              </el-icon>
            </div>
            <div v-show="expandedGroups.includes(group.key)" class="group-items">
              <div 
                v-for="module in getModulesInGroup(group.key)" 
                :key="module.source_module"
                class="sidebar-item sub-item"
                :class="{ active: currentPage === module.source_module }"
                @click="switchPage(module.source_module)"
              >
                <span>{{ module.name }}</span>
              </div>
            </div>
          </div>
        </template>
        </div>
        <div class="sidebar-footer">
          <div class="sidebar-item back-item" @click="$router.push('/apps')">
            <el-icon><ArrowLeft /></el-icon>
            <span>返回APP列表</span>
          </div>
        </div>
      </div>

      <!-- 右侧内容区 -->
      <div class="content-area">
        <!-- 工作台模式 -->
        <template v-if="activeTab === 'workspace'">
          <Workspace :app-id="appId" :app-info="appInfo" />
        </template>

        <!-- 配置中心模式 -->
        <template v-else-if="activeTab === 'config'">
        <!-- 概览页面 -->
        <div v-if="currentPage === 'overview'" class="page-content">
          <h2 class="page-title">APP概览</h2>
          
          <!-- 统计卡片 -->
          <div class="stats-cards">
            <div class="stat-card">
              <div class="stat-icon users"><el-icon><User /></el-icon></div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.userCount }}</div>
                <div class="stat-label">用户数</div>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon modules"><el-icon><Grid /></el-icon></div>
              <div class="stat-info">
                <div class="stat-value">{{ appModules.length }}</div>
                <div class="stat-label">启用模块</div>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon requests"><el-icon><DataLine /></el-icon></div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.todayRequests }}</div>
                <div class="stat-label">今日请求</div>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon errors"><el-icon><Warning /></el-icon></div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.todayErrors }}</div>
                <div class="stat-label">今日异常</div>
              </div>
            </div>
          </div>

          <!-- APP信息 -->
          <div class="info-section">
            <h3>APP信息</h3>
            <div class="info-grid">
              <div class="info-item">
                <label>APP名称</label>
                <span>{{ appInfo.name }}</span>
              </div>
              <div class="info-item">
                <label>APP标识</label>
                <span class="copyable" @click="copyText(appInfo.app_id)">
                  {{ appInfo.app_id }}
                  <el-icon><CopyDocument /></el-icon>
                </span>
              </div>
              <div class="info-item">
                <label>AppSecret</label>
                <span class="copyable" @click="copyText(appInfo.app_secret)">
                  {{ maskSecret(appInfo.app_secret) }}
                  <el-icon><CopyDocument /></el-icon>
                </span>
              </div>
              <div class="info-item">
                <label>包名</label>
                <span>{{ appInfo.package_name || '-' }}</span>
              </div>
              <div class="info-item">
                <label>状态</label>
                <el-tag :type="appInfo.status === 1 ? 'success' : 'info'" size="small">
                  {{ appInfo.status === 1 ? '正常' : '禁用' }}
                </el-tag>
              </div>
              <div class="info-item">
                <label>创建时间</label>
                <span>{{ formatDate(appInfo.created_at) }}</span>
              </div>
            </div>
          </div>

          <!-- 已启用模块 -->
          <div class="info-section">
            <h3>已启用模块</h3>
            <div class="module-tags">
              <el-tag 
                v-for="module in appModules" 
                :key="module.id"
                type="primary"
                effect="plain"
              >
                {{ module.module_name || moduleNameMap[module.module_code] || module.name }}
              </el-tag>
              <el-empty v-if="appModules.length === 0" description="暂无启用模块" />
            </div>
          </div>
        </div>

        <!-- 基础配置页面 -->
        <div v-else-if="currentPage === 'basic'" class="page-content">
          <h2 class="page-title">基础配置</h2>
          <p class="page-desc">配置APP的基本信息和通用设置</p>
          
          <el-form :model="basicConfig" label-width="140px" class="config-form">
            <div class="form-section">
              <h4>基本信息</h4>
              <el-form-item label="APP名称">
                <el-input v-model="basicConfig.name" placeholder="请输入APP名称" />
              </el-form-item>
              <el-form-item label="APP描述">
                <el-input v-model="basicConfig.description" type="textarea" :rows="3" placeholder="请输入APP描述" />
              </el-form-item>
              <el-form-item label="包名">
                <el-input v-model="basicConfig.package_name" placeholder="如：com.example.app" />
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>安全设置</h4>
              <el-form-item label="启用签名验证">
                <el-switch v-model="basicConfig.enableSignature" />
                <span class="form-hint">启用后所有API请求需要携带签名</span>
              </el-form-item>
              <el-form-item label="IP白名单">
                <el-input v-model="basicConfig.ipWhitelist" type="textarea" :rows="2" placeholder="每行一个IP，留空表示不限制" />
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveBasicConfig">保存配置</el-button>
              <el-button @click="resetBasicConfig">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 用户管理配置 -->
        <div v-else-if="currentPage === 'user_management'" class="page-content">
          <h2 class="page-title">用户管理配置</h2>
          <p class="page-desc">配置用户注册、登录和管理相关设置</p>
          
          <el-form :model="userConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>🔐 登录配置</h4>
              <el-form-item label="密码最小长度">
                <el-input-number v-model="userConfig.passwordMinLength" :min="6" :max="32" />
                <span class="form-hint">建议8位以上</span>
              </el-form-item>
              <el-form-item label="密码复杂度要求">
                <el-checkbox-group v-model="userConfig.passwordRequirements">
                  <el-checkbox label="number">必须包含数字</el-checkbox>
                  <el-checkbox label="letter">必须包含字母</el-checkbox>
                  <el-checkbox label="special">必须包含特殊字符</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label="登录失败锁定">
                <el-switch v-model="userConfig.enableLoginLock" />
                <span class="form-hint">防止暴力破解</span>
              </el-form-item>
              <el-form-item v-if="userConfig.enableLoginLock" label="失败次数限制">
                <el-input-number v-model="userConfig.maxLoginAttempts" :min="3" :max="10" />
                <span class="form-hint">次</span>
              </el-form-item>
              <el-form-item v-if="userConfig.enableLoginLock" label="锁定时长">
                <el-input-number v-model="userConfig.lockDuration" :min="5" :max="1440" />
                <span class="form-hint">分钟</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>👤 用户信息管理</h4>
              <el-form-item label="必填字段">
                <el-checkbox-group v-model="userConfig.requiredFields">
                  <el-checkbox label="nickname">昵称</el-checkbox>
                  <el-checkbox label="avatar">头像</el-checkbox>
                  <el-checkbox label="gender">性别</el-checkbox>
                  <el-checkbox label="birthday">生日</el-checkbox>
                  <el-checkbox label="region">地区</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label="允许修改用户名">
                <el-switch v-model="userConfig.allowChangeUsername" />
                <span class="form-hint">关闭后用户名不可修改</span>
              </el-form-item>
              <el-form-item label="昵称敏感词过滤">
                <el-switch v-model="userConfig.enableNicknameFilter" />
              </el-form-item>
              <el-form-item label="头像审核">
                <el-switch v-model="userConfig.enableAvatarReview" />
                <span class="form-hint">自动检测头像是否违规</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>🪪 实名认证配置</h4>
              <el-form-item label="启用实名认证">
                <el-switch v-model="userConfig.enableRealName" />
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>🗑️ 账号注销配置</h4>
              <el-form-item label="允许账号注销">
                <el-switch v-model="userConfig.allowAccountDeletion" />
              </el-form-item>
              <el-form-item v-if="userConfig.allowAccountDeletion" label="注销冷静期">
                <el-input-number v-model="userConfig.deletionCooldown" :min="0" :max="30" />
                <span class="form-hint">天，0表示立即注销</span>
              </el-form-item>
              <el-form-item v-if="userConfig.allowAccountDeletion" label="注销前置条件">
                <el-checkbox-group v-model="userConfig.deletionRequirements">
                  <el-checkbox label="clearData">清空个人数据</el-checkbox>
                  <el-checkbox label="unbindThirdParty">解除第三方账号</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item v-if="userConfig.allowAccountDeletion" label="注销确认方式">
                <el-radio-group v-model="userConfig.deletionConfirmMethod">
                  <el-radio label="sms">短信验证码</el-radio>
                  <el-radio label="password">密码验证</el-radio>
                </el-radio-group>
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('user_management')">保存配置</el-button>
              <el-button @click="testConfig('user_management')">测试配置</el-button>
              <el-button @click="resetConfig('user_management')">重置</el-button>
              <el-button @click="showConfigHistory('user_management')">查看历史</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 消息中心配置 -->
        <div v-else-if="currentPage === 'message_center'" class="page-content">
          <h2 class="page-title">消息中心配置</h2>
          <p class="page-desc">配置站内消息和通知相关设置</p>
          
          <el-form :model="messageConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>📬 基础配置</h4>
              <el-form-item label="启用消息服务">
                <el-switch v-model="messageConfig.enabled" />
              </el-form-item>
              <el-form-item label="消息保留天数">
                <el-input-number v-model="messageConfig.retentionDays" :min="7" :max="365" />
                <span class="form-hint">天</span>
              </el-form-item>
              <el-form-item label="单用户消息上限">
                <el-input-number v-model="messageConfig.maxMessagesPerUser" :min="100" :max="10000" />
                <span class="form-hint">条</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>📝 消息类型</h4>
              <el-form-item label="支持的消息类型">
                <el-checkbox-group v-model="messageConfig.supportedTypes">
                  <el-checkbox label="system">系统通知</el-checkbox>
                  <el-checkbox label="activity">活动消息</el-checkbox>
                  <el-checkbox label="transaction">交易消息</el-checkbox>
                  <el-checkbox label="social">社交消息</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('message_center')">保存配置</el-button>
              <el-button @click="resetConfig('message_center')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 推送服务配置 -->
        <div v-else-if="currentPage === 'push_service'" class="page-content">
          <h2 class="page-title">推送服务配置</h2>
          <p class="page-desc">配置APP推送通知服务</p>
          
          <el-form :model="pushConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>🔔 基础配置</h4>
              <el-form-item label="启用推送服务">
                <el-switch v-model="pushConfig.enabled" />
              </el-form-item>
              <el-form-item label="推送服务商">
                <el-select v-model="pushConfig.provider" placeholder="请选择">
                  <el-option label="极光推送" value="jpush" />
                  <el-option label="个推" value="getui" />
                  <el-option label="友盟推送" value="umeng" />
                  <el-option label="Firebase" value="firebase" />
                </el-select>
              </el-form-item>
              <el-form-item label="AppKey">
                <el-input v-model="pushConfig.appKey" placeholder="请输入AppKey" />
              </el-form-item>
              <el-form-item label="MasterSecret">
                <el-input v-model="pushConfig.masterSecret" type="password" placeholder="请输入MasterSecret" show-password />
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>⏰ 推送策略</h4>
              <el-form-item label="静默时段">
                <el-switch v-model="pushConfig.enableQuietHours" />
                <span class="form-hint">在指定时段不发送推送</span>
              </el-form-item>
              <el-form-item v-if="pushConfig.enableQuietHours" label="静默时间">
                <el-time-picker v-model="pushConfig.quietStart" placeholder="开始时间" format="HH:mm" />
                <span style="margin: 0 8px;">至</span>
                <el-time-picker v-model="pushConfig.quietEnd" placeholder="结束时间" format="HH:mm" />
              </el-form-item>
              <el-form-item label="每日推送上限">
                <el-input-number v-model="pushConfig.dailyLimit" :min="1" :max="100" />
                <span class="form-hint">条/用户</span>
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('push_service')">保存配置</el-button>
              <el-button @click="testConfig('push_service')">测试推送</el-button>
              <el-button @click="resetConfig('push_service')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 支付中心配置 -->
        <div v-else-if="currentPage === 'payment'" class="page-content">
          <h2 class="page-title">支付中心配置</h2>
          <p class="page-desc">配置支付渠道和安全设置</p>
          
          <el-form :model="paymentConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>🔐 安全验证配置</h4>
              <el-form-item label="启用安全验证">
                <el-switch v-model="paymentConfig.enableSecurityVerify" />
              </el-form-item>
              <el-form-item v-if="paymentConfig.enableSecurityVerify" label="验证方式">
                <el-checkbox-group v-model="paymentConfig.verifyMethods">
                  <el-checkbox label="password">支付密码</el-checkbox>
                  <el-checkbox label="fingerprint">指纹识别</el-checkbox>
                  <el-checkbox label="face">面容识别</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label="验证触发金额">
                <el-input-number v-model="paymentConfig.verifyThreshold" :min="0" :max="100000" />
                <span class="form-hint">元，0表示所有支付都需要验证</span>
              </el-form-item>
              <el-form-item label="密码错误锁定">
                <el-input-number v-model="paymentConfig.maxPasswordAttempts" :min="3" :max="10" />
                <span class="form-hint">次</span>
              </el-form-item>
              <el-form-item label="锁定时长">
                <el-input-number v-model="paymentConfig.lockDuration" :min="5" :max="1440" />
                <span class="form-hint">分钟</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>💰 限额控制配置</h4>
              <el-form-item label="启用限额控制">
                <el-switch v-model="paymentConfig.enableLimitControl" />
              </el-form-item>
              <el-form-item v-if="paymentConfig.enableLimitControl" label="单笔支付限额">
                <el-input-number v-model="paymentConfig.singleLimit" :min="0" :max="1000000" />
                <span class="form-hint">元，0表示不限制</span>
              </el-form-item>
              <el-form-item v-if="paymentConfig.enableLimitControl" label="每日支付限额">
                <el-input-number v-model="paymentConfig.dailyLimit" :min="0" :max="10000000" />
                <span class="form-hint">元，0表示不限制</span>
              </el-form-item>
              <el-form-item v-if="paymentConfig.enableLimitControl" label="每月支付限额">
                <el-input-number v-model="paymentConfig.monthlyLimit" :min="0" :max="100000000" />
                <span class="form-hint">元，0表示不限制</span>
              </el-form-item>
              <el-form-item v-if="paymentConfig.enableLimitControl" label="每日支付次数">
                <el-input-number v-model="paymentConfig.dailyCount" :min="0" :max="1000" />
                <span class="form-hint">次，0表示不限制</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>🔗 回调配置</h4>
              <el-form-item label="支付成功回调">
                <el-input v-model="paymentConfig.successCallback" placeholder="请输入支付成功回调地址" />
              </el-form-item>
              <el-form-item label="支付失败回调">
                <el-input v-model="paymentConfig.failCallback" placeholder="请输入支付失败回调地址" />
              </el-form-item>
              <el-form-item label="退款回调">
                <el-input v-model="paymentConfig.refundCallback" placeholder="请输入退款回调地址" />
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>⚙️ 其他配置</h4>
              <el-form-item label="支付超时时间">
                <el-input-number v-model="paymentConfig.timeout" :min="5" :max="60" />
                <span class="form-hint">分钟</span>
              </el-form-item>
              <el-form-item label="启用自动退款">
                <el-switch v-model="paymentConfig.enableAutoRefund" />
                <span class="form-hint">订单超时自动退款</span>
              </el-form-item>
              <el-form-item label="启用支付日志">
                <el-switch v-model="paymentConfig.enablePaymentLog" />
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('payment')">保存配置</el-button>
              <el-button @click="testConfig('payment')">测试配置</el-button>
              <el-button @click="resetConfig('payment')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 短信服务配置 -->
        <div v-else-if="currentPage === 'sms_service'" class="page-content">
          <h2 class="page-title">短信服务配置</h2>
          <p class="page-desc">配置短信发送服务和验证码设置</p>
          
          <el-form :model="smsConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>📱 基础配置</h4>
              <el-form-item label="启用短信服务">
                <el-switch v-model="smsConfig.enabled" />
              </el-form-item>
              <el-form-item label="短信服务提供商">
                <el-select v-model="smsConfig.provider" placeholder="请选择">
                  <el-option label="阿里云短信" value="aliyun" />
                  <el-option label="腾讯云短信" value="tencent" />
                  <el-option label="华为云短信" value="huawei" />
                </el-select>
              </el-form-item>
              <el-form-item label="AccessKey">
                <el-input v-model="smsConfig.accessKey" placeholder="请输入AccessKey" />
              </el-form-item>
              <el-form-item label="SecretKey">
                <el-input v-model="smsConfig.secretKey" type="password" placeholder="请输入SecretKey" show-password />
              </el-form-item>
              <el-form-item label="短信签名">
                <el-input v-model="smsConfig.signName" placeholder="例如：我的应用" />
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>🔢 验证码短信配置</h4>
              <el-form-item label="验证码长度">
                <el-input-number v-model="smsConfig.codeLength" :min="4" :max="8" />
                <span class="form-hint">位</span>
              </el-form-item>
              <el-form-item label="验证码有效期">
                <el-input-number v-model="smsConfig.codeExpiry" :min="1" :max="30" />
                <span class="form-hint">分钟</span>
              </el-form-item>
              <el-form-item label="验证码模板ID">
                <el-input v-model="smsConfig.codeTemplateId" placeholder="例如：SMS_123456789" />
              </el-form-item>
              <el-form-item label="发送间隔">
                <el-input-number v-model="smsConfig.sendInterval" :min="30" :max="300" />
                <span class="form-hint">秒</span>
              </el-form-item>
              <el-form-item label="每日发送限制">
                <el-input-number v-model="smsConfig.dailyLimit" :min="1" :max="50" />
                <span class="form-hint">条/手机号</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>📢 通知短信配置</h4>
              <el-form-item label="启用通知短信">
                <el-switch v-model="smsConfig.enableNotification" />
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>⚙️ 高级配置</h4>
              <el-form-item label="失败重试次数">
                <el-input-number v-model="smsConfig.retryCount" :min="0" :max="5" />
                <span class="form-hint">次</span>
              </el-form-item>
              <el-form-item label="请求超时时间">
                <el-input-number v-model="smsConfig.timeout" :min="5" :max="60" />
                <span class="form-hint">秒</span>
              </el-form-item>
              <el-form-item label="状态回调URL">
                <el-input v-model="smsConfig.callbackUrl" placeholder="请输入状态回调地址" />
              </el-form-item>
              <el-form-item label="余额告警阈值">
                <el-input-number v-model="smsConfig.balanceAlert" :min="100" :max="10000" />
                <span class="form-hint">条</span>
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('sms_service')">保存配置</el-button>
              <el-button @click="testConfig('sms_service')">测试发送</el-button>
              <el-button @click="resetConfig('sms_service')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 数据埋点配置 -->
        <div v-else-if="currentPage === 'data_tracking'" class="page-content">
          <h2 class="page-title">数据埋点配置</h2>
          <p class="page-desc">配置用户行为埋点和数据分析</p>
          
          <el-form :model="trackingConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>📊 基础配置</h4>
              <el-form-item label="启用数据埋点">
                <el-switch v-model="trackingConfig.enabled" />
              </el-form-item>
              <el-form-item label="数据上报方式">
                <el-radio-group v-model="trackingConfig.reportMethod">
                  <el-radio label="realtime">实时上报</el-radio>
                  <el-radio label="batch">批量上报</el-radio>
                </el-radio-group>
              </el-form-item>
              <el-form-item v-if="trackingConfig.reportMethod === 'batch'" label="批量上报间隔">
                <el-input-number v-model="trackingConfig.batchInterval" :min="10" :max="300" />
                <span class="form-hint">秒</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>🎯 事件配置</h4>
              <el-form-item label="自动采集事件">
                <el-checkbox-group v-model="trackingConfig.autoEvents">
                  <el-checkbox label="pageView">页面浏览</el-checkbox>
                  <el-checkbox label="click">点击事件</el-checkbox>
                  <el-checkbox label="scroll">滚动事件</el-checkbox>
                  <el-checkbox label="error">错误事件</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('data_tracking')">保存配置</el-button>
              <el-button @click="resetConfig('data_tracking')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 日志服务配置 -->
        <div v-else-if="currentPage === 'log_service'" class="page-content">
          <h2 class="page-title">日志服务配置</h2>
          <p class="page-desc">配置日志收集、存储、分析等功能</p>
          
          <el-form :model="logConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>📝 基础配置</h4>
              <el-form-item label="启用日志服务">
                <el-switch v-model="logConfig.enabled" />
              </el-form-item>
              <el-form-item label="日志级别">
                <el-select v-model="logConfig.level" placeholder="请选择">
                  <el-option label="DEBUG" value="debug" />
                  <el-option label="INFO" value="info" />
                  <el-option label="WARN" value="warn" />
                  <el-option label="ERROR" value="error" />
                </el-select>
              </el-form-item>
              <el-form-item label="日志存储方式">
                <el-checkbox-group v-model="logConfig.storageTypes">
                  <el-checkbox label="local">本地存储</el-checkbox>
                  <el-checkbox label="cloud">云端存储</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label="日志保留时间">
                <el-input-number v-model="logConfig.retentionDays" :min="7" :max="365" />
                <span class="form-hint">天</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>📤 上报配置</h4>
              <el-form-item label="实时上报">
                <el-switch v-model="logConfig.realtimeReport" />
              </el-form-item>
              <el-form-item label="批量上报数量">
                <el-input-number v-model="logConfig.batchSize" :min="10" :max="1000" />
                <span class="form-hint">条</span>
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('log_service')">保存配置</el-button>
              <el-button @click="resetConfig('log_service')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 监控告警配置 -->
        <div v-else-if="currentPage === 'monitor_alert'" class="page-content">
          <h2 class="page-title">监控告警配置</h2>
          <p class="page-desc">配置应用监控和告警通知</p>
          
          <el-form :model="monitorConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>📡 监控配置</h4>
              <el-form-item label="启用监控服务">
                <el-switch v-model="monitorConfig.enabled" />
              </el-form-item>
              <el-form-item label="监控指标">
                <el-checkbox-group v-model="monitorConfig.metrics">
                  <el-checkbox label="cpu">CPU使用率</el-checkbox>
                  <el-checkbox label="memory">内存使用率</el-checkbox>
                  <el-checkbox label="api">API响应时间</el-checkbox>
                  <el-checkbox label="error">错误率</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label="采集间隔">
                <el-input-number v-model="monitorConfig.interval" :min="10" :max="300" />
                <span class="form-hint">秒</span>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>🚨 告警配置</h4>
              <el-form-item label="启用告警">
                <el-switch v-model="monitorConfig.alertEnabled" />
              </el-form-item>
              <el-form-item v-if="monitorConfig.alertEnabled" label="告警方式">
                <el-checkbox-group v-model="monitorConfig.alertMethods">
                  <el-checkbox label="email">邮件</el-checkbox>
                  <el-checkbox label="sms">短信</el-checkbox>
                  <el-checkbox label="webhook">Webhook</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item v-if="monitorConfig.alertEnabled" label="告警接收人">
                <el-input v-model="monitorConfig.alertReceivers" type="textarea" :rows="2" placeholder="多个接收人用逗号分隔" />
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('monitor_alert')">保存配置</el-button>
              <el-button @click="testConfig('monitor_alert')">测试告警</el-button>
              <el-button @click="resetConfig('monitor_alert')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 文件存储配置 -->
        <div v-else-if="currentPage === 'file_storage'" class="page-content">
          <h2 class="page-title">文件存储配置</h2>
          <p class="page-desc">配置文件上传、下载和存储服务</p>
          
          <el-form :model="storageConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>☁️ 存储配置</h4>
              <el-form-item label="启用文件存储">
                <el-switch v-model="storageConfig.enabled" />
              </el-form-item>
              <el-form-item label="存储服务商">
                <el-select v-model="storageConfig.provider" placeholder="请选择">
                  <el-option label="阿里云OSS" value="aliyun" />
                  <el-option label="腾讯云COS" value="tencent" />
                  <el-option label="七牛云" value="qiniu" />
                  <el-option label="AWS S3" value="aws" />
                </el-select>
              </el-form-item>
              <el-form-item label="Bucket名称">
                <el-input v-model="storageConfig.bucket" placeholder="请输入Bucket名称" />
              </el-form-item>
              <el-form-item label="AccessKey">
                <el-input v-model="storageConfig.accessKey" placeholder="请输入AccessKey" />
              </el-form-item>
              <el-form-item label="SecretKey">
                <el-input v-model="storageConfig.secretKey" type="password" placeholder="请输入SecretKey" show-password />
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>📁 上传限制</h4>
              <el-form-item label="最大文件大小">
                <el-input-number v-model="storageConfig.maxFileSize" :min="1" :max="1024" />
                <span class="form-hint">MB</span>
              </el-form-item>
              <el-form-item label="允许的文件类型">
                <el-input v-model="storageConfig.allowedTypes" placeholder="例如：jpg,png,pdf" />
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('file_storage')">保存配置</el-button>
              <el-button @click="testConfig('file_storage')">测试连接</el-button>
              <el-button @click="resetConfig('file_storage')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 配置管理配置 -->
        <div v-else-if="currentPage === 'config_management'" class="page-content">
          <h2 class="page-title">配置管理</h2>
          <p class="page-desc">管理远程配置下发和动态配置</p>
          
          <el-form :model="configMgmtConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>⚙️ 基础配置</h4>
              <el-form-item label="启用配置管理">
                <el-switch v-model="configMgmtConfig.enabled" />
              </el-form-item>
              <el-form-item label="配置刷新间隔">
                <el-input-number v-model="configMgmtConfig.refreshInterval" :min="60" :max="3600" />
                <span class="form-hint">秒</span>
              </el-form-item>
              <el-form-item label="启用配置缓存">
                <el-switch v-model="configMgmtConfig.enableCache" />
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('config_management')">保存配置</el-button>
              <el-button @click="resetConfig('config_management')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 版本管理配置 -->
        <div v-else-if="currentPage === 'version_management'" class="page-content">
          <h2 class="page-title">版本管理配置</h2>
          <p class="page-desc">配置APP版本发布和更新策略</p>
          
          <el-form :model="versionConfig" label-width="160px" class="config-form">
            <div class="form-section">
              <h4>📦 更新配置</h4>
              <el-form-item label="启用版本管理">
                <el-switch v-model="versionConfig.enabled" />
              </el-form-item>
              <el-form-item label="强制更新">
                <el-switch v-model="versionConfig.forceUpdate" />
                <span class="form-hint">开启后用户必须更新到最新版本</span>
              </el-form-item>
              <el-form-item label="更新提示方式">
                <el-radio-group v-model="versionConfig.promptType">
                  <el-radio label="dialog">弹窗提示</el-radio>
                  <el-radio label="toast">轻提示</el-radio>
                  <el-radio label="silent">静默更新</el-radio>
                </el-radio-group>
              </el-form-item>
            </div>

            <div class="form-section">
              <h4>📥 下载配置</h4>
              <el-form-item label="Android下载地址">
                <el-input v-model="versionConfig.androidUrl" placeholder="请输入Android安装包下载地址" />
              </el-form-item>
              <el-form-item label="iOS下载地址">
                <el-input v-model="versionConfig.iosUrl" placeholder="请输入iOS App Store地址" />
              </el-form-item>
            </div>

            <el-form-item>
              <el-button type="primary" @click="saveModuleConfig('version_management')">保存配置</el-button>
              <el-button @click="resetConfig('version_management')">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 默认页面 -->
        <div v-else class="page-content">
          <el-empty description="请从左侧选择配置项" />
        </div>
        </template>
      </div>
    </div>

    <!-- 配置历史记录对话框 -->
    <el-dialog 
      v-model="historyDialogVisible" 
      title="配置历史记录" 
      width="800px"
      :close-on-click-modal="false"
    >
      <el-table 
        :data="configHistory" 
        v-loading="loadingHistory"
        style="width: 100%"
      >
        <el-table-column prop="version" label="版本" width="80" />
        <el-table-column prop="operator" label="操作人" width="120" />
        <el-table-column label="配置内容" min-width="200">
          <template #default="{ row }">
            <pre style="margin: 0; font-size: 12px; max-height: 100px; overflow: auto;">{{ formatConfig(row.config) }}</pre>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button 
              type="primary" 
              size="small" 
              @click="rollbackToHistory(row.id)"
            >
              回滚
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <template #footer>
        <el-button @click="historyDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { 
  ArrowLeft, ArrowRight, House, Setting, User, UserFilled, 
  CreditCard, ChatDotRound, DataLine, Document, Monitor, 
  FolderOpened, Tools, Box, Grid, Warning, CopyDocument 
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import Workspace from './Workspace.vue'

const route = useRoute()
const appId = computed(() => route.params.id)

const activeTab = ref('config') // 默认显示配置中心
const currentPage = ref('overview')
const expandedGroups = ref(['user', 'message', 'data', 'system', 'storage'])
const adminName = ref(localStorage.getItem('adminName') || 'Admin')

// 模块分组定义
const moduleGroups = [
  { key: 'user', name: '用户与权限', icon: 'UserFilled', modules: ['user_management'] },
  { key: 'payment', name: '交易与支付', icon: 'CreditCard', modules: ['payment'] },
  { key: 'message', name: '消息与通知', icon: 'ChatDotRound', modules: ['message_center', 'push_service', 'sms_service'] },
  { key: 'data', name: '数据与分析', icon: 'DataLine', modules: ['data_tracking'] },
  { key: 'system', name: '系统与运维', icon: 'Monitor', modules: ['log_service', 'monitor_alert'] },
  { key: 'storage', name: '存储服务', icon: 'FolderOpened', modules: ['file_storage', 'config_management', 'version_management'] }
]

// 模块名称映射
const moduleNameMap = {
  user_management: '用户管理',
  message_center: '消息中心',
  push_service: '推送服务',
  data_tracking: '数据埋点',
  log_service: '日志服务',
  monitor_alert: '监控告警',
  file_storage: '文件存储',
  config_management: '配置管理',
  version_management: '版本管理',
  payment: '支付中心',
  sms_service: '短信服务'
}

const appInfo = ref({
  name: '',
  app_id: '',
  app_secret: '',
  package_name: '',
  description: '',
  status: 1,
  created_at: ''
})

const appModules = ref([])

const stats = ref({
  userCount: 0,
  todayRequests: 0,
  todayErrors: 0
})

// 各模块配置表单
const basicConfig = ref({
  name: '',
  description: '',
  package_name: '',
  enableSignature: false,
  ipWhitelist: ''
})

const userConfig = ref({
  passwordMinLength: 8,
  passwordRequirements: ['number', 'letter'],
  enableLoginLock: true,
  maxLoginAttempts: 5,
  lockDuration: 30,
  requiredFields: ['nickname'],
  allowChangeUsername: false,
  enableNicknameFilter: true,
  enableAvatarReview: false,
  enableRealName: false,
  allowAccountDeletion: true,
  deletionCooldown: 7,
  deletionRequirements: ['clearData'],
  deletionConfirmMethod: 'sms'
})

const messageConfig = ref({
  enabled: true,
  retentionDays: 30,
  maxMessagesPerUser: 1000,
  supportedTypes: ['system', 'activity']
})

const pushConfig = ref({
  enabled: true,
  provider: 'jpush',
  appKey: '',
  masterSecret: '',
  enableQuietHours: false,
  quietStart: null,
  quietEnd: null,
  dailyLimit: 10
})

const paymentConfig = ref({
  enableSecurityVerify: true,
  verifyMethods: ['password'],
  verifyThreshold: 500,
  maxPasswordAttempts: 5,
  lockDuration: 30,
  enableLimitControl: true,
  singleLimit: 50000,
  dailyLimit: 100000,
  monthlyLimit: 500000,
  dailyCount: 100,
  successCallback: '',
  failCallback: '',
  refundCallback: '',
  timeout: 30,
  enableAutoRefund: false,
  enablePaymentLog: true
})

const smsConfig = ref({
  enabled: true,
  provider: 'aliyun',
  accessKey: '',
  secretKey: '',
  signName: '',
  codeLength: 6,
  codeExpiry: 5,
  codeTemplateId: '',
  sendInterval: 60,
  dailyLimit: 10,
  enableNotification: true,
  retryCount: 3,
  timeout: 10,
  callbackUrl: '',
  balanceAlert: 1000
})

const trackingConfig = ref({
  enabled: true,
  reportMethod: 'batch',
  batchInterval: 60,
  autoEvents: ['pageView', 'click']
})

const logConfig = ref({
  enabled: true,
  level: 'info',
  storageTypes: ['local'],
  retentionDays: 30,
  realtimeReport: false,
  batchSize: 100
})

const monitorConfig = ref({
  enabled: true,
  metrics: ['api', 'error'],
  interval: 60,
  alertEnabled: true,
  alertMethods: ['email'],
  alertReceivers: ''
})

const storageConfig = ref({
  enabled: true,
  provider: 'aliyun',
  bucket: '',
  accessKey: '',
  secretKey: '',
  maxFileSize: 100,
  allowedTypes: 'jpg,png,gif,pdf,doc,docx'
})

const configMgmtConfig = ref({
  enabled: true,
  refreshInterval: 300,
  enableCache: true
})

const versionConfig = ref({
  enabled: true,
  forceUpdate: false,
  promptType: 'dialog',
  androidUrl: '',
  iosUrl: ''
})

// 切换页面
const switchPage = (page) => {
  currentPage.value = page
  // 切换到模块配置页面时加载配置
  if (page !== 'overview' && page !== 'basic') {
    loadModuleConfig(page)
  }
}

// 切换分组展开/收起
const toggleGroup = (groupKey) => {
  const index = expandedGroups.value.indexOf(groupKey)
  if (index > -1) {
    expandedGroups.value.splice(index, 1)
  } else {
    expandedGroups.value.push(groupKey)
  }
}

// 检查分组是否有模块
const hasModulesInGroup = (groupKey) => {
  const group = moduleGroups.find(g => g.key === groupKey)
  if (!group) return false
  // 使用module_code匹配（后端返回的字段）
  return appModules.value.some(m => group.modules.includes(m.module_code))
}

// 获取分组内的模块
const getModulesInGroup = (groupKey) => {
  const group = moduleGroups.find(g => g.key === groupKey)
  if (!group) return []
  
  // 去重：使用Map确保每个module_code只出现一次
  const uniqueModules = new Map()
  appModules.value
    .filter(m => group.modules.includes(m.module_code))
    .forEach(m => {
      if (!uniqueModules.has(m.module_code)) {
        uniqueModules.set(m.module_code, {
          ...m,
          source_module: m.module_code, // 兼容侧边栏点击
          name: moduleNameMap[m.module_code] || m.module_name || m.name
        })
      }
    })
  
  return Array.from(uniqueModules.values())
}

// 获取APP信息
const fetchAppInfo = async () => {
  try {
    const res = await request.get(`/apps/${appId.value}`)
    if (res.code === 0 && res.data) {
      appInfo.value = res.data
      basicConfig.value.name = res.data.name
      basicConfig.value.description = res.data.description || ''
      basicConfig.value.package_name = res.data.package_name || ''
    }
  } catch (error) {
    console.error('获取APP信息失败:', error)
  }
}

// 获取APP模块列表
const fetchAppModules = async () => {
  try {
    const res = await request.get(`/apps/${appId.value}/modules`)
    if (res.code === 0 && res.data) {
      appModules.value = res.data
    }
  } catch (error) {
    console.error('获取APP模块失败:', error)
  }
}

// 复制文本
const copyText = (text) => {
  if (!text) return
  navigator.clipboard.writeText(text)
  ElMessage.success('已复制到剪贴板')
}

// 遮盖密钥
const maskSecret = (secret) => {
  if (!secret) return '-'
  if (secret.length <= 8) return '********'
  return secret.substring(0, 4) + '****' + secret.substring(secret.length - 4)
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

// 保存基础配置
const saveBasicConfig = async () => {
  try {
    await request.put(`/apps/${appId.value}`, basicConfig.value)
    ElMessage({
      message: '基础配置保存成功',
      type: 'success',
      duration: 3000,
      showClose: true
    })
    fetchAppInfo()
  } catch (error) {
    ElMessage({
      message: '配置保存失败，请稍后重试',
      type: 'error',
      duration: 5000,
      showClose: true
    })
  }
}

// 重置基础配置
const resetBasicConfig = () => {
  basicConfig.value.name = appInfo.value.name
  basicConfig.value.description = appInfo.value.description || ''
  basicConfig.value.package_name = appInfo.value.package_name || ''
}

// 获取模块配置数据
const getModuleConfigData = (moduleKey) => {
  const configMap = {
    'user_management': userConfig.value,
    'message_center': messageConfig.value,
    'push_service': pushConfig.value,
    'payment': paymentConfig.value,
    'sms_service': smsConfig.value,
    'data_tracking': trackingConfig.value,
    'log_service': logConfig.value,
    'monitor_alert': monitorConfig.value,
    'file_storage': fileConfig.value,
    'config_management': configMgmtConfig.value,
    'version_management': versionConfig.value
  }
  return configMap[moduleKey] || {}
}

// 加载模块配置
const loadModuleConfig = async (moduleKey) => {
  try {
    const res = await request.get(`/apps/${appId.value}/modules/${moduleKey}/config`)
    if (res.code === 0 && res.data && res.data.config) {
      const config = typeof res.data.config === 'string' ? JSON.parse(res.data.config) : res.data.config
      // 更新对应的配置对象
      const configMap = {
        'user_management': userConfig,
        'message_center': messageConfig,
        'push_service': pushConfig,
        'payment': paymentConfig,
        'sms_service': smsConfig,
        'data_tracking': trackingConfig,
        'log_service': logConfig,
        'monitor_alert': monitorConfig,
        'file_storage': fileConfig,
        'config_management': configMgmtConfig,
        'version_management': versionConfig
      }
      if (configMap[moduleKey]) {
        Object.assign(configMap[moduleKey].value, config)
      }
    }
  } catch (error) {
    console.error('加载配置失败:', error)
  }
}

// 保存模块配置
const saveModuleConfig = async (moduleKey) => {
  try {
    const configData = getModuleConfigData(moduleKey)
    await request.put(`/apps/${appId.value}/modules/${moduleKey}/config`, {
      config: configData
    })
    ElMessage({
      message: '配置保存成功',
      type: 'success',
      duration: 3000,
      showClose: true
    })
  } catch (error) {
    console.error('保存配置失败:', error)
    ElMessage({
      message: '配置保存失败，请稍后重试',
      type: 'error',
      duration: 5000,
      showClose: true
    })
  }
}

// 测试配置
const testConfig = (moduleKey) => {
  ElMessage.info('测试功能开发中...')
}

// 重置配置
const resetConfig = async (moduleKey) => {
  try {
    await ElMessageBox.confirm('确定要重置此模块的配置吗？重置后将恢复为默认配置。', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await request.post(`/apps/${appId.value}/modules/${moduleKey}/config/reset`)
    ElMessage.success('配置已重置')
    // 重新加载配置
    await loadModuleConfig(moduleKey)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重置配置失败:', error)
      ElMessage.error('重置失败')
    }
  }
}

// 配置历史记录
const historyDialogVisible = ref(false)
const currentHistoryModule = ref('')
const configHistory = ref([])
const loadingHistory = ref(false)

// 显示配置历史
const showConfigHistory = async (moduleKey) => {
  currentHistoryModule.value = moduleKey
  historyDialogVisible.value = true
  await loadConfigHistory(moduleKey)
}

// 加载配置历史
const loadConfigHistory = async (moduleKey) => {
  loadingHistory.value = true
  try {
    const res = await request.get(`/apps/${appId.value}/modules/${moduleKey}/config/history`)
    if (res.code === 0 && res.data) {
      configHistory.value = res.data
    }
  } catch (error) {
    console.error('加载历史记录失败:', error)
    ElMessage.error('加载失败')
  } finally {
    loadingHistory.value = false
  }
}

// 回滚配置
const rollbackToHistory = async (historyId) => {
  try {
    await ElMessageBox.confirm('确定要回滚到该版本的配置吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await request.post(`/apps/${appId.value}/modules/${currentHistoryModule.value}/config/rollback/${historyId}`)
    ElMessage.success('配置已回滚')
    historyDialogVisible.value = false
    // 重新加载配置
    await loadModuleConfig(currentHistoryModule.value)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('回滚配置失败:', error)
      ElMessage.error('回滚失败')
    }
  }
}

// 格式化配置显示
const formatConfig = (config) => {
  try {
    const configObj = typeof config === 'string' ? JSON.parse(config) : config
    return JSON.stringify(configObj, null, 2)
  } catch {
    return config
  }
}

onMounted(() => {
  fetchAppInfo()
  fetchAppModules()
})
</script>

<style lang="scss" scoped>
.app-detail {
  min-height: 100vh;
  background: #f5f7fa;
}

.top-header {
  display: flex;
  align-items: center;
  padding: 0 24px;
  height: 60px;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.header-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  
  .back-btn {
    color: white;
    &:hover {
      background: rgba(255, 255, 255, 0.1);
    }
  }
  
  .app-icon {
    width: 36px;
    height: 36px;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    font-size: 16px;
  }
  
  .app-info {
    display: flex;
    align-items: center;
    gap: 12px;
  }
  
  .app-name {
    font-size: 18px;
    font-weight: 600;
  }
  
  .app-id {
    font-size: 12px;
    opacity: 0.8;
    background: rgba(255, 255, 255, 0.1);
    padding: 2px 8px;
    border-radius: 4px;
  }
}

.header-nav {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: auto;
  
  .nav-item {
    padding: 8px 20px;
    font-size: 14px;
    font-weight: 500;
    color: rgba(255, 255, 255, 0.7);
    cursor: pointer;
    border-radius: 6px;
    transition: all 0.2s;
    
    &:hover {
      color: white;
      background: rgba(255, 255, 255, 0.1);
    }
    
    &.active {
      color: white;
      background: #409eff;
      font-weight: 600;
    }
  }
}

.main-container {
  display: flex;
  height: calc(100vh - 60px);
}

.sidebar {
  width: 240px;
  background: white;
  border-right: 1px solid #e4e7ed;
  overflow-y: auto;
  padding: 16px 0;
  display: flex;
  flex-direction: column;
}

.sidebar-menu {
  flex: 1;
}

.sidebar-footer {
  border-top: 1px solid #e4e7ed;
  padding-top: 8px;
  margin-top: 8px;
  
  .back-item {
    color: #909399;
    
    &:hover {
      color: #409eff;
      background: #f5f7fa;
    }
  }
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  cursor: pointer;
  color: #606266;
  transition: all 0.2s;
  
  &:hover {
    background: #f5f7fa;
    color: #409eff;
  }
  
  &.active {
    background: #ecf5ff;
    color: #409eff;
    border-right: 3px solid #409eff;
  }
  
  &.sub-item {
    padding-left: 48px;
    font-size: 14px;
  }
}

.sidebar-group {
  .group-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 20px;
    cursor: pointer;
    color: #303133;
    font-weight: 500;
    
    &:hover {
      background: #f5f7fa;
    }
    
    .group-title {
      display: flex;
      align-items: center;
      gap: 10px;
    }
    
    .expand-icon {
      transition: transform 0.2s;
      
      &.expanded {
        transform: rotate(90deg);
      }
    }
  }
  
  .group-items {
    background: #fafafa;
  }
}

.content-area {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.page-content {
  max-width: 900px;
}

.page-title {
  font-size: 22px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.page-desc {
  color: #909399;
  margin-bottom: 24px;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  
  .stat-icon {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    
    &.users { background: #e6f7ff; color: #1890ff; }
    &.modules { background: #f6ffed; color: #52c41a; }
    &.requests { background: #fff7e6; color: #fa8c16; }
    &.errors { background: #fff1f0; color: #f5222d; }
  }
  
  .stat-info {
    .stat-value {
      font-size: 28px;
      font-weight: 600;
      color: #303133;
    }
    
    .stat-label {
      font-size: 14px;
      color: #909399;
    }
  }
}

.info-section {
  background: white;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  
  h3 {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid #ebeef5;
  }
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  
  label {
    font-size: 13px;
    color: #909399;
  }
  
  span {
    font-size: 14px;
    color: #303133;
    
    &.copyable {
      cursor: pointer;
      display: flex;
      align-items: center;
      gap: 4px;
      
      &:hover {
        color: #409eff;
      }
    }
  }
}

.module-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.config-form {
  background: white;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.form-section {
  margin-bottom: 32px;
  
  h4 {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 20px;
    padding-bottom: 12px;
    border-bottom: 1px solid #ebeef5;
  }
}

.form-hint {
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}

:deep(.el-form-item) {
  margin-bottom: 20px;
}

:deep(.el-input-number) {
  width: 150px;
}
</style>


.workspace-content {
  padding: 24px;
  
  .page-title {
    font-size: 24px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 20px;
  }
}
