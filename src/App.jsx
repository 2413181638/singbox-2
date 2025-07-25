import React, { useState } from 'react';

function App() {
  const [token, setToken] = useState('');
  const [config, setConfig] = useState('');
  const [status, setStatus] = useState('未运行');

  const fetchConfig = async () => {
    const res = await window.api.xboardFetch(token);
    if (res.error) {
      setConfig('获取失败: ' + res.error);
    } else {
      setConfig(JSON.stringify(res, null, 2));
    }
  };

  const startSingbox = async () => {
    const res = await window.api.singboxStart(config);
    if (res.error) setStatus('启动失败: ' + res.error);
    else setStatus('运行中');
  };

  const stopSingbox = async () => {
    const res = await window.api.singboxStop();
    if (res.error) setStatus('停止失败: ' + res.error);
    else setStatus('未运行');
  };

  return (
    <div style={{ padding: 24 }}>
      <h2>Singbox XBoard 客户端</h2>
      <div>
        <input
          type="text"
          placeholder="输入 xboard token"
          value={token}
          onChange={e => setToken(e.target.value)}
          style={{ width: 300 }}
        />
        <button onClick={fetchConfig}>获取配置</button>
      </div>
      <div style={{ marginTop: 16 }}>
        <textarea
          rows={12}
          cols={60}
          value={config}
          onChange={e => setConfig(e.target.value)}
        />
      </div>
      <div style={{ marginTop: 16 }}>
        <button onClick={startSingbox}>启动 sing-box</button>
        <button onClick={stopSingbox} style={{ marginLeft: 8 }}>停止 sing-box</button>
        <span style={{ marginLeft: 16 }}>状态: {status}</span>
      </div>
    </div>
  );
}

export default App;