const { createApp } = Vue;

createApp({
    data() {
        return {
            // 状态
            isRunning: false,
            loading: false,
            uptime: 0,
            
            // 订阅
            subscriptionUrl: '',
            lastUpdate: null,
            
            // 用户信息
            userInfo: null,
            
            // 节点
            nodes: [],
            selectedNodeId: null,
            
            // 流量统计
            stats: {
                upload: 0,
                download: 0
            },
            
            // 消息提示
            message: '',
            messageType: 'info',
            messageTimer: null,
            
            // WebSocket
            ws: null
        };
    },
    
    mounted() {
        this.init();
        this.connectWebSocket();
    },
    
    beforeUnmount() {
        if (this.ws) {
            this.ws.close();
        }
    },
    
    methods: {
        // 初始化
        async init() {
            await this.getStatus();
            await this.getSubscription();
            await this.getConfig();
        },
        
        // 获取状态
        async getStatus() {
            try {
                const response = await fetch('/api/status');
                const data = await response.json();
                
                if (data.success) {
                    this.isRunning = data.data.running;
                    this.uptime = data.data.uptime || 0;
                    this.stats = data.data.stats || { upload: 0, download: 0 };
                    this.userInfo = data.data.user || null;
                    this.lastUpdate = data.data.lastUpdate || null;
                }
            } catch (error) {
                console.error('获取状态失败:', error);
            }
        },
        
        // 获取配置
        async getConfig() {
            try {
                const response = await fetch('/api/config');
                const data = await response.json();
                
                if (data.success && data.data.subscription) {
                    this.subscriptionUrl = data.data.subscription.url || '';
                }
            } catch (error) {
                console.error('获取配置失败:', error);
            }
        },
        
        // 获取订阅信息
        async getSubscription() {
            try {
                const response = await fetch('/api/subscription');
                const data = await response.json();
                
                if (data.success) {
                    this.lastUpdate = data.data.lastUpdate;
                    this.nodes = data.data.nodes || [];
                }
            } catch (error) {
                console.error('获取订阅信息失败:', error);
            }
        },
        
        // 更新订阅
        async updateSubscription() {
            if (!this.subscriptionUrl) {
                this.showMessage('请输入订阅地址', 'error');
                return;
            }
            
            this.loading = true;
            try {
                const response = await fetch('/api/subscription', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        url: this.subscriptionUrl
                    })
                });
                
                const data = await response.json();
                
                if (data.success) {
                    this.showMessage('订阅更新成功', 'success');
                    await this.getSubscription();
                } else {
                    this.showMessage(data.error || '订阅更新失败', 'error');
                }
            } catch (error) {
                this.showMessage('订阅更新失败: ' + error.message, 'error');
            } finally {
                this.loading = false;
            }
        },
        
        // 刷新订阅
        async refreshSubscription() {
            this.loading = true;
            try {
                const response = await fetch('/api/subscription/refresh', {
                    method: 'POST'
                });
                
                const data = await response.json();
                
                if (data.success) {
                    this.showMessage('订阅刷新成功', 'success');
                    await this.getSubscription();
                } else {
                    this.showMessage(data.error || '订阅刷新失败', 'error');
                }
            } catch (error) {
                this.showMessage('订阅刷新失败: ' + error.message, 'error');
            } finally {
                this.loading = false;
            }
        },
        
        // 选择节点
        async selectNode(nodeId) {
            this.selectedNodeId = nodeId;
            
            try {
                const response = await fetch('/api/node/select', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        node_id: nodeId
                    })
                });
                
                const data = await response.json();
                
                if (data.success) {
                    this.showMessage('节点切换成功', 'success');
                } else {
                    this.showMessage(data.error || '节点切换失败', 'error');
                }
            } catch (error) {
                this.showMessage('节点切换失败: ' + error.message, 'error');
            }
        },
        
        // 切换 sing-box 状态
        async toggleSingbox() {
            if (this.isRunning) {
                await this.stopSingbox();
            } else {
                await this.startSingbox();
            }
        },
        
        // 启动 sing-box
        async startSingbox() {
            this.loading = true;
            try {
                const response = await fetch('/api/singbox/start', {
                    method: 'POST'
                });
                
                const data = await response.json();
                
                if (data.success) {
                    this.showMessage('启动成功', 'success');
                    this.isRunning = true;
                } else {
                    this.showMessage(data.error || '启动失败', 'error');
                }
            } catch (error) {
                this.showMessage('启动失败: ' + error.message, 'error');
            } finally {
                this.loading = false;
            }
        },
        
        // 停止 sing-box
        async stopSingbox() {
            this.loading = true;
            try {
                const response = await fetch('/api/singbox/stop', {
                    method: 'POST'
                });
                
                const data = await response.json();
                
                if (data.success) {
                    this.showMessage('停止成功', 'success');
                    this.isRunning = false;
                    this.uptime = 0;
                } else {
                    this.showMessage(data.error || '停止失败', 'error');
                }
            } catch (error) {
                this.showMessage('停止失败: ' + error.message, 'error');
            } finally {
                this.loading = false;
            }
        },
        
        // 重启 sing-box
        async restartSingbox() {
            this.loading = true;
            try {
                const response = await fetch('/api/singbox/restart', {
                    method: 'POST'
                });
                
                const data = await response.json();
                
                if (data.success) {
                    this.showMessage('重启成功', 'success');
                } else {
                    this.showMessage(data.error || '重启失败', 'error');
                }
            } catch (error) {
                this.showMessage('重启失败: ' + error.message, 'error');
            } finally {
                this.loading = false;
            }
        },
        
        // WebSocket 连接
        connectWebSocket() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = `${protocol}//${window.location.host}/api/ws`;
            
            this.ws = new WebSocket(wsUrl);
            
            this.ws.onopen = () => {
                console.log('WebSocket 已连接');
            };
            
            this.ws.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    if (data.type === 'status') {
                        this.isRunning = data.running;
                        this.uptime = data.uptime || 0;
                        this.stats = data.stats || { upload: 0, download: 0 };
                    }
                } catch (error) {
                    console.error('WebSocket 消息解析失败:', error);
                }
            };
            
            this.ws.onerror = (error) => {
                console.error('WebSocket 错误:', error);
            };
            
            this.ws.onclose = () => {
                console.log('WebSocket 已断开');
                // 5秒后重连
                setTimeout(() => {
                    this.connectWebSocket();
                }, 5000);
            };
        },
        
        // 显示消息
        showMessage(message, type = 'info') {
            this.message = message;
            this.messageType = type;
            
            if (this.messageTimer) {
                clearTimeout(this.messageTimer);
            }
            
            this.messageTimer = setTimeout(() => {
                this.message = '';
            }, 3000);
        },
        
        // 格式化时间
        formatUptime(seconds) {
            const hours = Math.floor(seconds / 3600);
            const minutes = Math.floor((seconds % 3600) / 60);
            const secs = Math.floor(seconds % 60);
            
            return `${hours}小时 ${minutes}分钟 ${secs}秒`;
        },
        
        // 格式化日期
        formatDate(dateString) {
            if (!dateString) return '-';
            const date = new Date(dateString);
            return date.toLocaleString('zh-CN');
        },
        
        // 格式化字节
        formatBytes(bytes) {
            if (bytes === 0) return '0 B';
            
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }
    }
}).mount('#app');