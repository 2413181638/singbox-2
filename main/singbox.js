const { spawn } = require('child_process');
const path = require('path');
let singboxProcess = null;

function getSingboxPath() {
  // 根据平台选择二进制
  const platform = process.platform;
  if (platform === 'win32') return path.join(__dirname, '../sing-box/sing-box.exe');
  if (platform === 'darwin') return path.join(__dirname, '../sing-box/sing-box-mac');
  return path.join(__dirname, '../sing-box/sing-box-linux');
}

async function start(config) {
  if (singboxProcess) return { error: 'sing-box 已在运行' };
  const singboxPath = getSingboxPath();
  // 假设 config 已为 JSON 字符串
  singboxProcess = spawn(singboxPath, ['run', '-c', '-'], { stdio: ['pipe', 'pipe', 'pipe'] });
  singboxProcess.stdin.write(config);
  singboxProcess.stdin.end();
  singboxProcess.stdout.on('data', (data) => {
    console.log(`[sing-box] ${data}`);
  });
  singboxProcess.stderr.on('data', (data) => {
    console.error(`[sing-box] ${data}`);
  });
  singboxProcess.on('close', (code) => {
    console.log(`sing-box 进程退出，code=${code}`);
    singboxProcess = null;
  });
  return { ok: true };
}

async function stop() {
  if (!singboxProcess) return { error: 'sing-box 未运行' };
  singboxProcess.kill();
  singboxProcess = null;
  return { ok: true };
}

module.exports = { start, stop };