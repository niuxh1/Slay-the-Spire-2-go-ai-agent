package models

// ActionRequest 是发送给 POST /action 的请求体
type ActionRequest struct {
	Action        string                 `json:"action"`
	CardIndex     *int                   `json:"card_index,omitempty"`
	TargetIndex   *int                   `json:"target_index,omitempty"`
	OptionIndex   *int                   `json:"option_index,omitempty"`
	ClientContext map[string]interface{} `json:"client_context,omitempty"`
}

// ActionResponse 是动作执行后的响应数据
type ActionResponse struct {
	Action  string    `json:"action"`
	Status  string    `json:"status"`
	Stable  bool      `json:"stable"`
	Message string    `json:"message"`
	State   GameState `json:"state"`
}
