export const moduleCategories = [
  { id: 'user', name: '用户与权限', icon: '👤', description: '用户管理和权限控制' },
  { id: 'message', name: '消息与推送', icon: '📬', description: '消息中心和推送服务' },
  { id: 'data', name: '数据与分析', icon: '📊', description: '埋点和数据分析' },
  { id: 'system', name: '系统服务', icon: '⚙️', description: '日志、监控等系统服务' },
  { id: 'storage', name: '存储服务', icon: '📁', description: '文件存储和配置管理' },
  { id: 'other', name: '其他', icon: '📦', description: '其他功能模块' }
]

export const getGroupedModules = (modules) => {
  return moduleCategories.map(cat => ({
    ...cat,
    modules: modules.filter(m => {
      if (cat.id === 'user') return m.module_code?.includes('user')
      if (cat.id === 'message') return m.module_code?.includes('message') || m.module_code?.includes('push')
      if (cat.id === 'data') return m.module_code?.includes('event') || m.module_code?.includes('stats')
      if (cat.id === 'system') return m.module_code?.includes('log') || m.module_code?.includes('monitor')
      if (cat.id === 'storage') return m.module_code?.includes('file') || m.module_code?.includes('config')
      return true
    })
  }))
}
