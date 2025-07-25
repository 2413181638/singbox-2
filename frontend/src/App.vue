<template>
  <div id="app">
    <el-container class="app-container">
      <el-aside width="200px" class="sidebar">
        <div class="logo">
          <h2>SingBox Client</h2>
        </div>
        <el-menu
          :default-active="$route.path"
          class="sidebar-menu"
          router
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409EFF"
        >
          <el-menu-item index="/dashboard">
            <el-icon><Odometer /></el-icon>
            <span>仪表板</span>
          </el-menu-item>
          <el-menu-item index="/servers">
            <el-icon><Server /></el-icon>
            <span>服务器</span>
          </el-menu-item>
          <el-menu-item index="/logs">
            <el-icon><Document /></el-icon>
            <span>连接日志</span>
          </el-menu-item>
          <el-menu-item index="/settings">
            <el-icon><Setting /></el-icon>
            <span>设置</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      
      <el-container>
        <el-header class="header">
          <div class="header-left">
            <h3>{{ getPageTitle() }}</h3>
          </div>
          <div class="header-right">
            <el-button 
              :type="connectionStore.isRunning ? 'danger' : 'primary'"
              @click="toggleConnection"
              :loading="connectionStore.loading"
            >
              {{ connectionStore.isRunning ? '断开连接' : '连接' }}
            </el-button>
            <el-dropdown>
              <span class="user-info">
                <el-icon><User /></el-icon>
                {{ userStore.userInfo?.email || '未登录' }}
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>
        
        <el-main class="main-content">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { useConnectionStore } from '@/stores/connection'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const connectionStore = useConnectionStore()

const getPageTitle = () => {
  const titles: Record<string, string> = {
    '/dashboard': '仪表板',
    '/servers': '服务器管理',
    '/logs': '连接日志',
    '/settings': '设置'
  }
  return titles[route.path] || 'SingBox Client'
}

const toggleConnection = async () => {
  try {
    if (connectionStore.isRunning) {
      await connectionStore.stopConnection()
      ElMessage.success('连接已断开')
    } else {
      await connectionStore.startConnection()
      ElMessage.success('连接已建立')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

const logout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.app-container {
  height: 100vh;
}

.sidebar {
  background-color: #304156;
  color: white;
}

.logo {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid #434a5a;
}

.logo h2 {
  margin: 0;
  color: white;
  font-size: 18px;
}

.sidebar-menu {
  border: none;
}

.header {
  background: white;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}

.header-left h3 {
  margin: 0;
  color: #303133;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: #606266;
}

.main-content {
  background: #f5f5f5;
  padding: 20px;
}
</style>