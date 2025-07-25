import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:yaml/yaml.dart';

import '../models/singbox_outbound.dart';

class XBoardProvider {
  const XBoardProvider();

  /// 拉取并解析订阅，返回可直接写入 sing-box 配置的 Outbound 列表
  Future<List<SingboxOutbound>> fetchNodes(String subUrl) async {
    final resp = await http.get(Uri.parse(subUrl));
    if (resp.statusCode != 200) {
      throw Exception('xboard 返回 ${resp.statusCode}');
    }

    final ct = resp.headers['content-type'] ?? '';
    final body = resp.body.trim();

    // --- Clash / Singbox YAML ---
    if (ct.contains('yaml') || body.startsWith('proxies:')) {
      return _parseClashYaml(body);
    }

    // --- 多行或整体 base64 Encode ---
    final lines = body.contains('\n')
        ? body.split('\n').where((e) => e.isNotEmpty).toList()
        : utf8.decode(base64.decode(body)).split('\n');

    return _parseUris(lines);
  }

  // ----------------- private helpers -----------------
  List<SingboxOutbound> _parseClashYaml(String yamlStr) {
    final y = loadYaml(yamlStr) as YamlMap;
    final proxies = (y['proxies'] ?? []) as YamlList;
    return proxies.map((p) => SingboxOutbound.fromClash(p)).toList();
  }

  List<SingboxOutbound> _parseUris(List<String> lines) {
    final out = <SingboxOutbound>[];
    for (final l in lines) {
      if (l.startsWith('#') || l.trim().isEmpty) continue;
      final uri = Uri.parse(l.trim());
      switch (uri.scheme) {
        case 'ss':
          out.add(SingboxOutbound.fromShadowsocks(uri));
          break;
        case 'vmess':
          out.add(SingboxOutbound.fromVmess(uri));
          break;
        case 'trojan':
          out.add(SingboxOutbound.fromTrojan(uri));
          break;
        default:
          // TODO: 更多协议
          break;
      }
    }
    return out;
  }
}