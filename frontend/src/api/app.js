import request from '@/utils/request'

// APP管理
export const getAppList = () => request.get('/apps')
export const createApp = (data) => request.post('/apps', data)
export const updateApp = (id, data) => request.put(`/apps/${id}`, data)
export const deleteApp = (id) => request.delete(`/apps/${id}`)
export const resetAppSecret = (id) => request.post(`/apps/${id}/reset-secret`)
export const getAppDetail = (id) => request.get(`/apps/${id}`)

// 用户管理
export const getUserList = (params) => request.get('/users', { params })
export const getUserDetail = (id) => request.get(`/users/${id}`)
export const updateUserStatus = (id, status) => request.put(`/users/${id}/status`, { status })
export const getUserStats = (appId) => request.get('/users/stats', { params: { app_id: appId } })

// 日志服务
export const getLogList = (params) => request.get('/logs', { params })
export const reportLog = (data) => request.post('/logs', data)
export const getLogStats = (appId) => request.get('/logs/stats', { params: { app_id: appId } })
export const exportLogs = (params) => request.get('/logs/export', { params })
export const cleanLogs = (data) => request.delete('/logs/clean', { data })

// 消息中心
export const getMessageList = (params) => request.get('/messages', { params })
export const sendMessage = (data) => request.post('/messages', data)
export const getMessageDetail = (id) => request.get(`/messages/${id}`)
export const markMessageRead = (id) => request.put(`/messages/${id}/read`)
export const getMessageStats = (appId) => request.get('/messages/stats', { params: { app_id: appId } })
export const batchSendMessage = (data) => request.post('/messages/batch', data)

// 版本管理
export const getVersionList = (params) => request.get('/versions', { params })
export const createVersion = (data) => request.post('/versions', data)
export const publishVersion = (id) => request.post(`/versions/${id}/publish`)
export const offlineVersion = (id) => request.post(`/versions/${id}/offline`)
export const checkUpdate = (params) => request.get('/versions/check', { params })

// 推送服务
export const getPushList = (params) => request.get('/push', { params })
export const createPush = (data) => request.post('/push', data)
export const getPushDetail = (id) => request.get(`/push/${id}`)
export const sendPush = (id) => request.post(`/push/${id}/send`)
export const cancelPush = (id) => request.post(`/push/${id}/cancel`)
export const getPushStats = (appId) => request.get('/push/stats', { params: { app_id: appId } })

// 数据埋点
export const reportEvent = (data) => request.post('/events', data)
export const batchReportEvents = (data) => request.post('/events/batch', data)
export const getEventList = (params) => request.get('/events', { params })
export const getEventStats = (appId) => request.get('/events/stats', { params: { app_id: appId } })
export const getFunnelAnalysis = (params) => request.get('/events/funnel', { params })

// 监控告警
export const getMonitorMetrics = (params) => request.get('/monitor/metrics', { params })
export const reportMetrics = (data) => request.post('/monitor/metrics', data)
export const getAlertList = (params) => request.get('/monitor/alerts', { params })
export const createAlert = (data) => request.post('/monitor/alerts', data)
export const getMonitorStats = (appId) => request.get('/monitor/stats', { params: { app_id: appId } })

// 配置管理
export const getModuleConfig = (appId, moduleCode) => request.get('/configs', { params: { app_id: appId, module_code: moduleCode } })
export const saveModuleConfig = (data) => request.post('/configs', data)
export const resetModuleConfig = (appId, moduleCode) => request.post('/configs/reset', { app_id: appId, module_code: moduleCode })
export const getConfigHistory = (appId, moduleCode) => request.get('/configs/history', { params: { app_id: appId, module_code: moduleCode } })

// 统计数据
export const getDashboardStats = () => request.get('/stats')

// 存储服务
export const uploadFile = (formData) => request.post('/files', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
export const getFileList = (params) => request.get('/files', { params })
export const downloadFile = (id) => request.get(`/files/download/${id}`, { responseType: 'blob' })
export const deleteFile = (id) => request.delete(`/files/${id}`)
export const getFileStats = (appId) => request.get('/files/stats', { params: { app_id: appId } })
export const batchDeleteFiles = (ids) => request.post('/files/batch-delete', { ids })

// 事件定义管理
export const getEventDefinitions = (params) => request.get('/events/definitions', { params })
export const createEventDefinition = (data) => request.post('/events/definitions', data)
export const updateEventDefinition = (id, data) => request.put(`/events/definitions/${id}`, data)
export const deleteEventDefinition = (id) => request.delete(`/events/definitions/${id}`)

// 告警规则管理
export const updateAlert = (id, data) => request.put(`/monitor/alerts/${id}`, data)
export const deleteAlert = (id) => request.delete(`/monitor/alerts/${id}`)
export const getHealthCheck = () => request.get('/monitor/health')

// 审计日志API
export const getAuditLogs = (params) => request.get('/audit', { params })
export const getAuditStats = (params) => request.get('/audit/stats', { params })
export const exportAuditLogs = (params) => request.get('/audit/export', { params, responseType: 'blob' })
