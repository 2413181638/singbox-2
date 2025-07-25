const axios = require('axios');

async function fetchConfig(token) {
  // 假设 xboard 面板 API 地址如下，token 为鉴权
  const url = 'https://your-xboard-panel.com/api/client/subscribe?token=' + token;
  try {
    const res = await axios.get(url);
    // 这里假设返回 sing-box 配置
    return res.data;
  } catch (e) {
    return { error: e.message };
  }
}

module.exports = { fetchConfig };