<template>
  <div class="workspace">
    <!-- 工作台侧边栏 -->
    <div class="workspace-sidebar">
      <div class="sidebar-menu">
        <div 
          v-for="item in menuItems" 
          :key="item.key"
          class="menu-item"
          :class="{ active: currentMenu === item.key }"
          :data-testid="'menu-' + item.key"
          :data-menu-key="item.key"
          @click="currentMenu = item.key"
        >
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.label }}</span>
        </div>
      </div>
      <div class="sidebar-footer">
        <div class="menu-item back-item" @click="$router.push('/apps')">
          <el-icon><ArrowLeft /></el-icon>
          <span>返回APP列表</span>
        </div>
      </div>
    </div>

    <!-- 工作台内容区 -->
    <div class="workspace-content">
      <!-- 数据概览 -->
      <div v-if="currentMenu === 'overview'" class="content-section">
        <div class="section-header">
          <h2>数据概览</h2>
          <p>实时监控应用运行状态</p>
        </div>
        
        <!-- 统计卡片 -->
        <div class="stats-grid">
          <div class="stat-card blue">
            <div class="stat-content">
              <div class="stat-value">{{ stats.userCount.toLocaleString() }}</div>
              <div class="stat-label">总用户数</div>
              <div class="stat-trend up">
                <el-icon><Top /></el-icon>
                <span>+12.5%</span>
              </div>
            </div>
            <div class="stat-icon">
              <el-icon><User /></el-icon>
            </div>
          </div>
          
          <div class="stat-card green">
            <div class="stat-content">
              <div class="stat-value">{{ stats.activeUsers.toLocaleString() }}</div>
              <div class="stat-label">活跃用户</div>
              <div class="stat-trend up">
                <el-icon><Top /></el-icon>
                <span>+8.3%</span>
              </div>
            </div>
            <div class="stat-icon">
              <el-icon><UserFilled /></el-icon>
            </div>
          </div>
          
          <div class="stat-card orange">
            <div class="stat-content">
              <div class="stat-value">{{ stats.todayRequests.toLocaleString() }}</div>
              <div class="stat-label">今日请求</div>
              <div class="stat-trend up">
                <el-icon><Top /></el-icon>
                <span>+15.2%</span>
              </div>
            </div>
            <div class="stat-icon">
              <el-icon><DataLine /></el-icon>
            </div>
          </div>
          
          <div class="stat-card red">
            <div class="stat-content">
              <div class="stat-value">{{ stats.todayErrors }}</div>
              <div class="stat-label">今日异常</div>
              <div class="stat-trend down">
                <el-icon><Bottom /></el-icon>
                <span>-5.1%</span>
              </div>
            </div>
            <div class="stat-icon">
              <el-icon><Warning /></el-icon>
            </div>
          </div>
        </div>

        <!-- 图表区域 -->
        <div class="charts-row">
          <div class="chart-card">
            <div class="chart-header">
              <h3>请求趋势</h3>
              <el-radio-group v-model="chartPeriod" size="small">
                <el-radio-button label="7d">7天</el-radio-button>
                <el-radio-button label="30d">30天</el-radio-button>
              </el-radio-group>
            </div>
            <div class="chart-body" ref="requestChartRef"></div>
          </div>
          
          <div class="chart-card">
            <div class="chart-header">
              <h3>模块调用分布</h3>
            </div>
            <div class="chart-body" ref="moduleChartRef"></div>
          </div>
        </div>
      </div>

      <!-- 用户管理 -->
      <div v-else-if="currentMenu === 'users'" class="content-section">
        <div class="section-header">
          <h2>用户管理</h2>
          <p>管理应用用户数据</p>
        </div>
        
        <!-- 搜索和操作栏 -->
        <div class="toolbar">
          <div class="search-area">
            <el-input 
              v-model="userSearch" 
              placeholder="搜索用户名/手机号/邮箱" 
              prefix-icon="Search"
              clearable
              style="width: 300px"
              @keyup.enter="fetchUserList"
            />
            <el-select v-model="userStatus" placeholder="用户状态" clearable style="width: 120px" @change="fetchUserList">
              <el-option label="全部" value="" />
              <el-option label="正常" value="1" />
              <el-option label="禁用" value="0" />
            </el-select>
            <el-button type="primary" @click="fetchUserList">搜索</el-button>
          </div>
          <div class="action-area">
            <el-button type="primary" :icon="Plus">添加用户</el-button>
            <el-button :icon="Download">导出</el-button>
          </div>
        </div>

        <!-- 用户表格 -->
        <el-table :data="userList" stripe style="width: 100%" v-loading="userLoading">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="nickname" label="昵称" width="150" />
          <el-table-column prop="phone" label="手机号" width="140" />
          <el-table-column prop="email" label="邮箱" min-width="180" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                {{ row.status === 1 ? '正常' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="注册时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small">编辑</el-button>
              <el-button 
                :type="row.status === 1 ? 'danger' : 'success'" 
                link 
                size="small"
                @click="toggleUserStatus(row)"
              >
                {{ row.status === 1 ? '禁用' : '启用' }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination">
          <el-pagination
            v-model:current-page="userPage"
            v-model:page-size="userPageSize"
            :total="userTotal"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="fetchUserList"
            @current-change="fetchUserList"
          />
        </div>
      </div>

      <!-- 消息推送 -->
      <div v-else-if="currentMenu === 'messages'" class="content-section">
        <div class="section-header">
          <h2>消息推送</h2>
          <p>管理应用消息和推送通知</p>
        </div>

        <el-tabs v-model="messageTab">
          <el-tab-pane label="发送消息" name="send">
            <div class="message-form">
              <el-form :model="messageForm" label-width="100px" :rules="messageRules" ref="messageFormRef">
                <el-form-item label="推送类型" prop="type">
                  <el-radio-group v-model="messageForm.type">
                    <el-radio label="all">全部用户</el-radio>
                    <el-radio label="group">用户分组</el-radio>
                    <el-radio label="user">指定用户</el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-form-item v-if="messageForm.type === 'user'" label="用户ID" prop="userIds">
                  <el-input v-model="messageForm.userIds" placeholder="多个用户ID用逗号分隔" />
                </el-form-item>
                <el-form-item label="消息标题" prop="title">
                  <el-input v-model="messageForm.title" placeholder="请输入消息标题" />
                </el-form-item>
                <el-form-item label="消息内容" prop="content">
                  <el-input v-model="messageForm.content" type="textarea" :rows="4" placeholder="请输入消息内容" />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" @click="sendMessageNow" :loading="messageSending">立即推送</el-button>
                  <el-button>定时推送</el-button>
                </el-form-item>
              </el-form>
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="推送记录" name="history">
            <el-table :data="messageHistory" stripe v-loading="messageLoading">
              <el-table-column prop="title" label="标题" min-width="200" />
              <el-table-column prop="type" label="推送类型" width="120" />
              <el-table-column prop="content" label="内容" min-width="200" show-overflow-tooltip />
              <el-table-column prop="created_at" label="推送时间" width="180">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                    {{ row.status === 1 ? '已读' : '未读' }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- 日志查询 -->
      <div v-else-if="currentMenu === 'logs'" class="content-section">
        <div class="section-header">
          <h2>日志查询</h2>
          <p>查看应用运行日志和操作记录</p>
        </div>

        <!-- 日志统计卡片 -->
        <div class="log-stats">
          <div class="log-stat-item">
            <span class="label">总日志数</span>
            <span class="value">{{ logStats.total || 0 }}</span>
          </div>
          <div class="log-stat-item error">
            <span class="label">错误日志</span>
            <span class="value">{{ logStats.error_count || 0 }}</span>
          </div>
          <div class="log-stat-item warn">
            <span class="label">警告日志</span>
            <span class="value">{{ logStats.warn_count || 0 }}</span>
          </div>
          <div class="log-stat-item info">
            <span class="label">今日日志</span>
            <span class="value">{{ logStats.today_count || 0 }}</span>
          </div>
        </div>

        <!-- 日志筛选 -->
        <div class="toolbar">
          <div class="search-area">
            <el-date-picker
              v-model="logDateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              style="width: 260px"
              value-format="YYYY-MM-DD"
              @change="fetchLogList"
            />
            <el-select v-model="logLevel" placeholder="日志级别" clearable style="width: 120px" @change="fetchLogList">
              <el-option label="全部" value="" />
              <el-option label="DEBUG" value="debug">
                <el-tag type="info" size="small">DEBUG</el-tag>
              </el-option>
              <el-option label="INFO" value="info">
                <el-tag type="success" size="small">INFO</el-tag>
              </el-option>
              <el-option label="WARN" value="warn">
                <el-tag type="warning" size="small">WARN</el-tag>
              </el-option>
              <el-option label="ERROR" value="error">
                <el-tag type="danger" size="small">ERROR</el-tag>
              </el-option>
            </el-select>
            <el-input 
              v-model="logSearch" 
              placeholder="搜索日志内容" 
              prefix-icon="Search"
              clearable
              style="width: 200px"
              @keyup.enter="fetchLogList"
            />
            <el-button type="primary" @click="fetchLogList">搜索</el-button>
          </div>
          <div class="action-area">
            <el-button :icon="Refresh" @click="fetchLogList">刷新</el-button>
            <el-button :icon="Download" @click="exportLogs">导出</el-button>
          </div>
        </div>

        <!-- 日志列表 -->
        <div class="log-list" v-loading="logLoading">
          <div v-if="logList.length === 0" class="empty-logs">
            <el-empty description="暂无日志数据" />
          </div>
          <div v-else>
            <div v-for="log in logList" :key="log.id" class="log-item" :class="log.level">
              <div class="log-time">{{ formatDate(log.created_at) }}</div>
              <el-tag :type="getLogTagType(log.level)" size="small">{{ (log.level || 'info').toUpperCase() }}</el-tag>
              <div class="log-module" v-if="log.module">{{ log.module }}</div>
              <div class="log-content">{{ log.message }}</div>
            </div>
          </div>
        </div>

        <!-- 日志分页 -->
        <div class="pagination">
          <el-pagination
            v-model:current-page="logPage"
            v-model:page-size="logPageSize"
            :total="logTotal"
            :page-sizes="[20, 50, 100, 200]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="fetchLogList"
            @current-change="fetchLogList"
          />
        </div>
      </div>

      <!-- 版本管理 -->
      <div v-else-if="currentMenu === 'versions'" class="content-section">
        <div class="section-header">
          <h2>版本管理</h2>
          <p>管理应用版本和更新</p>
        </div>

        <div class="toolbar">
          <div class="search-area">
            <el-select v-model="versionPlatform" placeholder="选择平台" clearable style="width: 120px" @change="fetchVersionList">
              <el-option label="全部" value="" />
              <el-option label="Android" value="android" />
              <el-option label="iOS" value="ios" />
            </el-select>
            <el-select v-model="versionStatus" placeholder="版本状态" clearable style="width: 120px" @change="fetchVersionList">
              <el-option label="全部" value="" />
              <el-option label="已发布" value="published" />
              <el-option label="草稿" value="draft" />
              <el-option label="已下线" value="offline" />
            </el-select>
          </div>
          <div class="action-area">
            <el-button type="primary" :icon="Plus" @click="showVersionDialog = true">发布新版本</el-button>
          </div>
        </div>

        <el-table :data="versionList" stripe v-loading="versionLoading">
          <el-table-column prop="version" label="版本号" width="120" />
          <el-table-column prop="platform" label="平台" width="100">
            <template #default="{ row }">
              <el-tag :type="row.platform === 'android' ? 'success' : 'primary'" size="small">
                {{ row.platform === 'android' ? 'Android' : 'iOS' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="更新说明" min-width="250" show-overflow-tooltip />
          <el-table-column prop="download_url" label="下载地址" min-width="200" show-overflow-tooltip />
          <el-table-column prop="force_update" label="强制更新" width="100">
            <template #default="{ row }">
              <el-tag :type="row.force_update ? 'danger' : 'info'" size="small">
                {{ row.force_update ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getVersionStatusType(row.status)" size="small">
                {{ getVersionStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="发布时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="editVersion(row)">编辑</el-button>
              <el-button 
                v-if="row.status !== 'published'" 
                type="success" 
                link 
                size="small"
                @click="publishVersionAction(row)"
              >
                发布
              </el-button>
              <el-button 
                v-if="row.status === 'published'" 
                type="danger" 
                link 
                size="small"
                @click="offlineVersionAction(row)"
              >
                下线
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 版本分页 -->
        <div class="pagination">
          <el-pagination
            v-model:current-page="versionPage"
            v-model:page-size="versionPageSize"
            :total="versionTotal"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="fetchVersionList"
            @current-change="fetchVersionList"
          />
        </div>
      </div>

      <!-- 存储服务 -->
      <div v-else-if="currentMenu === 'storage'" class="content-section">
        <div class="section-header">
          <h2>存储服务</h2>
          <p>管理应用文件存储</p>
        </div>

        <!-- 存储统计 -->
        <div class="stats-grid small">
          <div class="stat-card blue">
            <div class="stat-content">
              <div class="stat-value">{{ fileStats.total_count || 0 }}</div>
              <div class="stat-label">文件总数</div>
            </div>
            <div class="stat-icon"><el-icon><FolderOpened /></el-icon></div>
          </div>
          <div class="stat-card green">
            <div class="stat-content">
              <div class="stat-value">{{ formatFileSize(fileStats.total_size || 0) }}</div>
              <div class="stat-label">占用空间</div>
            </div>
            <div class="stat-icon"><el-icon><PieChart /></el-icon></div>
          </div>
          <div class="stat-card orange">
            <div class="stat-content">
              <div class="stat-value">{{ fileStats.today_count || 0 }}</div>
              <div class="stat-label">今日上传</div>
            </div>
            <div class="stat-icon"><el-icon><Upload /></el-icon></div>
          </div>
        </div>

        <!-- 操作栏 -->
        <div class="toolbar">
          <div class="search-area">
            <el-input v-model="fileSearch" placeholder="搜索文件名" prefix-icon="Search" clearable style="width: 200px" @keyup.enter="fetchFileList" />
            <el-select v-model="fileType" placeholder="文件类型" clearable style="width: 120px" @change="fetchFileList">
              <el-option label="全部" value="" />
              <el-option label="图片" value="image" />
              <el-option label="文档" value="document" />
              <el-option label="视频" value="video" />
              <el-option label="音频" value="audio" />
              <el-option label="其他" value="other" />
            </el-select>
            <el-button type="primary" @click="fetchFileList">搜索</el-button>
          </div>
          <div class="action-area">
            <el-upload
              :action="uploadUrl"
              :headers="uploadHeaders"
              :on-success="handleUploadSuccess"
              :on-error="handleUploadError"
              :show-file-list="false"
              multiple
            >
              <el-button type="primary" :icon="Upload">上传文件</el-button>
            </el-upload>
            <el-button :icon="Delete" :disabled="selectedFiles.length === 0" @click="batchDeleteFilesAction">批量删除</el-button>
          </div>
        </div>

        <!-- 文件列表 -->
        <el-table :data="fileList" stripe v-loading="fileLoading" @selection-change="handleFileSelectionChange">
          <el-table-column type="selection" width="50" />
          <el-table-column prop="original_name" label="文件名" min-width="200" show-overflow-tooltip />
          <el-table-column prop="mime_type" label="类型" width="120">
            <template #default="{ row }">
              <el-tag size="small">{{ getFileTypeLabel(row.mime_type) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="file_size" label="大小" width="100">
            <template #default="{ row }">{{ formatFileSize(row.file_size) }}</template>
          </el-table-column>
          <el-table-column prop="created_at" label="上传时间" width="180">
            <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="downloadFileAction(row)">下载</el-button>
              <el-button type="danger" link size="small" @click="deleteFileAction(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination">
          <el-pagination
            v-model:current-page="filePage"
            v-model:page-size="filePageSize"
            :total="fileTotal"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="fetchFileList"
            @current-change="fetchFileList"
          />
        </div>
      </div>

      <!-- 数据埋点 -->
      <div v-else-if="currentMenu === 'events'" class="content-section">
        <div class="section-header">
          <h2>数据埋点</h2>
          <p>查看和分析应用事件数据</p>
        </div>

        <!-- 埋点统计 -->
        <div class="stats-grid small">
          <div class="stat-card blue">
            <div class="stat-content">
              <div class="stat-value">{{ eventStats.total_events || 0 }}</div>
              <div class="stat-label">事件总数</div>
            </div>
            <div class="stat-icon"><el-icon><DataAnalysis /></el-icon></div>
          </div>
          <div class="stat-card green">
            <div class="stat-content">
              <div class="stat-value">{{ eventStats.today_events || 0 }}</div>
              <div class="stat-label">今日事件</div>
            </div>
            <div class="stat-icon"><el-icon><TrendCharts /></el-icon></div>
          </div>
          <div class="stat-card orange">
            <div class="stat-content">
              <div class="stat-value">{{ eventStats.unique_users || 0 }}</div>
              <div class="stat-label">独立用户</div>
            </div>
            <div class="stat-icon"><el-icon><User /></el-icon></div>
          </div>
        </div>

        <el-tabs v-model="eventTab">
          <el-tab-pane label="事件列表" name="list">
            <div class="toolbar">
              <div class="search-area">
                <el-date-picker
                  v-model="eventDateRange"
                  type="daterange"
                  range-separator="至"
                  start-placeholder="开始日期"
                  end-placeholder="结束日期"
                  style="width: 260px"
                  value-format="YYYY-MM-DD"
                  @change="fetchEventList"
                />
                <el-select v-model="eventName" placeholder="事件类型" clearable style="width: 150px" @change="fetchEventList">
                  <el-option label="全部" value="" />
                  <el-option label="页面浏览" value="page_view" />
                  <el-option label="按钮点击" value="button_click" />
                  <el-option label="表单提交" value="form_submit" />
                  <el-option label="登录" value="login" />
                  <el-option label="注册" value="register" />
                </el-select>
                <el-button type="primary" @click="fetchEventList">搜索</el-button>
              </div>
              <div class="action-area">
                <el-button :icon="Download" @click="exportEvents">导出</el-button>
              </div>
            </div>

            <el-table :data="eventList" stripe v-loading="eventLoading">
              <el-table-column prop="event_name" label="事件名称" width="150" />
              <el-table-column prop="event_type" label="事件类型" width="120">
                <template #default="{ row }">
                  <el-tag size="small">{{ row.event_type || '自定义' }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="user_id" label="用户ID" width="120" />
              <el-table-column prop="device_id" label="设备ID" width="150" show-overflow-tooltip />
              <el-table-column prop="properties" label="属性" min-width="200">
                <template #default="{ row }">
                  <el-tooltip :content="JSON.stringify(row.properties)" placement="top">
                    <span class="properties-preview">{{ JSON.stringify(row.properties).slice(0, 50) }}...</span>
                  </el-tooltip>
                </template>
              </el-table-column>
              <el-table-column prop="created_at" label="时间" width="180">
                <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
              </el-table-column>
            </el-table>

            <div class="pagination">
              <el-pagination
                v-model:current-page="eventPage"
                v-model:page-size="eventPageSize"
                :total="eventTotal"
                :page-sizes="[20, 50, 100, 200]"
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="fetchEventList"
                @current-change="fetchEventList"
              />
            </div>
          </el-tab-pane>

          <el-tab-pane label="漏斗分析" name="funnel">
            <div class="funnel-config">
              <el-form :inline="true">
                <el-form-item label="漏斗步骤">
                  <el-select v-model="funnelSteps" multiple placeholder="选择事件步骤" style="width: 400px">
                    <el-option label="页面浏览" value="page_view" />
                    <el-option label="按钮点击" value="button_click" />
                    <el-option label="表单提交" value="form_submit" />
                    <el-option label="登录" value="login" />
                    <el-option label="注册" value="register" />
                    <el-option label="支付" value="payment" />
                  </el-select>
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" @click="analyzeFunnel">分析</el-button>
                </el-form-item>
              </el-form>
            </div>
            <div class="funnel-chart" ref="funnelChartRef"></div>
          </el-tab-pane>

          <el-tab-pane label="事件定义" name="definitions">
            <div class="toolbar">
              <div class="search-area"></div>
              <div class="action-area">
                <el-button type="primary" :icon="Plus" @click="showEventDefDialog = true">新建事件</el-button>
              </div>
            </div>

            <el-table :data="eventDefinitions" stripe v-loading="eventDefLoading">
              <el-table-column prop="name" label="事件名称" width="150" />
              <el-table-column prop="code" label="事件编码" width="150" />
              <el-table-column prop="description" label="描述" min-width="200" />
              <el-table-column prop="properties" label="属性定义" min-width="200">
                <template #default="{ row }">
                  <el-tag v-for="prop in (row.properties || []).slice(0, 3)" :key="prop" size="small" style="margin-right: 4px">{{ prop }}</el-tag>
                  <span v-if="(row.properties || []).length > 3">...</span>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="120" fixed="right">
                <template #default="{ row }">
                  <el-button type="danger" link size="small" @click="deleteEventDefAction(row)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- 监控告警 -->
      <div v-else-if="currentMenu === 'monitor'" class="content-section">
        <div class="section-header">
          <div class="header-left">
            <h2>监控告警</h2>
            <p>实时监控应用运行状态</p>
          </div>
          <div class="header-right">
            <el-tag :type="wsConnected ? 'success' : 'danger'" size="small">
              <el-icon style="margin-right: 4px;"><Connection /></el-icon>
              {{ wsConnected ? 'WebSocket已连接' : 'WebSocket未连接' }}
            </el-tag>
          </div>
        </div>

        <!-- 监控指标卡片 -->
        <div class="stats-grid">
          <div class="stat-card" :class="healthStatus === 'healthy' ? 'green' : 'red'">
            <div class="stat-content">
              <div class="stat-value">{{ healthStatus === 'healthy' ? '正常' : '异常' }}</div>
              <div class="stat-label">系统状态</div>
            </div>
            <div class="stat-icon"><el-icon><Connection /></el-icon></div>
          </div>
          <div class="stat-card blue">
            <div class="stat-content">
              <div class="stat-value">{{ monitorStats.cpu_usage || 0 }}%</div>
              <div class="stat-label">CPU使用率</div>
            </div>
            <div class="stat-icon"><el-icon><Monitor /></el-icon></div>
          </div>
          <div class="stat-card orange">
            <div class="stat-content">
              <div class="stat-value">{{ monitorStats.memory_usage || 0 }}%</div>
              <div class="stat-label">内存使用率</div>
            </div>
            <div class="stat-icon"><el-icon><PieChart /></el-icon></div>
          </div>
          <div class="stat-card red">
            <div class="stat-content">
              <div class="stat-value">{{ monitorStats.active_alerts || 0 }}</div>
              <div class="stat-label">活跃告警</div>
            </div>
            <div class="stat-icon"><el-icon><Warning /></el-icon></div>
          </div>
        </div>

        <el-tabs v-model="monitorTab">
          <el-tab-pane label="监控图表" name="charts">
            <div class="charts-row">
              <div class="chart-card">
                <div class="chart-header">
                  <h3>请求量趋势</h3>
                  <el-radio-group v-model="monitorPeriod" size="small" @change="fetchMonitorMetrics">
                    <el-radio-button label="1h">1小时</el-radio-button>
                    <el-radio-button label="24h">24小时</el-radio-button>
                    <el-radio-button label="7d">7天</el-radio-button>
                  </el-radio-group>
                </div>
                <div class="chart-body" ref="requestMetricChartRef"></div>
              </div>
              <div class="chart-card">
                <div class="chart-header">
                  <h3>错误率趋势</h3>
                </div>
                <div class="chart-body" ref="errorMetricChartRef"></div>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="告警列表" name="alerts">
            <div class="toolbar">
              <div class="search-area">
                <el-select v-model="alertStatus" placeholder="告警状态" clearable style="width: 120px" @change="fetchAlertList">
                  <el-option label="全部" value="" />
                  <el-option label="未处理" value="0" />
                  <el-option label="已处理" value="1" />
                </el-select>
              </div>
              <div class="action-area">
                <el-button :icon="Refresh" @click="fetchAlertList">刷新</el-button>
              </div>
            </div>

            <el-table :data="alertList" stripe v-loading="alertLoading">
              <el-table-column prop="name" label="告警名称" width="150" />
              <el-table-column prop="metric_name" label="监控指标" width="150" />
              <el-table-column prop="metric_value" label="当前值" width="100" />
              <el-table-column prop="threshold" label="阈值" width="100" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.status === 0 ? 'danger' : 'success'" size="small">
                    {{ row.status === 0 ? '未处理' : '已处理' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="created_at" label="触发时间" width="180">
                <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
              </el-table-column>
              <el-table-column label="操作" width="120" fixed="right">
                <template #default="{ row }">
                  <el-button v-if="row.status === 0" type="primary" link size="small" @click="resolveAlert(row)">处理</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="告警规则" name="rules">
            <div class="toolbar">
              <div class="search-area"></div>
              <div class="action-area">
                <el-button type="primary" :icon="Plus" @click="showAlertRuleDialog = true">新建规则</el-button>
              </div>
            </div>

            <el-table :data="alertRules" stripe v-loading="alertRuleLoading">
              <el-table-column prop="name" label="规则名称" width="150" />
              <el-table-column prop="metric_name" label="监控指标" width="150" />
              <el-table-column prop="condition_type" label="条件" width="100">
                <template #default="{ row }">
                  {{ row.condition_type === 'gt' ? '>' : row.condition_type === 'lt' ? '<' : row.condition_type }}
                </template>
              </el-table-column>
              <el-table-column prop="threshold" label="阈值" width="100" />
              <el-table-column prop="duration" label="持续时间(秒)" width="120" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-switch v-model="row.status" :active-value="1" :inactive-value="0" @change="toggleAlertRule(row)" />
                </template>
              </el-table-column>
              <el-table-column label="操作" width="120" fixed="right">
                <template #default="{ row }">
                  <el-button type="danger" link size="small" @click="deleteAlertRule(row)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- 审计日志 -->
      <div v-else-if="currentMenu === 'audit'" class="content-section">
        <div class="section-header">
          <h2>审计日志</h2>
          <p>查看系统操作记录和安全审计</p>
        </div>

        <!-- 审计统计卡片 -->
        <div class="stats-grid">
          <div class="stat-card blue">
            <div class="stat-content">
              <div class="stat-value">{{ auditStats.total_count?.toLocaleString() || 0 }}</div>
              <div class="stat-label">总操作数</div>
            </div>
            <div class="stat-icon"><el-icon><DataLine /></el-icon></div>
          </div>
          <div class="stat-card green">
            <div class="stat-content">
              <div class="stat-value">{{ auditStats.user_stats?.length || 0 }}</div>
              <div class="stat-label">活跃用户</div>
            </div>
            <div class="stat-icon"><el-icon><User /></el-icon></div>
          </div>
          <div class="stat-card orange">
            <div class="stat-content">
              <div class="stat-value">{{ auditStats.action_stats?.length || 0 }}</div>
              <div class="stat-label">操作类型</div>
            </div>
            <div class="stat-icon"><el-icon><Management /></el-icon></div>
          </div>
          <div class="stat-card purple">
            <div class="stat-content">
              <div class="stat-value">{{ auditStats.resource_stats?.length || 0 }}</div>
              <div class="stat-label">资源类型</div>
            </div>
            <div class="stat-icon"><el-icon><FolderOpened /></el-icon></div>
          </div>
        </div>

        <!-- 搜索和筛选 -->
        <div class="filter-bar">
          <el-input v-model="auditSearch.keyword" placeholder="搜索描述、用户名、IP" style="width: 200px" clearable>
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
          <el-select v-model="auditSearch.action" placeholder="操作类型" clearable style="width: 120px">
            <el-option label="查看" value="view" />
            <el-option label="创建" value="create" />
            <el-option label="更新" value="update" />
            <el-option label="删除" value="delete" />
            <el-option label="登录" value="login" />
            <el-option label="登出" value="logout" />
            <el-option label="导出" value="export" />
          </el-select>
          <el-select v-model="auditSearch.resource" placeholder="资源类型" clearable style="width: 120px">
            <el-option label="用户" value="user" />
            <el-option label="应用" value="app" />
            <el-option label="配置" value="config" />
            <el-option label="消息" value="message" />
            <el-option label="推送" value="push" />
            <el-option label="文件" value="file" />
            <el-option label="版本" value="version" />
          </el-select>
          <el-date-picker v-model="auditSearch.dateRange" type="daterange" range-separator="-" start-placeholder="开始日期" end-placeholder="结束日期" style="width: 240px" />
          <el-button type="primary" @click="fetchAuditLogs"><el-icon><Search /></el-icon>查询</el-button>
          <el-button @click="exportAuditLogsData"><el-icon><Download /></el-icon>导出</el-button>
        </div>

        <!-- 审计日志列表 -->
        <el-table :data="auditLogs" v-loading="auditLoading" stripe style="margin-top: 16px">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="created_at" label="时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column prop="user_name" label="用户" width="120" />
          <el-table-column prop="action" label="操作" width="100">
            <template #default="{ row }">
              <el-tag :type="getActionTagType(row.action)" size="small">{{ getActionLabel(row.action) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="resource" label="资源" width="100">
            <template #default="{ row }">
              {{ getResourceLabel(row.resource) }}
            </template>
          </el-table-column>
          <el-table-column prop="description" label="描述" min-width="150" />
          <el-table-column prop="ip_address" label="IP地址" width="140" />
          <el-table-column prop="status_code" label="状态码" width="80">
            <template #default="{ row }">
              <el-tag :type="row.status_code < 400 ? 'success' : 'danger'" size="small">{{ row.status_code }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="duration" label="耗时(ms)" width="100" />
        </el-table>

        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="auditPage"
            v-model:page-size="auditPageSize"
            :total="auditTotal"
            :page-sizes="[20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            @size-change="fetchAuditLogs"
            @current-change="fetchAuditLogs"
          />
        </div>
      </div>
    </div>

    <!-- 发布新版本对话框 -->
    <el-dialog v-model="showVersionDialog" :title="editingVersion ? '编辑版本' : '发布新版本'" width="600px">
      <el-form :model="versionForm" label-width="100px" :rules="versionRules" ref="versionFormRef">
        <el-form-item label="版本号" prop="version">
          <el-input v-model="versionForm.version" placeholder="如：1.0.0" />
        </el-form-item>
        <el-form-item label="平台" prop="platform">
          <el-radio-group v-model="versionForm.platform">
            <el-radio label="android">Android</el-radio>
            <el-radio label="ios">iOS</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="下载地址" prop="download_url">
          <el-input v-model="versionForm.download_url" placeholder="请输入下载地址" />
        </el-form-item>
        <el-form-item label="更新说明" prop="description">
          <el-input v-model="versionForm.description" type="textarea" :rows="4" placeholder="请输入更新说明" />
        </el-form-item>
        <el-form-item label="强制更新">
          <el-switch v-model="versionForm.force_update" />
          <span class="form-tip">开启后，用户必须更新到此版本才能使用</span>
        </el-form-item>
        <el-form-item label="灰度发布">
          <el-switch v-model="versionForm.gray_release" />
        </el-form-item>
        <el-form-item v-if="versionForm.gray_release" label="灰度比例">
          <el-slider v-model="versionForm.gray_percent" :min="1" :max="100" :format-tooltip="val => `${val}%`" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showVersionDialog = false">取消</el-button>
        <el-button type="primary" @click="submitVersion" :loading="versionSubmitting">
          {{ editingVersion ? '保存' : '发布' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 新建事件定义对话框 -->
    <el-dialog v-model="showEventDefDialog" title="新建事件定义" width="600px">
      <el-form :model="eventDefForm" label-width="100px" :rules="eventDefRules" ref="eventDefFormRef">
        <el-form-item label="事件名称" prop="name">
          <el-input v-model="eventDefForm.name" placeholder="请输入事件名称，如：用户登录" />
        </el-form-item>
        <el-form-item label="事件编码" prop="code">
          <el-input v-model="eventDefForm.code" placeholder="请输入事件编码，如：user_login" />
        </el-form-item>
        <el-form-item label="事件描述" prop="description">
          <el-input v-model="eventDefForm.description" type="textarea" :rows="3" placeholder="请输入事件描述" />
        </el-form-item>
        <el-form-item label="事件属性">
          <div class="property-list">
            <div v-for="(prop, index) in eventDefForm.properties" :key="index" class="property-item">
              <el-input v-model="eventDefForm.properties[index]" placeholder="属性名称" style="width: 200px" />
              <el-button type="danger" :icon="Delete" circle size="small" @click="eventDefForm.properties.splice(index, 1)" />
            </div>
            <el-button type="primary" link @click="eventDefForm.properties.push('')">+ 添加属性</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEventDefDialog = false">取消</el-button>
        <el-button type="primary" @click="submitEventDef" :loading="eventDefSubmitting">创建</el-button>
      </template>
    </el-dialog>

    <!-- 新建告警规则对话框 -->
    <el-dialog v-model="showAlertRuleDialog" title="新建告警规则" width="600px">
      <el-form :model="alertRuleForm" label-width="100px" :rules="alertRuleRules" ref="alertRuleFormRef">
        <el-form-item label="规则名称" prop="name">
          <el-input v-model="alertRuleForm.name" placeholder="请输入规则名称" />
        </el-form-item>
        <el-form-item label="监控指标" prop="metric_name">
          <el-select v-model="alertRuleForm.metric_name" placeholder="请选择监控指标" style="width: 100%">
            <el-option label="CPU使用率" value="cpu_usage" />
            <el-option label="内存使用率" value="memory_usage" />
            <el-option label="请求延迟" value="request_latency" />
            <el-option label="错误率" value="error_rate" />
            <el-option label="QPS" value="qps" />
          </el-select>
        </el-form-item>
        <el-form-item label="触发条件" prop="condition_type">
          <el-select v-model="alertRuleForm.condition_type" placeholder="请选择条件" style="width: 120px">
            <el-option label="大于" value="gt" />
            <el-option label="小于" value="lt" />
            <el-option label="等于" value="eq" />
            <el-option label="大于等于" value="gte" />
            <el-option label="小于等于" value="lte" />
          </el-select>
        </el-form-item>
        <el-form-item label="阈值" prop="threshold">
          <el-input-number v-model="alertRuleForm.threshold" :min="0" :max="100" placeholder="请输入阈值" />
          <span style="margin-left: 8px; color: #909399;">{{ alertRuleForm.metric_name?.includes('rate') || alertRuleForm.metric_name?.includes('usage') ? '%' : '' }}</span>
        </el-form-item>
        <el-form-item label="持续时间" prop="duration">
          <el-input-number v-model="alertRuleForm.duration" :min="0" :max="3600" placeholder="持续时间" />
          <span style="margin-left: 8px; color: #909399;">秒</span>
        </el-form-item>
        <el-form-item label="告警级别" prop="level">
          <el-radio-group v-model="alertRuleForm.level">
            <el-radio label="info">提示</el-radio>
            <el-radio label="warning">警告</el-radio>
            <el-radio label="critical">严重</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="通知方式">
          <el-checkbox-group v-model="alertRuleForm.notify_channels">
            <el-checkbox label="email">邮件</el-checkbox>
            <el-checkbox label="sms">短信</el-checkbox>
            <el-checkbox label="webhook">Webhook</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAlertRuleDialog = false">取消</el-button>
        <el-button type="primary" @click="submitAlertRule" :loading="alertRuleSubmitting">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import wsClient from '@/utils/websocket'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as echarts from 'echarts'
import { 
  DataLine, User, UserFilled, Warning, Top, Bottom,
  Plus, Download, Refresh, Search, Delete, Upload, View, Edit,
  House, Management, Bell, Document, Promotion, ArrowLeft,
  FolderOpened, DataAnalysis, Monitor, PieChart, TrendCharts, Timer, Connection, Lock
} from '@element-plus/icons-vue'
import {
  getUserList, getUserStats, updateUserStatus,
  getLogList, getLogStats, exportLogs as exportLogsApi,
  getMessageList, sendMessage, batchSendMessage,
  getVersionList, createVersion, publishVersion, offlineVersion,
  // 存储服务
  uploadFile, getFileList, downloadFile, deleteFile, getFileStats, batchDeleteFiles,
  // 数据埋点
  getEventList, getEventStats, getFunnelAnalysis, getEventDefinitions, createEventDefinition, deleteEventDefinition,
  // 监控告警
  getMonitorMetrics, getAlertList, createAlert, updateAlert, deleteAlert, getMonitorStats, getHealthCheck,
  // 审计日志
  getAuditLogs, getAuditStats, exportAuditLogs
} from '@/api/app'

const props = defineProps({
  appId: {
    type: String,
    default: ''
  },
  appInfo: {
    type: Object,
    default: () => ({})
  },
  initialMenu: {
    type: String,
    default: 'overview'
  }
})

// 菜单配置
const menuItems = [
  { key: 'overview', label: '数据概览', icon: House },
  { key: 'users', label: '用户管理', icon: User },
  { key: 'messages', label: '消息推送', icon: Bell },
  { key: 'storage', label: '存储服务', icon: FolderOpened },
  { key: 'events', label: '数据埋点', icon: DataAnalysis },
  { key: 'monitor', label: '监控告警', icon: Monitor },
  { key: 'logs', label: '日志查询', icon: Document },
  { key: 'versions', label: '版本管理', icon: Promotion },
  { key: 'audit', label: '审计日志', icon: Lock }
]

const currentMenu = ref(props.initialMenu || 'overview')
const chartPeriod = ref('7d')

// 统计数据
const stats = ref({
  userCount: 0,
  activeUsers: 0,
  todayRequests: 0,
  todayErrors: 0
})

// 图表引用
const requestChartRef = ref(null)
const moduleChartRef = ref(null)
let requestChart = null
let moduleChart = null

// 用户管理
const userSearch = ref('')
const userStatus = ref('')
const userPage = ref(1)
const userPageSize = ref(10)
const userTotal = ref(0)
const userList = ref([])
const userLoading = ref(false)

// 消息推送
const messageTab = ref('send')
const messageForm = ref({
  type: 'all',
  userIds: '',
  title: '',
  content: ''
})
const messageRules = {
  title: [{ required: true, message: '请输入消息标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入消息内容', trigger: 'blur' }]
}
const messageFormRef = ref(null)
const messageSending = ref(false)
const messageHistory = ref([])
const messageLoading = ref(false)

// 日志查询
const logDateRange = ref([])
const logLevel = ref('')
const logSearch = ref('')
const logList = ref([])
const logLoading = ref(false)
const logPage = ref(1)
const logPageSize = ref(20)
const logTotal = ref(0)
const logStats = ref({})

// 版本管理
const versionList = ref([])
const versionLoading = ref(false)
const versionPage = ref(1)
const versionPageSize = ref(10)
const versionTotal = ref(0)
const versionPlatform = ref('')
const versionStatus = ref('')
const showVersionDialog = ref(false)
const editingVersion = ref(null)
const versionForm = ref({
  version: '',
  platform: 'android',
  download_url: '',
  description: '',
  force_update: false,
  gray_release: false,
  gray_percent: 10
})
const versionRules = {
  version: [{ required: true, message: '请输入版本号', trigger: 'blur' }],
  platform: [{ required: true, message: '请选择平台', trigger: 'change' }],
  download_url: [{ required: true, message: '请输入下载地址', trigger: 'blur' }],
  description: [{ required: true, message: '请输入更新说明', trigger: 'blur' }]
}
const versionFormRef = ref(null)
const versionSubmitting = ref(false)

// 存储服务
const fileList = ref([])
const fileLoading = ref(false)
const filePage = ref(1)
const filePageSize = ref(10)
const fileTotal = ref(0)
const fileSearch = ref('')
const fileType = ref('')
const fileStats = ref({ total_count: 0, total_size: 0, today_count: 0 })
const selectedFiles = ref([])
const uploadUrl = '/api/v1/files'
const uploadHeaders = { Authorization: `Bearer ${localStorage.getItem('token') || ''}` }

// 数据埋点
const eventTab = ref('list')
const eventList = ref([])
const eventLoading = ref(false)
const eventPage = ref(1)
const eventPageSize = ref(20)
const eventTotal = ref(0)
const eventDateRange = ref([])
const eventName = ref('')
const eventStats = ref({ total_events: 0, today_events: 0, unique_users: 0 })
const funnelSteps = ref([])
const funnelChartRef = ref(null)
let funnelChart = null
const eventDefinitions = ref([])
const eventDefLoading = ref(false)
const showEventDefDialog = ref(false)

// 监控告警
const monitorTab = ref('charts')
const monitorPeriod = ref('24h')
const monitorStats = ref({ cpu_usage: 0, memory_usage: 0, active_alerts: 0 })
const healthStatus = ref('healthy')
const requestMetricChartRef = ref(null)
const errorMetricChartRef = ref(null)
let requestMetricChart = null
let errorMetricChart = null
const alertList = ref([])
const alertLoading = ref(false)
const alertStatus = ref('')
const alertRules = ref([])
const alertRuleLoading = ref(false)

// 审计日志相关变量
const auditLogs = ref([])
const auditLoading = ref(false)
const auditPage = ref(1)
const auditPageSize = ref(20)
const auditTotal = ref(0)
const auditStats = ref({})
const auditSearch = ref({
  keyword: '',
  action: '',
  resource: '',
  dateRange: null
})
const showAlertRuleDialog = ref(false)
const alertRuleForm = ref({
  name: '',
  metric_name: '',
  condition_type: 'gt',
  threshold: 80,
  duration: 60,
  level: 'warning',
  notify_channels: []
})
const alertRuleRules = {
  name: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
  metric_name: [{ required: true, message: '请选择监控指标', trigger: 'change' }],
  condition_type: [{ required: true, message: '请选择触发条件', trigger: 'change' }],
  threshold: [{ required: true, message: '请输入阈值', trigger: 'blur' }]
}
const alertRuleFormRef = ref(null)
const alertRuleSubmitting = ref(false)

// 事件定义表单
const eventDefForm = ref({
  name: '',
  code: '',
  description: '',
  properties: []
})
const eventDefRules = {
  name: [{ required: true, message: '请输入事件名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入事件编码', trigger: 'blur' }]
}
const eventDefFormRef = ref(null)
const eventDefSubmitting = ref(false)

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', { 
    year: 'numeric', 
    month: '2-digit', 
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const getLogTagType = (level) => {
  const types = { debug: 'info', info: 'success', warn: 'warning', error: 'danger' }
  return types[level] || 'info'
}

const getVersionStatusType = (status) => {
  const types = { published: 'success', draft: 'info', offline: 'danger' }
  return types[status] || 'info'
}

const getVersionStatusText = (status) => {
  const texts = { published: '已发布', draft: '草稿', offline: '已下线' }
  return texts[status] || status
}

// 获取用户列表
const fetchUserList = async () => {
  if (!props.appId) {
    console.warn('[UserManagement] appId is not set, skipping fetch')
    return
  }
  userLoading.value = true
  console.log('[UserManagement] Fetching user list for app:', props.appId, {
    page: userPage.value,
    size: userPageSize.value,
    status: userStatus.value,
    search: userSearch.value
  })
  try {
    const res = await getUserList({
      app_id: props.appId,
      page: userPage.value,
      size: userPageSize.value,
      status: userStatus.value,
      search: userSearch.value
    })
    console.log('[UserManagement] API response:', res)
    // request.js已解包，res直接是数据对象
    userList.value = res.list || []
    userTotal.value = res.total || 0
    console.log('[UserManagement] Loaded', userList.value.length, 'users, total:', userTotal.value)
  } catch (error) {
    console.error('[UserManagement] Failed to fetch user list:', error)
    console.error('[UserManagement] Error details:', {
      message: error.message,
      response: error.response?.data,
      status: error.response?.status
    })
    // 不显示错误提示，因为 request.js 已经处理了
  } finally {
    userLoading.value = false
  }
}

// 切换用户状态
const toggleUserStatus = async (row) => {
  const newStatus = row.status === 1 ? 0 : 1
  const action = newStatus === 0 ? '禁用' : '启用'
  try {
    await ElMessageBox.confirm(`确定要${action}该用户吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await updateUserStatus(row.id, newStatus)
    // request.js已解包，成功时不会抛出异常
    ElMessage.success(`${action}成功`)
    fetchUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('更新用户状态失败:', error)
    }
  }
}

// 获取用户统计
const fetchUserStats = async () => {
  if (!props.appId || props.appId === '') return
  try {
    const res = await getUserStats(props.appId)
    // request.js已解包，res直接是数据对象
    stats.value.userCount = res.total || 0
    stats.value.activeUsers = res.active || 0
  } catch (error) {
    console.error('获取用户统计失败:', error)
  }
}

// 获取日志列表
const fetchLogList = async () => {
  if (!props.appId || props.appId === '') return
  logLoading.value = true
  try {
    const params = {
      app_id: props.appId,
      page: logPage.value,
      size: logPageSize.value
    }
    if (logLevel.value) params.level = logLevel.value
    if (logSearch.value) params.keyword = logSearch.value
    if (logDateRange.value && logDateRange.value.length === 2) {
      params.start_time = logDateRange.value[0]
      params.end_time = logDateRange.value[1]
    }
    
    const res = await getLogList(params)
    // request.js已解包，res直接是数据对象
    logList.value = res.list || []
    logTotal.value = res.total || 0
  } catch (error) {
    console.error('获取日志列表失败:', error)
  } finally {
    logLoading.value = false
  }
}

// 获取日志统计
const fetchLogStats = async () => {
  if (!props.appId || props.appId === '') return
  try {
    const res = await getLogStats(props.appId)
    // request.js已解包，res直接是数据对象
    logStats.value = res || {}
    stats.value.todayErrors = res.error_count || 0
  } catch (error) {
    console.error('获取日志统计失败:', error)
  }
}

// 导出日志
const exportLogs = async () => {
  if (!props.appId || props.appId === '') return
  try {
    const params = {
      app_id: props.appId,
      level: logLevel.value
    }
    if (logDateRange.value && logDateRange.value.length === 2) {
      params.start_time = logDateRange.value[0]
      params.end_time = logDateRange.value[1]
    }
    const res = await exportLogsApi(params)
    // request.js已解包，res直接是数据对象
    ElMessage.success(`导出成功，共${res.count || 0}条日志`)
    // 实际项目中这里应该下载文件
  } catch (error) {
    console.error('导出日志失败:', error)
  }
}

// 获取消息列表
const fetchMessageList = async () => {
  if (!props.appId || props.appId === '') return
  messageLoading.value = true
  try {
const res = await getMessageList({
      app_id: props.appId,
      page: 1,
      size: 50
    })
    // request.js已解包，res直接是数据对象
    messageHistory.value = res.list || []
  } catch (error) {
    console.error('获取消息列表失败:', error)
  } finally {
    messageLoading.value = false
  }
}

// 发送消息
const sendMessageNow = async () => {
  if (!messageFormRef.value) return
  await messageFormRef.value.validate()
  
  messageSending.value = true
  try {
    let res
    if (messageForm.value.type === 'all') {
      res = await batchSendMessage({
        app_id: parseInt(props.appId),
        title: messageForm.value.title,
        content: messageForm.value.content,
        type: 'system'
      })
    } else if (messageForm.value.type === 'user' && messageForm.value.userIds) {
      const userIds = messageForm.value.userIds.split(',').map(id => parseInt(id.trim()))
      res = await batchSendMessage({
        app_id: parseInt(props.appId),
        user_ids: userIds,
        title: messageForm.value.title,
        content: messageForm.value.content,
        type: 'system'
      })
    } else {
      res = await sendMessage({
        app_id: parseInt(props.appId),
        title: messageForm.value.title,
        content: messageForm.value.content,
        type: 'system'
      })
    }
    
    // request.js已解包，成功时不会抛出异常
    ElMessage.success('消息推送成功')
    messageForm.value = { type: 'all', userIds: '', title: '', content: '' }
    fetchMessageList()
  } catch (error) {
    console.error('发送消息失败:', error)
  } finally {
    messageSending.value = false
  }
}

// 获取版本列表
const fetchVersionList = async () => {
  if (!props.appId || props.appId === '') return
  versionLoading.value = true
  try {
    const params = {
      app_id: props.appId,
      page: versionPage.value,
      size: versionPageSize.value
    }
    if (versionPlatform.value) params.platform = versionPlatform.value
    if (versionStatus.value) params.status = versionStatus.value
    
    const res = await getVersionList(params)
    // request.js已解包，res直接是数据对象
    versionList.value = res.list || []
    versionTotal.value = res.total || versionList.value.length
  } catch (error) {
    console.error('获取版本列表失败:', error)
  } finally {
    versionLoading.value = false
  }
}

// 编辑版本
const editVersion = (row) => {
  editingVersion.value = row
  versionForm.value = {
    version: row.version,
    platform: row.platform,
    download_url: row.download_url || '',
    description: row.description,
    force_update: row.force_update,
    gray_release: row.gray_release || false,
    gray_percent: row.gray_percent || 10
  }
  showVersionDialog.value = true
}

// 提交版本
const submitVersion = async () => {
  if (!versionFormRef.value) return
  await versionFormRef.value.validate()
  
  versionSubmitting.value = true
  try {
    const data = {
      app_id: parseInt(props.appId),
      ...versionForm.value
    }
    
    await createVersion(data)
    // request.js已解包，成功时不会抛出异常
    ElMessage.success(editingVersion.value ? '版本更新成功' : '版本发布成功')
    showVersionDialog.value = false
    editingVersion.value = null
    versionForm.value = {
      version: '',
      platform: 'android',
      download_url: '',
      description: '',
      force_update: false,
      gray_release: false,
      gray_percent: 10
    }
    fetchVersionList()
  } catch (error) {
    console.error('提交版本失败:', error)
  } finally {
    versionSubmitting.value = false
  }
}

// 发布版本
const publishVersionAction = async (row) => {
  try {
    await ElMessageBox.confirm('确定要发布该版本吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await publishVersion(row.id)
    // request.js已解包，成功时不会抛出异常
    ElMessage.success('版本发布成功')
    fetchVersionList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('发布版本失败:', error)
    }
  }
}

// 下线版本
const offlineVersionAction = async (row) => {
  try {
    await ElMessageBox.confirm('确定要下线该版本吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await offlineVersion(row.id)
    // request.js已解包，成功时不会抛出异常
    ElMessage.success('版本已下线')
    fetchVersionList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('下线版本失败:', error)
    }
  }
}

// ========== 存储服务函数 ==========

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(2) + ' ' + units[i]
}

// 获取文件类型标签
const getFileTypeLabel = (mimeType) => {
  if (!mimeType) return '其他'
  if (mimeType.startsWith('image/')) return '图片'
  if (mimeType.startsWith('video/')) return '视频'
  if (mimeType.startsWith('audio/')) return '音频'
  if (mimeType.includes('pdf') || mimeType.includes('document') || mimeType.includes('word')) return '文档'
  return '其他'
}

// 获取文件列表
const fetchFileList = async () => {
  if (!props.appId || props.appId === '') return
  fileLoading.value = true
  try {
    const params = {
      app_id: props.appId,
      page: filePage.value,
      size: filePageSize.value
    }
    if (fileSearch.value) params.keyword = fileSearch.value
    if (fileType.value) params.type = fileType.value
    
    const res = await getFileList(params)
    // request.js已解包，res直接是数据对象
    fileList.value = res.list || []
    fileTotal.value = res.total || 0
  } catch (error) {
    console.error('获取文件列表失败:', error)
  } finally {
    fileLoading.value = false
  }
}

// 获取文件统计
const fetchFileStats = async () => {
  if (!props.appId || props.appId === '') return
  try {
    const res = await getFileStats({ app_id: props.appId })
    // request.js已解包，res直接是数据对象
    fileStats.value = res || { total_count: 0, total_size: 0, today_count: 0 }
  } catch (error) {
    console.error('获取文件统计失败:', error)
  }
}

// 文件选择变化
const handleFileSelectionChange = (selection) => {
  selectedFiles.value = selection
}

// 上传成功
const handleUploadSuccess = (response) => {
  if (response.code === 0) {
    ElMessage.success('文件上传成功')
    fetchFileList()
    fetchFileStats()
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

// 上传失败
const handleUploadError = () => {
  ElMessage.error('文件上传失败')
}

// 下载文件
const downloadFileAction = async (row) => {
  try {
    const res = await downloadFile(row.id)
    // request.js已解包，res直接是数据对象
    if (res && res.url) {
      window.open(res.url, '_blank')
    }
  } catch (error) {
    console.error('下载文件失败:', error)
  }
}

// 删除文件
const deleteFileAction = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该文件吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteFile(row.id)
    // request.js已解包，成功时不会抛出异常
    ElMessage.success('文件删除成功')
    fetchFileList()
    fetchFileStats()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除文件失败:', error)
    }
  }
}

// 批量删除文件
const batchDeleteFilesAction = async () => {
  if (selectedFiles.value.length === 0) return
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedFiles.value.length} 个文件吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const ids = selectedFiles.value.map(f => f.id)
    await batchDeleteFiles({ ids })
    // request.js已解包，成功时不会抛出异常
    ElMessage.success('文件删除成功')
    selectedFiles.value = []
    fetchFileList()
    fetchFileStats()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除文件失败:', error)
    }
  }
}

// ========== 数据埋点函数 ==========

// 获取事件列表
const fetchEventList = async () => {
  if (!props.appId || props.appId === '') return
  eventLoading.value = true
  try {
    const params = {
      app_id: props.appId,
      page: eventPage.value,
      size: eventPageSize.value
    }
    if (eventDateRange.value && eventDateRange.value.length === 2) {
      params.start_date = eventDateRange.value[0]
      params.end_date = eventDateRange.value[1]
    }
    if (eventName.value) params.event_name = eventName.value
    
    const res = await getEventList(params)
    // request.js已解包，res直接是数据对象
    eventList.value = res.list || []
    eventTotal.value = res.total || 0
  } catch (error) {
    console.error('获取事件列表失败:', error)
  } finally {
    eventLoading.value = false
  }
}

// 获取事件统计
const fetchEventStats = async () => {
  if (!props.appId || props.appId === '') return
  try {
    const res = await getEventStats({ app_id: props.appId })
    // request.js已解包，res直接是数据对象
    eventStats.value = res || { total_events: 0, today_events: 0, unique_users: 0 }
  } catch (error) {
    console.error('获取事件统计失败:', error)
  }
}

// 导出事件
const exportEvents = () => {
  ElMessage.info('导出功能开发中...')
}

// 漏斗分析
const analyzeFunnel = async () => {
  if (funnelSteps.value.length < 2) {
    ElMessage.warning('请至少选择2个步骤')
    return
  }
  try {
    const res = await getFunnelAnalysis({
      app_id: props.appId,
      steps: funnelSteps.value
    })
    // request.js已解包，res直接是数据对象
    initFunnelChart(res)
  } catch (error) {
    console.error('漏斗分析失败:', error)
  }
}

// 初始化漏斗图表
const initFunnelChart = (data) => {
  if (!funnelChartRef.value) return
  if (funnelChart) funnelChart.dispose()
  funnelChart = echarts.init(funnelChartRef.value)
  funnelChart.setOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    series: [{
      name: '漏斗分析',
      type: 'funnel',
      left: '10%',
      width: '80%',
      label: { formatter: '{b}: {c}' },
      data: data || [
        { value: 100, name: '页面浏览' },
        { value: 80, name: '按钮点击' },
        { value: 60, name: '表单提交' },
        { value: 40, name: '注册' },
        { value: 20, name: '支付' }
      ]
    }]
  })
}

// 获取事件定义列表
const fetchEventDefinitions = async () => {
  if (!props.appId || props.appId === '') return
  eventDefLoading.value = true
  try {
    const res = await getEventDefinitions({ app_id: props.appId })
    // request.js已解包，res直接是数据数组
    eventDefinitions.value = res || []
  } catch (error) {
    console.error('获取事件定义失败:', error)
  } finally {
    eventDefLoading.value = false
  }
}

// 删除事件定义
const deleteEventDefAction = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该事件定义吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteEventDefinition(row.id)
    // request.js已解包，成功时不会抛出异常
    ElMessage.success('删除成功')
    fetchEventDefinitions()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除事件定义失败:', error)
    }
  }
}

// ========== 监控告警函数 ==========

// 获取监控统计
const fetchMonitorStats = async () => {
  if (!props.appId || props.appId === '') return
  try {
    const res = await getMonitorStats({ app_id: props.appId })
    // request.js已解包，res直接是数据对象
    monitorStats.value = res || { cpu_usage: 0, memory_usage: 0, active_alerts: 0 }
  } catch (error) {
    console.error('获取监控统计失败:', error)
  }
}

// 获取健康检查
const fetchHealthCheck = async () => {
  if (!props.appId || props.appId === '') return
  try {
    const res = await getHealthCheck({ app_id: props.appId })
    // request.js已解包，res直接是数据对象
    healthStatus.value = res.status || 'healthy'
  } catch (error) {
    healthStatus.value = 'unhealthy'
    console.error('健康检查失败:', error)
  }
}

// 获取监控指标
const fetchMonitorMetrics = async () => {
  if (!props.appId || props.appId === '') return
  try {
    const res = await getMonitorMetrics({
      app_id: props.appId,
      period: monitorPeriod.value
    })
    // request.js已解包，res直接是数据对象
    initMonitorCharts(res)
  } catch (error) {
    console.error('获取监控指标失败:', error)
  }
}

// 初始化监控图表
const initMonitorCharts = (data) => {
  // 请求量图表
  if (requestMetricChartRef.value) {
    if (requestMetricChart) requestMetricChart.dispose()
    requestMetricChart = echarts.init(requestMetricChartRef.value)
    requestMetricChart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: data?.times || ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00'] },
      yAxis: { type: 'value' },
      series: [{
        name: '请求量',
        type: 'line',
        smooth: true,
        data: data?.requests || [120, 200, 150, 80, 70, 110],
        areaStyle: { opacity: 0.3 }
      }]
    })
  }
  
  // 错误率图表
  if (errorMetricChartRef.value) {
    if (errorMetricChart) errorMetricChart.dispose()
    errorMetricChart = echarts.init(errorMetricChartRef.value)
    errorMetricChart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: data?.times || ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00'] },
      yAxis: { type: 'value', axisLabel: { formatter: '{value}%' } },
      series: [{
        name: '错误率',
        type: 'line',
        smooth: true,
        data: data?.errors || [0.5, 0.8, 0.3, 1.2, 0.6, 0.4],
        itemStyle: { color: '#f56c6c' },
        areaStyle: { opacity: 0.3, color: '#f56c6c' }
      }]
    })
  }
}

// 获取告警列表
const fetchAlertList = async () => {
  if (!props.appId || props.appId === '') return
  alertLoading.value = true
  try {
    const params = { app_id: props.appId }
    if (alertStatus.value !== '') params.status = parseInt(alertStatus.value)
    
    const res = await getAlertList(params)
    // request.js已解包，res直接是数据数组
    alertList.value = res || []
  } catch (error) {
    console.error('获取告警列表失败:', error)
  } finally {
    alertLoading.value = false
  }
}

// 处理告警
const resolveAlert = async (row) => {
  try {
    await updateAlert(row.id, { status: 1 })
    // request.js已解包，成功时不会抛出异常
    ElMessage.success('告警已处理')
    fetchAlertList()
    fetchMonitorStats()
  } catch (error) {
    console.error('处理告警失败:', error)
  }
}

// 切换告警规则状态
const toggleAlertRule = async (row) => {
  try {
    await updateAlert(row.id, { status: row.status })
    // request.js已解包，成功时不会抛出异常
    ElMessage.success(row.status === 1 ? '规则已启用' : '规则已禁用')
  } catch (error) {
    console.error('更新告警规则失败:', error)
  }
}

// 删除告警规则
const deleteAlertRule = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该告警规则吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteAlert(row.id)
    // request.js已解包，成功时不会抛出异常
    ElMessage.success('删除成功')
    fetchAlertRules()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除告警规则失败:', error)
    }
  }
}

// 获取告警规则列表
const fetchAlertRules = async () => {
  if (!props.appId || props.appId === '') return
  alertRuleLoading.value = true
  try {
    const res = await getAlertList({ app_id: props.appId, type: 'rule' })
    // request.js已解包，res直接是数据数组
    alertRules.value = res || []
  } catch (error) {
    console.error('获取告警规则失败:', error)
  } finally {
    alertRuleLoading.value = false
  }
}

// 提交告警规则
const submitAlertRule = async () => {
  if (!alertRuleFormRef.value) return
  try {
    await alertRuleFormRef.value.validate()
    alertRuleSubmitting.value = true
    await createAlert({
      app_id: props.appId,
      type: 'rule',
      ...alertRuleForm.value
    })
    ElMessage.success('告警规则创建成功')
    showAlertRuleDialog.value = false
    // 重置表单
    alertRuleForm.value = {
      name: '',
      metric_name: '',
      condition_type: 'gt',
      threshold: 80,
      duration: 60,
      level: 'warning',
      notify_channels: []
    }
    fetchAlertRules()
  } catch (error) {
    if (error !== 'cancel' && error !== false) {
      console.error('创建告警规则失败:', error)
      ElMessage.error('创建告警规则失败')
    }
  } finally {
    alertRuleSubmitting.value = false
  }
}

// 提交事件定义
const submitEventDef = async () => {
  if (!eventDefFormRef.value) return
  try {
    await eventDefFormRef.value.validate()
    eventDefSubmitting.value = true
    await createEventDefinition({
      app_id: props.appId,
      ...eventDefForm.value,
      properties: eventDefForm.value.properties.filter(p => p.trim() !== '')
    })
    ElMessage.success('事件定义创建成功')
    showEventDefDialog.value = false
    // 重置表单
    eventDefForm.value = {
      name: '',
      code: '',
      description: '',
      properties: []
    }
    fetchEventDefinitions()
  } catch (error) {
    if (error !== 'cancel' && error !== false) {
      console.error('创建事件定义失败:', error)
      ElMessage.error('创建事件定义失败')
    }
  } finally {
    eventDefSubmitting.value = false
  }
}

// 审计日志相关函数
const fetchAuditLogs = async () => {
  if (!props.appId || props.appId === '') return
  auditLoading.value = true
  try {
    const params = {
      app_id: props.appId,
      page: auditPage.value,
      page_size: auditPageSize.value,
      keyword: auditSearch.value.keyword,
      action: auditSearch.value.action,
      resource: auditSearch.value.resource
    }
    if (auditSearch.value.dateRange && auditSearch.value.dateRange.length === 2) {
      params.start_time = formatDateForApi(auditSearch.value.dateRange[0])
      params.end_time = formatDateForApi(auditSearch.value.dateRange[1])
    }
    const res = await getAuditLogs(params)
    // request.js已解包，res直接是数据对象
    auditLogs.value = res.list || []
    auditTotal.value = res.total || 0
  } catch (error) {
    console.error('获取审计日志失败:', error)
  } finally {
    auditLoading.value = false
  }
}

const fetchAuditStats = async () => {
  if (!props.appId || props.appId === '') return
  try {
    const res = await getAuditStats({ app_id: props.appId, days: 7 })
    // request.js已解包，res直接是数据对象
    auditStats.value = res || {}
  } catch (error) {
    console.error('获取审计统计失败:', error)
  }
}

const exportAuditLogsData = async () => {
  try {
    const params = {
      app_id: props.appId,
      format: 'csv'
    }
    if (auditSearch.value.dateRange && auditSearch.value.dateRange.length === 2) {
      params.start_time = formatDateForApi(auditSearch.value.dateRange[0])
      params.end_time = formatDateForApi(auditSearch.value.dateRange[1])
    }
    const res = await exportAuditLogs(params)
    const blob = new Blob([res], { type: 'text/csv' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `audit_logs_${new Date().toISOString().slice(0, 10)}.csv`
    a.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出审计日志失败:', error)
    ElMessage.error('导出失败')
  }
}

const formatDateForApi = (date) => {
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} 00:00:00`
}

const getActionTagType = (action) => {
  const types = {
    view: 'info',
    create: 'success',
    update: 'warning',
    delete: 'danger',
    login: 'success',
    logout: 'info',
    export: 'primary'
  }
  return types[action] || 'info'
}

const getActionLabel = (action) => {
  const labels = {
    view: '查看',
    create: '创建',
    update: '更新',
    delete: '删除',
    login: '登录',
    logout: '登出',
    export: '导出',
    send: '发送',
    publish: '发布'
  }
  return labels[action] || action
}

const getResourceLabel = (resource) => {
  const labels = {
    user: '用户',
    app: '应用',
    config: '配置',
    message: '消息',
    push: '推送',
    file: '文件',
    version: '版本',
    log: '日志',
    event: '事件',
    monitor: '监控'
  }
  return labels[resource] || resource
}

// 初始化图表
const initCharts = () => {
  if (requestChartRef.value) {
    requestChart = echarts.init(requestChartRef.value)
    requestChart.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: ['1/4', '1/5', '1/6', '1/7', '1/8', '1/9', '1/10']
      },
      yAxis: { type: 'value' },
      series: [{
        name: '请求数',
        type: 'line',
        smooth: true,
        areaStyle: { opacity: 0.3 },
        data: [32000, 35000, 38000, 42000, 45000, 43000, 45678],
        itemStyle: { color: '#409eff' }
      }]
    })
  }

  if (moduleChartRef.value) {
    moduleChart = echarts.init(moduleChartRef.value)
    moduleChart.setOption({
      tooltip: { trigger: 'item' },
      legend: { orient: 'vertical', left: 'left' },
      series: [{
        name: '模块调用',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
        label: { show: false, position: 'center' },
        emphasis: {
          label: { show: true, fontSize: 16, fontWeight: 'bold' }
        },
        data: [
          { value: 1048, name: '用户管理' },
          { value: 735, name: '消息推送' },
          { value: 580, name: '数据埋点' },
          { value: 484, name: '日志服务' },
          { value: 300, name: '版本管理' }
        ]
      }]
    })
  }
}

// 加载数据
const loadData = () => {
  fetchUserStats()
  fetchLogStats()
}

// WebSocket相关变量
const wsConnected = ref(false)
const realtimeAlerts = ref([])

// 初始化WebSocket连接
const initWebSocket = () => {
  if (!props.appId || props.appId === '') return
  
  // 请求通知权限
  wsClient.constructor.requestNotificationPermission()
  
  // 连接WebSocket
  wsClient.connect(props.appId, 'admin')
  
  // 监听连接事件
  wsClient.on('connected', () => {
    wsConnected.value = true
    console.log('[Workspace] WebSocket connected')
  })
  
  wsClient.on('disconnected', () => {
    wsConnected.value = false
    console.log('[Workspace] WebSocket disconnected')
  })
  
  // 监听监控数据
  wsClient.on('monitor', (data) => {
    console.log('[Workspace] Monitor data:', data)
    if (data.cpu_usage !== undefined) {
      monitorStats.value.cpu_usage = data.cpu_usage
    }
    if (data.memory_usage !== undefined) {
      monitorStats.value.memory_usage = data.memory_usage
    }
    if (data.request_count !== undefined) {
      // 更新请求数图表
      updateRequestChart(data)
    }
  })
  
  // 监听告警事件
  wsClient.on('alert', (data) => {
    console.log('[Workspace] Alert:', data)
    realtimeAlerts.value.unshift(data)
    if (realtimeAlerts.value.length > 10) {
      realtimeAlerts.value.pop()
    }
    // 刷新告警列表
    if (currentMenu.value === 'monitor' && monitorTab.value === 'alerts') {
      fetchAlertList()
    }
    // 显示告警提示
    const levelMap = { critical: 'error', warning: 'warning', info: 'info' }
    ElMessage[levelMap[data.level] || 'warning']({
      message: `告警: ${data.title || data.message}`,
      duration: 5000
    })
  })
  
  // 监听通知事件
  wsClient.on('notification', (data) => {
    console.log('[Workspace] Notification:', data)
    ElMessage.info({
      message: data.title || data.message,
      duration: 3000
    })
  })
}

// 更新请求数图表
const updateRequestChart = (data) => {
  if (requestMetricChart && data.timestamp && data.request_count !== undefined) {
    const chart = requestMetricChart
    const option = chart.getOption()
    if (option && option.xAxis && option.series) {
      const time = new Date(data.timestamp).toLocaleTimeString()
      option.xAxis[0].data.push(time)
      option.series[0].data.push(data.request_count)
      // 保持最近20个数据点
      if (option.xAxis[0].data.length > 20) {
        option.xAxis[0].data.shift()
        option.series[0].data.shift()
      }
      chart.setOption(option)
    }
  }
}

// 断开WebSocket连接
const disconnectWebSocket = () => {
  wsClient.disconnect()
  wsConnected.value = false
}

onMounted(() => {
  setTimeout(initCharts, 100)
  // 只有当appId有效时才加载数据
  if (props.appId && props.appId !== '') {
    // 根据当前菜单加载对应数据
    const menu = props.initialMenu || currentMenu.value
    if (menu === 'overview') {
      loadData()
    } else if (menu === 'users') {
      fetchUserList()
    } else if (menu === 'logs') {
      fetchLogList()
      fetchLogStats()
    } else if (menu === 'messages') {
      fetchMessageList()
    } else if (menu === 'versions') {
      fetchVersionList()
    } else if (menu === 'storage') {
      fetchFileList()
      fetchFileStats()
    } else if (menu === 'events') {
      fetchEventList()
      fetchEventStats()
    } else if (menu === 'monitor') {
      fetchMonitorStats()
      fetchHealthCheck()
      setTimeout(() => fetchMonitorMetrics(), 100)
    } else if (menu === 'audit') {
      fetchAuditLogs()
      fetchAuditStats()
    }
    initWebSocket()
  }
})

onUnmounted(() => {
  disconnectWebSocket()
})

watch(currentMenu, (val) => {
  // 确保appId有效时才加载数据
  if (!props.appId || props.appId === '') {
    console.log('currentMenu changed but appId is empty, skipping data load')
    return
  }
  
  if (val === 'overview') {
    setTimeout(initCharts, 100)
    loadData()
  } else if (val === 'users') {
    fetchUserList()
  } else if (val === 'logs') {
    fetchLogList()
    fetchLogStats()
  } else if (val === 'messages') {
    fetchMessageList()
  } else if (val === 'versions') {
    fetchVersionList()
  } else if (val === 'storage') {
    fetchFileList()
    fetchFileStats()
  } else if (val === 'events') {
    fetchEventList()
    fetchEventStats()
    if (eventTab.value === 'definitions') fetchEventDefinitions()
  } else if (val === 'monitor') {
    fetchMonitorStats()
    fetchHealthCheck()
    setTimeout(() => fetchMonitorMetrics(), 100)
    if (monitorTab.value === 'alerts') fetchAlertList()
    if (monitorTab.value === 'rules') fetchAlertRules()
  } else if (val === 'audit') {
    fetchAuditLogs()
    fetchAuditStats()
  }
})

watch(() => props.appId, (newVal, oldVal) => {
  // 确保appId有效时才加载数据
  if (!newVal || newVal === '') return
  // 如果appId变化了，重新加载数据
  if (oldVal && oldVal !== newVal) {
    loadData()
    if (currentMenu.value === 'users') fetchUserList()
    if (currentMenu.value === 'logs') fetchLogList()
    if (currentMenu.value === 'messages') fetchMessageList()
    if (currentMenu.value === 'versions') fetchVersionList()
    if (currentMenu.value === 'storage') { fetchFileList(); fetchFileStats() }
    if (currentMenu.value === 'events') { fetchEventList(); fetchEventStats() }
    if (currentMenu.value === 'monitor') { fetchMonitorStats(); fetchHealthCheck(); fetchMonitorMetrics() }
  }
})

// 监听父组件传入的initialMenu变化（移动端菜单切换）
watch(() => props.initialMenu, (newVal) => {
  if (newVal && newVal !== currentMenu.value) {
    currentMenu.value = newVal
    // watch currentMenu会自动处理数据加载
  }
})
</script>

<style lang="scss" scoped>
.workspace {
  display: flex;
  height: 100%;
  background: #f5f7fa;
}

.workspace-sidebar {
  width: 200px;
  background: white;
  border-right: 1px solid #e4e7ed;
  padding: 16px 0;
  display: flex;
  flex-direction: column;
  height: 100%;
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

.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 24px;
  cursor: pointer;
  color: #606266;
  font-size: 14px;
  transition: all 0.2s;
  
  &:hover {
    background: #f5f7fa;
    color: #409eff;
  }
  
  &.active {
    background: linear-gradient(90deg, #ecf5ff 0%, transparent 100%);
    color: #409eff;
    border-left: 3px solid #409eff;
    font-weight: 500;
  }
}

.workspace-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.content-section {
  max-width: 1400px;
}

.section-header {
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  
  .header-left {
    flex: 1;
  }
  
  .header-right {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  
  h2 {
    font-size: 22px;
    font-weight: 600;
    color: #1a1a2e;
    margin: 0 0 8px;
  }
  
  p {
    font-size: 14px;
    color: #909399;
    margin: 0;
  }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  transition: all 0.3s;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  }
  
  &.blue {
    border-left: 4px solid #409eff;
    .stat-icon { background: #ecf5ff; color: #409eff; }
  }
  
  &.green {
    border-left: 4px solid #67c23a;
    .stat-icon { background: #f0f9eb; color: #67c23a; }
  }
  
  &.orange {
    border-left: 4px solid #e6a23c;
    .stat-icon { background: #fdf6ec; color: #e6a23c; }
  }
  
  &.red {
    border-left: 4px solid #f56c6c;
    .stat-icon { background: #fef0f0; color: #f56c6c; }
  }
  
  .stat-content {
    .stat-value {
      font-size: 32px;
      font-weight: 700;
      color: #1a1a2e;
      line-height: 1;
    }
    
    .stat-label {
      font-size: 14px;
      color: #909399;
      margin-top: 8px;
    }
    
    .stat-trend {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 13px;
      margin-top: 8px;
      
      &.up { color: #67c23a; }
      &.down { color: #f56c6c; }
    }
  }
  
  .stat-icon {
    width: 56px;
    height: 56px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 28px;
  }
}

.charts-row {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
}

.chart-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  
  .chart-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    
    h3 {
      font-size: 16px;
      font-weight: 600;
      color: #1a1a2e;
      margin: 0;
    }
  }
  
  .chart-body {
    height: 300px;
  }
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px 20px;
  background: white;
  border-radius: 8px;
  
  .search-area {
    display: flex;
    gap: 12px;
  }
  
  .action-area {
    display: flex;
    gap: 12px;
  }
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.message-form {
  background: white;
  padding: 24px;
  border-radius: 8px;
}

.log-stats {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
  
  .log-stat-item {
    background: white;
    padding: 16px 24px;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-width: 120px;
    
    .label {
      font-size: 13px;
      color: #909399;
    }
    
    .value {
      font-size: 24px;
      font-weight: 600;
      color: #1a1a2e;
    }
    
    &.error {
      border-left: 3px solid #f56c6c;
      .value { color: #f56c6c; }
    }
    
    &.warn {
      border-left: 3px solid #e6a23c;
      .value { color: #e6a23c; }
    }
    
    &.info {
      border-left: 3px solid #409eff;
      .value { color: #409eff; }
    }
  }
}

.log-list {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  min-height: 200px;
}

.empty-logs {
  padding: 40px;
}

.log-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
  font-size: 13px;
  
  &:last-child {
    border-bottom: none;
  }
  
  &.error {
    background: #fff1f0;
  }
  
  &.warn {
    background: #fffbe6;
  }
  
  .log-time {
    color: #909399;
    font-family: monospace;
    white-space: nowrap;
    min-width: 160px;
  }
  
  .log-module {
    color: #606266;
    background: #f0f0f0;
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 12px;
  }
  
  .log-content {
    flex: 1;
    color: #303133;
  }
}

.form-tip {
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}

// 存储服务和新模块样式
.stats-grid.small {
  grid-template-columns: repeat(3, 1fr);
}

.properties-preview {
  font-family: monospace;
  font-size: 12px;
  color: #606266;
}

.funnel-config {
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
  margin-bottom: 20px;
}

.funnel-chart {
  height: 400px;
  background: white;
  border-radius: 8px;
  padding: 16px;
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .charts-row {
    grid-template-columns: 1fr;
  }
  
  .log-stats {
    flex-wrap: wrap;
  }
}

@media (max-width: 768px) {
  .workspace {
    flex-direction: column;
  }
  
  .workspace-sidebar {
    display: none; /* 移动端隐藏侧边栏，使用汉堡菜单代替 */
  }
  
  .workspace-content {
    padding: 16px;
  }
  
  .menu-item {
    padding: 12px 16px;
    white-space: nowrap;
    border-left: none !important;
    
    &.active {
      border-bottom: 2px solid #409eff;
      background: transparent;
    }
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .toolbar {
    flex-direction: column;
    gap: 12px;
    
    .search-area, .action-area {
      width: 100%;
      flex-wrap: wrap;
    }
  }
}
</style>
