package models

type APIResponse struct {
	ErrorCode string      `json:"error_code"`
	Message   string      `json:"message"`
	Result    interface{} `json:"result"`
}
