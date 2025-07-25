<template>
  <div class="logs">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>连接日志</span>
          <el-button @click="refreshLogs" :loading="loading">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>
      
      <el-table :data="logs" style="width: 100%" v-loading="loading">
        <el-table-column prop="start_time" label="开始时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.start_time) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="end_time" label="结束时间" width="180">
          <template #default="{ row }">
            {{ row.end_time ? formatTime(row.end_time) : '进行中' }}
          </template>
        </el-table-column>
        
        <el-table-column prop="server_id" label="服务器" width="150">
          <template #default="{ row }">
            <el-tag size="small">服务器 #{{ row.server_id }}</el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="upload" label="上传" width="120">
          <template #default="{ row }">
            {{ formatBytes(row.upload) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="download" label="下载" width="120">
          <template #default="{ row }">
            {{ formatBytes(row.download) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="success" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.success ? 'success' : 'danger'" size="small">
              {{ row.success ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="error" label="错误信息" min-width="200">
          <template #default="{ row }">
            <span v-if="row.error" class="error-text">{{ row.error }}</span>
            <span v-else class="success-text">无错误</span>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="table-footer">
        <div class="log-stats">
          <span>总记录: {{ logs.length }}</span>
          <span>成功: {{ successLogs.length }}</span>
          <span>失败: {{ failedLogs.length }}</span>
        </div>
        
        <div class="actions">
          <el-button @click="clearLogs" type="danger" size="small">
            清空日志
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api, formatBytes } from '@/utils/api'

const logs = ref<any[]>([])
const loading = ref(false)

const successLogs = computed(() => logs.value.filter(log => log.success))
const failedLogs = computed(() => logs.value.filter(log => !log.success))

const formatTime = (timestamp: string) => {
  return new Date(timestamp).toLocaleString('zh-CN')
}

const refreshLogs = async () => {
  loading.value = true
  try {
    logs.value = await api.getConnectionLogs()
  } catch (error: any) {
    ElMessage.error('获取连接日志失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const clearLogs = async () => {
  try {
    await ElMessageBox.confirm('确定要清空所有连接日志吗？', '确认操作', {
      type: 'warning'
    })
    
    // 这里应该调用后端API清空日志，暂时只清空前端显示
    logs.value = []
    ElMessage.success('日志已清空')
  } catch {
    // 用户取消操作
  }
}

onMounted(() => {
  refreshLogs()
})
</script>

<style scoped>
.logs {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.error-text {
  color: #F56C6C;
}

.success-text {
  color: #67C23A;
}

.table-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #EBEEF5;
}

.log-stats {
  display: flex;
  gap: 20px;
  font-size: 14px;
  color: #606266;
}

.actions {
  display: flex;
  gap: 10px;
}
</style>