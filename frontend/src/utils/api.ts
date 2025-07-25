// API 接口定义，用于与 Wails 后端通信
declare global {
  interface Window {
    go: {
      app: {
        App: {
          Login: (email: string, password: string) => Promise<any>
          SyncSubscription: () => Promise<any>
          GetServers: () => Promise<any>
          ToggleServer: (id: number) => Promise<any>
          StartConnection: () => Promise<any>
          StopConnection: () => Promise<any>
          GetConnectionStatus: () => Promise<any>
          GetUserStatus: () => Promise<any>
          PingServer: (id: number) => Promise<any>
          GetConnectionLogs: () => Promise<any>
          GetConfig: () => Promise<any>
          UpdateConfig: (config: any) => Promise<any>
        }
      }
    }
  }
}

export const api = {
  // 用户认证
  login: async (email: string, password: string) => {
    return await window.go.app.App.Login(email, password)
  },

  // 订阅管理
  syncSubscription: async () => {
    return await window.go.app.App.SyncSubscription()
  },

  // 服务器管理
  getServers: async () => {
    return await window.go.app.App.GetServers()
  },

  toggleServer: async (id: number) => {
    return await window.go.app.App.ToggleServer(id)
  },

  pingServer: async (id: number) => {
    return await window.go.app.App.PingServer(id)
  },

  // 连接管理
  startConnection: async () => {
    return await window.go.app.App.StartConnection()
  },

  stopConnection: async () => {
    return await window.go.app.App.StopConnection()
  },

  getConnectionStatus: async () => {
    return await window.go.app.App.GetConnectionStatus()
  },

  // 用户信息
  getUserStatus: async () => {
    return await window.go.app.App.GetUserStatus()
  },

  // 日志
  getConnectionLogs: async () => {
    return await window.go.app.App.GetConnectionLogs()
  },

  // 配置
  getConfig: async () => {
    return await window.go.app.App.GetConfig()
  },

  updateConfig: async (config: any) => {
    return await window.go.app.App.UpdateConfig(config)
  }
}

// 格式化流量大小
export const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化时间
export const formatTime = (timestamp: number): string => {
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}