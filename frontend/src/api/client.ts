import axios from 'axios'
import router from '@/router'

const client = axios.create({
  baseURL: '/api',
  timeout: 30000,
})

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  // Correlation id traced through cloud -> bridge logs.
  if (typeof crypto !== 'undefined' && crypto.randomUUID) {
    config.headers['X-Request-Id'] = crypto.randomUUID()
  }
  return config
})

client.interceptors.response.use(
  (response) => response,
  (error) => {
    // Surface failed API calls in the devtools console (previously silent).
    console.error(
      '[api]',
      error.config?.method?.toUpperCase(),
      error.config?.url,
      error.response?.status,
      error.response?.data
    )
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

export default client
