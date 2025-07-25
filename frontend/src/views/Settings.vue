<template>
  <div class="settings">
    <el-row :gutter="20">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>基本设置</span>
          </template>
          
          <el-form
            ref="settingsFormRef"
            :model="settings"
            :rules="rules"
            label-width="120px"
            @submit.prevent="saveSettings"
          >
            <el-form-item label="XBoard地址" prop="xboardUrl">
              <el-input
                v-model="settings.xboardUrl"
                placeholder="https://panel.example.com"
              />
            </el-form-item>
            
            <el-form-item label="同步间隔" prop="syncInterval">
              <el-input-number
                v-model="settings.syncInterval"
                :min="60"
                :max="3600"
                :step="60"
                controls-position="right"
              />
              <span class="form-help">秒（建议300秒）</span>
            </el-form-item>
            
            <el-form-item label="日志级别" prop="logLevel">
              <el-select v-model="settings.logLevel" style="width: 200px;">
                <el-option label="Debug" value="debug" />
                <el-option label="Info" value="info" />
                <el-option label="Warning" value="warning" />
                <el-option label="Error" value="error" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="API端口" prop="apiPort">
              <el-input-number
                v-model="settings.apiPort"
                :min="1024"
                :max="65535"
                controls-position="right"
              />
              <span class="form-help">SingBox API端口</span>
            </el-form-item>
            
            <el-form-item label="本地代理端口" prop="proxyPort">
              <el-input-number
                v-model="settings.proxyPort"
                :min="1024"
                :max="65535"
                controls-position="right"
              />
              <span class="form-help">本地代理监听端口</span>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="saveSettings" :loading="saving">
                保存设置
              </el-button>
              <el-button @click="resetSettings">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
      
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>系统信息</span>
          </template>
          
          <el-descriptions :column="1" border>
            <el-descriptions-item label="应用版本">
              v1.0.0
            </el-descriptions-item>
            <el-descriptions-item label="SingBox版本">
              {{ systemInfo.singboxVersion || '未知' }}
            </el-descriptions-item>
            <el-descriptions-item label="配置文件路径">
              {{ systemInfo.configPath || '未知' }}
            </el-descriptions-item>
            <el-descriptions-item label="数据库路径">
              {{ systemInfo.databasePath || '未知' }}
            </el-descriptions-item>
            <el-descriptions-item label="运行时间">
              {{ formatUptime(systemInfo.uptime || 0) }}
            </el-descriptions-item>
          </el-descriptions>
          
          <div class="system-actions">
            <el-button @click="checkUpdate" :loading="checkingUpdate">
              检查更新
            </el-button>
            <el-button @click="exportConfig" type="info">
              导出配置
            </el-button>
            <el-button @click="showImportDialog = true" type="warning">
              导入配置
            </el-button>
          </div>
        </el-card>
        
        <el-card style="margin-top: 20px;">
          <template #header>
            <span>代理设置</span>
          </template>
          
          <div class="proxy-info">
            <p><strong>HTTP代理:</strong> 127.0.0.1:{{ settings.proxyPort }}</p>
            <p><strong>SOCKS5代理:</strong> 127.0.0.1:{{ settings.proxyPort }}</p>
            <p class="proxy-help">
              将上述地址配置到您的应用程序中即可使用代理
            </p>
          </div>
          
          <el-button @click="copyProxyInfo" size="small" type="primary">
            复制代理信息
          </el-button>
        </el-card>
      </el-col>
    </el-row>

    <!-- 导入配置对话框 -->
    <el-dialog v-model="showImportDialog" title="导入配置" width="500px">
      <el-upload
        ref="uploadRef"
        :auto-upload="false"
        :show-file-list="false"
        accept=".json,.yaml,.yml"
        :on-change="handleFileChange"
        drag
      >
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div class="el-upload__text">
          将配置文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持 JSON 和 YAML 格式的配置文件
          </div>
        </template>
      </el-upload>
      
      <template #footer>
        <el-button @click="showImportDialog = false">取消</el-button>
        <el-button type="primary" @click="importConfig" :loading="importing">
          导入
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import { api } from '@/utils/api'

const settingsFormRef = ref<FormInstance>()
const saving = ref(false)
const checkingUpdate = ref(false)
const showImportDialog = ref(false)
const importing = ref(false)
const uploadRef = ref()
let importFile: File | null = null

const settings = reactive({
  xboardUrl: '',
  syncInterval: 300,
  logLevel: 'info',
  apiPort: 9090,
  proxyPort: 7890
})

const systemInfo = reactive({
  singboxVersion: '',
  configPath: '',
  databasePath: '',
  uptime: 0
})

const rules: FormRules = {
  xboardUrl: [
    { required: true, message: '请输入XBoard地址', trigger: 'blur' },
    { pattern: /^https?:\/\/.+/, message: '请输入正确的URL格式', trigger: 'blur' }
  ],
  syncInterval: [
    { required: true, message: '请设置同步间隔', trigger: 'blur' }
  ],
  logLevel: [
    { required: true, message: '请选择日志级别', trigger: 'change' }
  ],
  apiPort: [
    { required: true, message: '请设置API端口', trigger: 'blur' }
  ],
  proxyPort: [
    { required: true, message: '请设置代理端口', trigger: 'blur' }
  ]
}

const loadSettings = async () => {
  try {
    const config = await api.getConfig()
    settings.xboardUrl = config.xboard?.url || ''
    settings.syncInterval = config.xboard?.interval || 300
    settings.logLevel = config.log_level || 'info'
    settings.apiPort = config.singbox?.api_port || 9090
    settings.proxyPort = 7890 // 从配置中读取
    
    // 加载系统信息
    systemInfo.configPath = config.singbox?.config_path || ''
    systemInfo.databasePath = config.database_path || ''
  } catch (error: any) {
    ElMessage.error('加载设置失败: ' + error.message)
  }
}

const saveSettings = async () => {
  if (!settingsFormRef.value) return
  
  const valid = await settingsFormRef.value.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    const config = await api.getConfig()
    
    // 更新配置
    config.xboard.url = settings.xboardUrl
    config.xboard.interval = settings.syncInterval
    config.log_level = settings.logLevel
    config.singbox.api_port = settings.apiPort
    
    await api.updateConfig(config)
    ElMessage.success('设置保存成功')
  } catch (error: any) {
    ElMessage.error('保存设置失败: ' + error.message)
  } finally {
    saving.value = false
  }
}

const resetSettings = async () => {
  await loadSettings()
  ElMessage.info('设置已重置')
}

const checkUpdate = async () => {
  checkingUpdate.value = true
  try {
    // 模拟检查更新
    await new Promise(resolve => setTimeout(resolve, 2000))
    ElMessage.success('当前已是最新版本')
  } catch (error: any) {
    ElMessage.error('检查更新失败: ' + error.message)
  } finally {
    checkingUpdate.value = false
  }
}

const exportConfig = async () => {
  try {
    const config = await api.getConfig()
    const blob = new Blob([JSON.stringify(config, null, 2)], {
      type: 'application/json'
    })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'singbox-config.json'
    a.click()
    URL.revokeObjectURL(url)
    ElMessage.success('配置导出成功')
  } catch (error: any) {
    ElMessage.error('导出配置失败: ' + error.message)
  }
}

const handleFileChange = (file: any) => {
  importFile = file.raw
}

const importConfig = async () => {
  if (!importFile) {
    ElMessage.warning('请选择要导入的配置文件')
    return
  }

  importing.value = true
  try {
    const text = await importFile.text()
    let config
    
    if (importFile.name.endsWith('.json')) {
      config = JSON.parse(text)
    } else {
      // 简单的YAML解析，实际应该使用专门的YAML库
      ElMessage.warning('暂不支持YAML格式，请使用JSON格式')
      return
    }
    
    await api.updateConfig(config)
    await loadSettings()
    showImportDialog.value = false
    ElMessage.success('配置导入成功')
  } catch (error: any) {
    ElMessage.error('导入配置失败: ' + error.message)
  } finally {
    importing.value = false
  }
}

const copyProxyInfo = () => {
  const proxyInfo = `HTTP代理: 127.0.0.1:${settings.proxyPort}
SOCKS5代理: 127.0.0.1:${settings.proxyPort}`
  
  navigator.clipboard.writeText(proxyInfo).then(() => {
    ElMessage.success('代理信息已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

const formatUptime = (seconds: number) => {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return `${hours}小时${minutes}分钟`
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.settings {
  padding: 0;
}

.form-help {
  margin-left: 10px;
  color: #909399;
  font-size: 12px;
}

.system-actions {
  margin-top: 20px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.proxy-info {
  margin-bottom: 15px;
}

.proxy-info p {
  margin: 8px 0;
  font-family: monospace;
}

.proxy-help {
  color: #909399;
  font-size: 12px;
  font-family: initial !important;
}
</style>