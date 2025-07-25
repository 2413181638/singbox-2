const { app, BrowserWindow, ipcMain } = require('electron');
const path = require('path');
const singbox = require('./main/singbox');
const xboard = require('./main/xboard');

function createWindow() {
  const win = new BrowserWindow({
    width: 900,
    height: 700,
    webPreferences: {
      preload: path.join(__dirname, 'main/preload.js'),
      nodeIntegration: false,
      contextIsolation: true
    }
  });
  win.loadFile('dist/index.html');
}

app.whenReady().then(() => {
  createWindow();
  app.on('activate', function () {
    if (BrowserWindow.getAllWindows().length === 0) createWindow();
  });
});

app.on('window-all-closed', function () {
  if (process.platform !== 'darwin') app.quit();
});

// IPC: 前端与 sing-box/xboard 通信
ipcMain.handle('singbox-start', async (event, config) => {
  return await singbox.start(config);
});
ipcMain.handle('singbox-stop', async () => {
  return await singbox.stop();
});
ipcMain.handle('xboard-fetch', async (event, token) => {
  return await xboard.fetchConfig(token);
});