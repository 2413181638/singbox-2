<template>
  <div class="servers">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>服务器管理</span>
          <div class="header-actions">
            <el-button @click="syncSubscription" :loading="syncing" type="primary">
              <el-icon><Refresh /></el-icon>
              同步订阅
            </el-button>
            <el-button @click="pingAllServers" :loading="pinging">
              <el-icon><Connection /></el-icon>
              测试全部
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table :data="servers" style="width: 100%" v-loading="loading">
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="name" label="服务器名称" min-width="200">
          <template #default="{ row }">
            <div class="server-name">
              <span>{{ row.name }}</span>
              <el-tag :type="getServerTypeColor(row.server_type)" size="small">
                {{ row.server_type.toUpperCase() }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="host" label="地址" min-width="150">
          <template #default="{ row }">
            <span>{{ row.host }}:{{ row.port }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="ping" label="延迟" width="100" sortable>
          <template #default="{ row }">
            <el-tag :type="getPingType(row.ping)" size="small">
              {{ row.ping === 0 ? '未测试' : row.ping + 'ms' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="is_active" label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_active"
              @change="toggleServer(row.id)"
              :loading="row.loading"
            />
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button-group>
              <el-button 
                size="small" 
                @click="pingServer(row)"
                :loading="row.pinging"
              >
                <el-icon><Connection /></el-icon>
                测试
              </el-button>
              <el-button 
                size="small" 
                type="info"
                @click="showServerDetails(row)"
              >
                <el-icon><View /></el-icon>
                详情
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="table-footer">
        <div class="server-stats">
          <span>总计: {{ servers.length }} 个服务器</span>
          <span>已启用: {{ activeServers.length }} 个</span>
          <span>平均延迟: {{ averagePing }}ms</span>
        </div>
        
        <div class="batch-actions">
          <el-button @click="enableAllServers" size="small">启用全部</el-button>
          <el-button @click="disableAllServers" size="small">禁用全部</el-button>
        </div>
      </div>
    </el-card>

    <!-- 服务器详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="服务器详情"
      width="600px"
    >
      <div v-if="selectedServer" class="server-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="服务器名称">
            {{ selectedServer.name }}
          </el-descriptions-item>
          <el-descriptions-item label="服务器类型">
            <el-tag :type="getServerTypeColor(selectedServer.server_type)">
              {{ selectedServer.server_type.toUpperCase() }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="服务器地址">
            {{ selectedServer.host }}
          </el-descriptions-item>
          <el-descriptions-item label="端口">
            {{ selectedServer.port }}
          </el-descriptions-item>
          <el-descriptions-item label="延迟">
            <el-tag :type="getPingType(selectedServer.ping)">
              {{ selectedServer.ping === 0 ? '未测试' : selectedServer.ping + 'ms' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="selectedServer.is_active ? 'success' : 'danger'">
              {{ selectedServer.is_active ? '已启用' : '已禁用' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
        
        <div class="server-config">
          <h4>服务器配置</h4>
          <el-input
            v-model="serverConfigText"
            type="textarea"
            :rows="10"
            readonly
            placeholder="配置信息"
          />
        </div>
      </div>
      
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
        <el-button 
          type="primary" 
          @click="pingServer(selectedServer)"
          :loading="selectedServer?.pinging"
        >
          测试连接
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api } from '@/utils/api'

const servers = ref<any[]>([])
const loading = ref(false)
const syncing = ref(false)
const pinging = ref(false)
const detailDialogVisible = ref(false)
const selectedServer = ref<any>(null)

const activeServers = computed(() => servers.value.filter(s => s.is_active))
const averagePing = computed(() => {
  const pings = servers.value.filter(s => s.ping > 0).map(s => s.ping)
  return pings.length > 0 ? Math.round(pings.reduce((a, b) => a + b, 0) / pings.length) : 0
})

const serverConfigText = computed(() => {
  if (!selectedServer.value) return ''
  try {
    const config = JSON.parse(selectedServer.value.settings || '{}')
    return JSON.stringify(config, null, 2)
  } catch {
    return selectedServer.value.settings || ''
  }
})

const getServerTypeColor = (type: string) => {
  const colors: Record<string, string> = {
    'vmess': 'primary',
    'vless': 'success',
    'shadowsocks': 'warning',
    'trojan': 'danger',
    'hysteria': 'info'
  }
  return colors[type] || 'info'
}

const getPingType = (ping: number) => {
  if (ping === 0) return 'info'
  if (ping < 100) return 'success'
  if (ping < 300) return 'warning'
  return 'danger'
}

const loadServers = async () => {
  loading.value = true
  try {
    servers.value = await api.getServers()
    // 为每个服务器添加加载状态
    servers.value.forEach(server => {
      server.loading = false
      server.pinging = false
    })
  } catch (error: any) {
    ElMessage.error('获取服务器列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const syncSubscription = async () => {
  syncing.value = true
  try {
    await api.syncSubscription()
    await loadServers()
    ElMessage.success('订阅同步成功')
  } catch (error: any) {
    ElMessage.error('同步订阅失败: ' + error.message)
  } finally {
    syncing.value = false
  }
}

const toggleServer = async (id: number) => {
  const server = servers.value.find(s => s.id === id)
  if (!server) return
  
  server.loading = true
  try {
    await api.toggleServer(id)
    ElMessage.success(`服务器已${server.is_active ? '启用' : '禁用'}`)
  } catch (error: any) {
    // 回滚状态
    server.is_active = !server.is_active
    ElMessage.error('切换服务器状态失败: ' + error.message)
  } finally {
    server.loading = false
  }
}

const pingServer = async (server: any) => {
  server.pinging = true
  try {
    const ping = await api.pingServer(server.id)
    server.ping = ping
    ElMessage.success(`延迟测试完成: ${ping}ms`)
  } catch (error: any) {
    ElMessage.error('延迟测试失败: ' + error.message)
  } finally {
    server.pinging = false
  }
}

const pingAllServers = async () => {
  pinging.value = true
  try {
    const promises = servers.value.map(server => {
      server.pinging = true
      return api.pingServer(server.id)
        .then(ping => {
          server.ping = ping
        })
        .catch(() => {
          server.ping = 0
        })
        .finally(() => {
          server.pinging = false
        })
    })
    
    await Promise.all(promises)
    ElMessage.success('全部延迟测试完成')
  } catch (error: any) {
    ElMessage.error('批量测试失败: ' + error.message)
  } finally {
    pinging.value = false
  }
}

const enableAllServers = async () => {
  try {
    await ElMessageBox.confirm('确定要启用所有服务器吗？', '确认操作')
    
    const promises = servers.value
      .filter(s => !s.is_active)
      .map(server => {
        server.loading = true
        return api.toggleServer(server.id)
          .then(() => {
            server.is_active = true
          })
          .finally(() => {
            server.loading = false
          })
      })
    
    await Promise.all(promises)
    ElMessage.success('已启用所有服务器')
  } catch {
    // 用户取消操作
  }
}

const disableAllServers = async () => {
  try {
    await ElMessageBox.confirm('确定要禁用所有服务器吗？', '确认操作')
    
    const promises = servers.value
      .filter(s => s.is_active)
      .map(server => {
        server.loading = true
        return api.toggleServer(server.id)
          .then(() => {
            server.is_active = false
          })
          .finally(() => {
            server.loading = false
          })
      })
    
    await Promise.all(promises)
    ElMessage.success('已禁用所有服务器')
  } catch {
    // 用户取消操作
  }
}

const showServerDetails = (server: any) => {
  selectedServer.value = server
  detailDialogVisible.value = true
}

onMounted(() => {
  loadServers()
})
</script>

<style scoped>
.servers {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.server-name {
  display: flex;
  align-items: center;
  gap: 10px;
}

.table-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #EBEEF5;
}

.server-stats {
  display: flex;
  gap: 20px;
  font-size: 14px;
  color: #606266;
}

.batch-actions {
  display: flex;
  gap: 10px;
}

.server-details {
  padding: 20px 0;
}

.server-config {
  margin-top: 20px;
}

.server-config h4 {
  margin: 0 0 10px 0;
  color: #303133;
}
</style>