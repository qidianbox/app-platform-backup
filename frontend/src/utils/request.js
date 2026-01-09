import axios from 'axios'
import { ElMessage } from 'element-plus'

// 根据环境自动选择API地址
const getBaseURL = () => {
  // 如果是通过代理访问（开发模式），使用相对路径
  if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
    return '/api/v1'
  }
  // 生产环境或通过公网访问时，使用绝对路径
  // 将前端域名中的5173或5174替换为8080
  const apiHost = window.location.origin.replace(/517[34]/, '8080')
  return `${apiHost}/api/v1`
}

const request = axios.create({
  baseURL: getBaseURL(),
  timeout: 30000
})

request.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  response => response.data,
  error => {
    ElMessage.error(error.response?.data?.message || '请求失败')
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default request
