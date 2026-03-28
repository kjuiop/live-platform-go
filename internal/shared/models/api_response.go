package models

type APIResponse struct {
	ErrorCode int         `json:"error_code"`
	Message   string      `json:"message"`
	Result    interface{} `json:"result"`
}
