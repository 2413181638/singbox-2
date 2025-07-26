package com.singbox.xboard;

import android.content.Intent;
import android.net.VpnService;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;
import android.widget.Toast;
import androidx.appcompat.app.AppCompatActivity;
import androidx.lifecycle.ViewModelProvider;
import com.google.android.material.textfield.TextInputEditText;
import com.singbox.xboard.databinding.ActivityMainBinding;
import com.singbox.xboard.viewmodel.MainViewModel;

public class MainActivity extends AppCompatActivity {
    private static final int VPN_REQUEST_CODE = 100;
    
    private ActivityMainBinding binding;
    private MainViewModel viewModel;
    
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        
        // 初始化视图绑定
        binding = ActivityMainBinding.inflate(getLayoutInflater());
        setContentView(binding.getRoot());
        
        // 初始化 ViewModel
        viewModel = new ViewModelProvider(this).get(MainViewModel.class);
        
        // 设置监听器
        setupListeners();
        
        // 观察数据变化
        observeViewModel();
        
        // 请求 VPN 权限
        checkVpnPermission();
    }
    
    private void setupListeners() {
        // 连接/断开按钮
        binding.btnConnect.setOnClickListener(v -> {
            if (viewModel.isConnected().getValue()) {
                disconnectVpn();
            } else {
                connectVpn();
            }
        });
        
        // 更新订阅按钮
        binding.btnUpdateSubscription.setOnClickListener(v -> {
            String url = binding.etSubscriptionUrl.getText().toString().trim();
            if (url.isEmpty()) {
                Toast.makeText(this, "请输入订阅地址", Toast.LENGTH_SHORT).show();
                return;
            }
            viewModel.updateSubscription(url);
        });
        
        // 刷新订阅按钮
        binding.btnRefreshSubscription.setOnClickListener(v -> {
            viewModel.refreshSubscription();
        });
        
        // 节点选择
        binding.spinnerNodes.setOnItemSelectedListener(new android.widget.AdapterView.OnItemSelectedListener() {
            @Override
            public void onItemSelected(android.widget.AdapterView<?> parent, View view, int position, long id) {
                // TODO: 实现节点选择
            }
            
            @Override
            public void onNothingSelected(android.widget.AdapterView<?> parent) {
            }
        });
    }
    
    private void observeViewModel() {
        // 连接状态
        viewModel.isConnected().observe(this, connected -> {
            binding.btnConnect.setText(connected ? "断开连接" : "连接");
            binding.tvStatus.setText(connected ? "已连接" : "未连接");
        });
        
        // 流量统计
        viewModel.getTrafficStats().observe(this, stats -> {
            binding.tvUpload.setText("上传: " + formatBytes(stats.getUpload()));
            binding.tvDownload.setText("下载: " + formatBytes(stats.getDownload()));
        });
        
        // 用户信息
        viewModel.getUserInfo().observe(this, userInfo -> {
            if (userInfo != null) {
                binding.tvUserEmail.setText(userInfo.getEmail());
                binding.tvUserTraffic.setText(
                    String.format("流量: %s / %s", 
                        formatBytes(userInfo.getUsed()), 
                        formatBytes(userInfo.getTotal()))
                );
            }
        });
        
        // 节点列表
        viewModel.getNodes().observe(this, nodes -> {
            // TODO: 更新节点选择器
        });
        
        // 错误信息
        viewModel.getError().observe(this, error -> {
            if (error != null && !error.isEmpty()) {
                Toast.makeText(this, error, Toast.LENGTH_LONG).show();
            }
        });
    }
    
    private void checkVpnPermission() {
        Intent intent = VpnService.prepare(this);
        if (intent != null) {
            startActivityForResult(intent, VPN_REQUEST_CODE);
        }
    }
    
    private void connectVpn() {
        Intent intent = VpnService.prepare(this);
        if (intent != null) {
            startActivityForResult(intent, VPN_REQUEST_CODE);
        } else {
            startVpnService();
        }
    }
    
    private void disconnectVpn() {
        viewModel.disconnect();
    }
    
    private void startVpnService() {
        Intent intent = new Intent(this, com.singbox.xboard.service.VpnService.class);
        intent.setAction(com.singbox.xboard.service.VpnService.ACTION_CONNECT);
        startService(intent);
        viewModel.connect();
    }
    
    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);
        if (requestCode == VPN_REQUEST_CODE) {
            if (resultCode == RESULT_OK) {
                startVpnService();
            } else {
                Toast.makeText(this, "需要 VPN 权限才能使用", Toast.LENGTH_SHORT).show();
            }
        }
    }
    
    private String formatBytes(long bytes) {
        if (bytes < 1024) return bytes + " B";
        if (bytes < 1024 * 1024) return String.format("%.2f KB", bytes / 1024.0);
        if (bytes < 1024 * 1024 * 1024) return String.format("%.2f MB", bytes / (1024.0 * 1024));
        return String.format("%.2f GB", bytes / (1024.0 * 1024 * 1024));
    }
    
    @Override
    protected void onDestroy() {
        super.onDestroy();
        binding = null;
    }
}