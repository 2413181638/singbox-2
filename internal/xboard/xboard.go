package xboard

import (
    "encoding/base64"
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "time"
)

// Profile represents a single outbound profile retrieved from XBoard.
// The structure covers only the information we need to generate sing-box config.
// When new protocol fields are added in the future, simply extend this struct.
//
// NOTE: This structure is intentionally kept protocol-agnostic; the Config
// generator will map it onto the exact sing-box outbound definitions.
//
// All fields are public so that JSON decoding works out-of-the-box.
// If the XBoard API you are using differs, adapt the struct / decoder here.
//
// 最简化示例（仅供参考）:
// {
//   "id": "node-uuid",
//   "type": "hysteria2",
//   "address": "h2.example.com",
//   "port": 443,
//   "password": "hunter2",
//   "sni": "bing.com",
//   "alpn": ["h3"],
//   "tls": true
// }
//
// 欢迎根据面板实际返回修改。
//
// --------------------------------------------------------------------------------
type Profile struct {
    ID      string   `json:"id"`
    Type    string   `json:"type"` // e.g. hysteria2, anytls, vless-reality

    Address string   `json:"address"`
    Port    int      `json:"port"`

    // Optional fields
    Password string   `json:"password,omitempty"`
    Flow     string   `json:"flow,omitempty"`
    SNI      string   `json:"sni,omitempty"`
    ALPN     []string `json:"alpn,omitempty"`
    TLS      bool     `json:"tls,omitempty"`
}

// FetchSubscription downloads the subscription from given URL, supports base64 encoded payloads.
func FetchSubscription(url string) ([]Profile, error) {
    client := &http.Client{Timeout: 15 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return nil, errors.New("unexpected HTTP status: " + resp.Status)
    }

    raw, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // Some panels base64-wrap JSON to avoid CDN corruption.
    if json.Valid(raw) {
        // already pure JSON list/object
        return decodeProfiles(raw)
    }

    // Attempt base64 decode.
    decoded, err := base64.StdEncoding.DecodeString(string(raw))
    if err != nil {
        return nil, errors.New("payload is neither valid JSON nor base64-encoded JSON")
    }

    return decodeProfiles(decoded)
}

func decodeProfiles(data []byte) ([]Profile, error) {
    var profiles []Profile
    if err := json.Unmarshal(data, &profiles); err != nil {
        return nil, err
    }
    return profiles, nil
}