const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('api', {
  singboxStart: (config) => ipcRenderer.invoke('singbox-start', config),
  singboxStop: () => ipcRenderer.invoke('singbox-stop'),
  xboardFetch: (token) => ipcRenderer.invoke('xboard-fetch', token)
});