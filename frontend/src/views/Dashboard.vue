<template>
  <div class="dashboard">
    <!-- 用户信息卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-icon user-icon">
              <el-icon><User /></el-icon>
            </div>
            <div class="stats-info">
              <h3>{{ userStore.userInfo?.email || '未登录' }}</h3>
              <p>用户账户</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-icon upload-icon">
              <el-icon><Top /></el-icon>
            </div>
            <div class="stats-info">
              <h3>{{ formatBytes(userStore.userInfo?.upload || 0) }}</h3>
              <p>上传流量</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-icon download-icon">
              <el-icon><Bottom /></el-icon>
            </div>
            <div class="stats-info">
              <h3>{{ formatBytes(userStore.userInfo?.download || 0) }}</h3>
              <p>下载流量</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-icon remaining-icon">
              <el-icon><Odometer /></el-icon>
            </div>
            <div class="stats-info">
              <h3>{{ formatBytes(userStore.userInfo?.remaining || 0) }}</h3>
              <p>剩余流量</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 连接状态和服务器信息 -->
    <el-row :gutter="20" class="content-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>连接状态</span>
              <el-button 
                :type="connectionStore.isRunning ? 'danger' : 'primary'"
                @click="toggleConnection"
                :loading="connectionStore.loading"
                size="small"
              >
                {{ connectionStore.isRunning ? '断开' : '连接' }}
              </el-button>
            </div>
          </template>
          
          <div class="connection-status">
            <div class="status-indicator">
              <div 
                :class="['status-dot', connectionStore.isRunning ? 'connected' : 'disconnected']"
              ></div>
              <span class="status-text">
                {{ connectionStore.isRunning ? '已连接' : '未连接' }}
              </span>
            </div>
            
            <div v-if="connectionStore.isRunning" class="traffic-stats">
              <div class="traffic-item">
                <span>实时上传:</span>
                <span>{{ formatBytes(connectionStore.stats?.up || 0) }}/s</span>
              </div>
              <div class="traffic-item">
                <span>实时下载:</span>
                <span>{{ formatBytes(connectionStore.stats?.down || 0) }}/s</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>服务器状态</span>
              <el-button @click="refreshServers" size="small" :loading="serversLoading">
                刷新
              </el-button>
            </div>
          </template>
          
          <div class="servers-summary">
            <div class="summary-item">
              <span>总服务器:</span>
              <span>{{ servers.length }}</span>
            </div>
            <div class="summary-item">
              <span>已启用:</span>
              <span>{{ activeServers.length }}</span>
            </div>
            <div class="summary-item">
              <span>平均延迟:</span>
              <span>{{ averagePing }}ms</span>
            </div>
          </div>
          
          <div class="server-list">
            <div 
              v-for="server in servers.slice(0, 3)" 
              :key="server.id"
              class="server-item"
            >
              <div class="server-info">
                <span class="server-name">{{ server.name }}</span>
                <span class="server-type">{{ server.server_type }}</span>
              </div>
              <div class="server-status">
                <span :class="['ping', getPingClass(server.ping)]">
                  {{ server.ping }}ms
                </span>
                <el-switch 
                  v-model="server.is_active" 
                  @change="toggleServer(server.id)"
                  size="small"
                />
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 流量使用图表 -->
    <el-row class="chart-row">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>流量使用情况</span>
          </template>
          <div class="chart-container">
            <v-chart :option="trafficChartOption" style="height: 300px;" />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { PieChart, LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components'
import VChart from 'vue-echarts'

import { useUserStore } from '@/stores/user'
import { useConnectionStore } from '@/stores/connection'
import { api, formatBytes } from '@/utils/api'

use([
  CanvasRenderer,
  PieChart,
  LineChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
])

const userStore = useUserStore()
const connectionStore = useConnectionStore()

const servers = ref<any[]>([])
const serversLoading = ref(false)
let updateTimer: number | null = null

const activeServers = computed(() => servers.value.filter(s => s.is_active))
const averagePing = computed(() => {
  const pings = servers.value.filter(s => s.ping > 0).map(s => s.ping)
  return pings.length > 0 ? Math.round(pings.reduce((a, b) => a + b, 0) / pings.length) : 0
})

const trafficChartOption = computed(() => ({
  title: {
    text: '流量使用分布',
    left: 'center'
  },
  tooltip: {
    trigger: 'item',
    formatter: '{a} <br/>{b}: {c} ({d}%)'
  },
  legend: {
    orient: 'vertical',
    left: 'left'
  },
  series: [
    {
      name: '流量使用',
      type: 'pie',
      radius: '50%',
      data: [
        { 
          value: userStore.userInfo?.upload || 0, 
          name: '已上传',
          itemStyle: { color: '#409EFF' }
        },
        { 
          value: userStore.userInfo?.download || 0, 
          name: '已下载',
          itemStyle: { color: '#67C23A' }
        },
        { 
          value: userStore.userInfo?.remaining || 0, 
          name: '剩余流量',
          itemStyle: { color: '#E6A23C' }
        }
      ],
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      }
    }
  ]
}))

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

const refreshServers = async () => {
  serversLoading.value = true
  try {
    servers.value = await api.getServers()
  } catch (error: any) {
    ElMessage.error('获取服务器列表失败')
  } finally {
    serversLoading.value = false
  }
}

const toggleServer = async (id: number) => {
  try {
    await api.toggleServer(id)
    await refreshServers()
  } catch (error: any) {
    ElMessage.error('切换服务器状态失败')
  }
}

const getPingClass = (ping: number) => {
  if (ping === 0) return 'ping-unknown'
  if (ping < 100) return 'ping-good'
  if (ping < 300) return 'ping-medium'
  return 'ping-bad'
}

const updateData = async () => {
  try {
    await connectionStore.updateStatus()
    if (userStore.isLoggedIn) {
      const userStatus = await api.getUserStatus()
      userStore.setUserInfo(userStatus)
    }
  } catch (error) {
    console.error('更新数据失败:', error)
  }
}

onMounted(async () => {
  await refreshServers()
  await updateData()
  
  // 每5秒更新一次数据
  updateTimer = window.setInterval(updateData, 5000)
})

onUnmounted(() => {
  if (updateTimer) {
    clearInterval(updateTimer)
  }
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.stats-row {
  margin-bottom: 20px;
}

.stats-card {
  height: 100px;
}

.stats-content {
  display: flex;
  align-items: center;
  height: 100%;
}

.stats-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
  font-size: 24px;
  color: white;
}

.user-icon { background: #409EFF; }
.upload-icon { background: #67C23A; }
.download-icon { background: #E6A23C; }
.remaining-icon { background: #F56C6C; }

.stats-info h3 {
  margin: 0 0 5px 0;
  font-size: 18px;
  font-weight: bold;
  color: #303133;
}

.stats-info p {
  margin: 0;
  font-size: 14px;
  color: #909399;
}

.content-row {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.connection-status {
  text-align: center;
}

.status-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}

.status-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  margin-right: 8px;
}

.status-dot.connected {
  background: #67C23A;
  box-shadow: 0 0 10px rgba(103, 194, 58, 0.5);
}

.status-dot.disconnected {
  background: #F56C6C;
}

.status-text {
  font-size: 16px;
  font-weight: bold;
}

.traffic-stats {
  display: flex;
  justify-content: space-around;
}

.traffic-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
}

.servers-summary {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #EBEEF5;
}

.summary-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
}

.server-list {
  max-height: 200px;
  overflow-y: auto;
}

.server-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #F5F7FA;
}

.server-item:last-child {
  border-bottom: none;
}

.server-info {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.server-name {
  font-weight: bold;
  color: #303133;
}

.server-type {
  font-size: 12px;
  color: #909399;
}

.server-status {
  display: flex;
  align-items: center;
  gap: 10px;
}

.ping {
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 3px;
  color: white;
}

.ping-good { background: #67C23A; }
.ping-medium { background: #E6A23C; }
.ping-bad { background: #F56C6C; }
.ping-unknown { background: #909399; }

.chart-row {
  margin-bottom: 20px;
}

.chart-container {
  width: 100%;
  height: 300px;
}
</style>