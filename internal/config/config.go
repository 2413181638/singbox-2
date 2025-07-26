package config

import (
    "encoding/json"
    "fmt"

    "github.com/yourusername/singxclient/internal/xboard"
)

// Generate converts slice of xboard.Profile into a sing-box JSON configuration.
// At minimum we emit:
//   * a mixed inbound listening on 0.0.0.0:2080 (SOCKS+HTTP)
//   * outbounds – one per profile
//   * selector outbound to allow run-time switching (if clients/UI in future)
//
// The returned byte slice is ready to be written to `config.json`.
func Generate(profiles []xboard.Profile) ([]byte, error) {
    conf := make(map[string]any)
    conf["log"] = map[string]any{
        "level": "info",
        "output": "box.log",
    }

    conf["inbounds"] = []any{
        map[string]any{
            "type": "mixed",
            "tag": "mixed-in",
            "listen": "0.0.0.0",
            "listen_port": 2080,
            "sniff": true,
        },
    }

    var outbounds []any
    var selectorNames []string

    for i, p := range profiles {
        ob, err := makeOutbound(p, i)
        if err != nil {
            return nil, err
        }
        outbounds = append(outbounds, ob)
        selectorNames = append(selectorNames, ob.(map[string]any)["tag"].(string))
    }

    // Add selector outbound as last element.
    outbounds = append(outbounds, map[string]any{
        "type":      "selector",
        "tag":       "auto",
        "outbounds": selectorNames,
    })

    conf["outbounds"] = outbounds

    // simple DNS & route to avoid proxying LAN.
    conf["route"] = map[string]any{
        "geoip": map[string]any{
            "download_url": "https://raw.githubusercontent.com/foxcpp/geolite2updater/master/GeoLite2-Country.mmdb",
        },
        "rules": []any{
            map[string]any{
                "type":       "direct",
                "outbound":   "direct",
                "ip_cidr":    []string{"geoip:private"},
                "domain":     []string{"geosite:cn"},
                "domain_suffix": []string{"lan"},
            },
        },
    }

    return json.MarshalIndent(conf, "", "  ")
}

// makeOutbound maps an XBoard profile into a sing-box outbound object.
func makeOutbound(p xboard.Profile, idx int) (any, error) {
    tag := fmt.Sprintf("node-%d", idx)

    switch p.Type {
    case "hysteria2":
        return map[string]any{
            "type":  "hysteria2",
            "tag":   tag,
            "server": map[string]any{
                "address": p.Address,
                "port":    p.Port,
            },
            "auth_str": p.Password,
            "tls": map[string]any{
                "enabled": true,
                "server_name": p.SNI,
                "insecure":   false,
                "alpn":       p.ALPN,
            },
        }, nil
    case "anytls":
        return map[string]any{
            "type": "anytls",
            "tag":  tag,
            "server": map[string]any{
                "address": p.Address,
                "port":    p.Port,
            },
            "username": p.ID,
            "password": p.Password,
            "tls": map[string]any{
                "enabled": p.TLS,
                "server_name": p.SNI,
                "insecure": false,
            },
        }, nil
    case "vless-reality", "vless+reality", "vless":
        return map[string]any{
            "type": "vless",
            "tag":  tag,
            "server": map[string]any{
                "address": p.Address,
                "port":    p.Port,
            },
            "uuid": p.ID,
            "flow": p.Flow,
            "tls": map[string]any{
                "enabled": true,
                "server_name": p.SNI,
                "alpn": p.ALPN,
                "insecure": false,
                "reality": map[string]any{
                    "enabled": true,
                    "public_key": p.Password, // using Password as pk for example – adjust if differs
                    "short_id":  "", // TODO – fill if your panel provides
                },
            },
        }, nil
    default:
        return nil, fmt.Errorf("unsupported outbound type: %s", p.Type)
    }
}