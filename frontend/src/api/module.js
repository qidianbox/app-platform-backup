import request from '@/utils/request'

export const getModuleTemplates = () => request.get('/modules/templates')
export const getAppModules = (appId) => request.get(`/apps/${appId}/modules`)
export const updateAppModule = (appId, moduleCode, data) => request.put(`/apps/${appId}/modules/${moduleCode}`, data)
export const getModuleConfig = (appId, moduleCode) => request.get(`/apps/${appId}/modules/${moduleCode}/config`)
export const saveModuleConfig = (appId, moduleCode, data) => request.put(`/apps/${appId}/modules/${moduleCode}/config`, data)
export const resetModuleConfig = (appId, moduleCode) => request.post(`/apps/${appId}/modules/${moduleCode}/config/reset`)
