import request from '@/utils/request'

export const getAppList = () => request.get('/apps')
export const createApp = (data) => request.post('/apps', data)
export const updateApp = (id, data) => request.put(`/apps/${id}`, data)
export const deleteApp = (id) => request.delete(`/apps/${id}`)
export const resetAppSecret = (id) => request.post(`/apps/${id}/reset-secret`)
