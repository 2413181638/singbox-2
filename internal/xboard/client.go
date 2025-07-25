package xboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"singbox-xboard-client/internal/config"
)

type Client struct {
	config     config.XBoardConfig
	httpClient *http.Client
}

type UserInfo struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	UUID     string `json:"uuid"`
	Transfer struct {
		Up        int64 `json:"up"`
		Down      int64 `json:"down"`
		Total     int64 `json:"total"`
		Remaining int64 `json:"remaining"`
	} `json:"transfer"`
	ExpiredAt int64 `json:"expired_at"`
	Status    int   `json:"status"`
}

type ServerNode struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Settings string `json:"settings"`
	Tags     string `json:"tags"`
	Rate     string `json:"rate"`
	Network  string `json:"network"`
	TLS      int    `json:"tls"`
}

type SubscriptionResponse struct {
	Data []ServerNode `json:"data"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func New(config config.XBoardConfig) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GetUserInfo(token string) (*UserInfo, error) {
	url := fmt.Sprintf("%s/api/v1/user/info", c.config.URL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.Code != 200 {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	userInfoBytes, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	var userInfo UserInfo
	if err := json.Unmarshal(userInfoBytes, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (c *Client) GetSubscription(token string) ([]ServerNode, error) {
	url := fmt.Sprintf("%s/api/v1/user/server/fetch", c.config.URL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.Code != 200 {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	nodesBytes, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	var nodes []ServerNode
	if err := json.Unmarshal(nodesBytes, &nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (c *Client) ReportTraffic(token string, upload, download int64) error {
	url := fmt.Sprintf("%s/api/v1/user/traffic", c.config.URL)
	
	data := map[string]interface{}{
		"u": upload,
		"d": download,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	if response.Code != 200 {
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

func (c *Client) Login(email, password string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/passport/auth/login", c.config.URL)
	
	data := map[string]string{
		"email":    email,
		"password": password,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if response.Code != 200 {
		return "", fmt.Errorf("登录失败: %s", response.Message)
	}

	// 提取 token
	dataMap, ok := response.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	token, ok := dataMap["auth_data"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	return token, nil
}