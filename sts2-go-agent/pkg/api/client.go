package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/niuxh/sts2-go-agent/pkg/models"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetState 获取当前游戏状态
func (c *Client) GetState() (*models.GameState, error) {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/state", c.BaseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp models.Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Ok {
		return nil, fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	var state models.GameState
	if err := json.Unmarshal(apiResp.Data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// SendAction 发送动作到游戏
func (c *Client) SendAction(req *models.ActionRequest) (*models.ActionResponse, error) {
	body, _ := json.Marshal(req)
	resp, err := c.HTTPClient.Post(fmt.Sprintf("%s/action", c.BaseURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp models.Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Ok {
		return nil, fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	var actionResp models.ActionResponse
	if err := json.Unmarshal(apiResp.Data, &actionResp); err != nil {
		return nil, err
	}

	return &actionResp, nil
}
