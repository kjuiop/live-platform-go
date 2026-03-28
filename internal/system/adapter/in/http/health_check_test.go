package http

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantCode   int
		wantStatus string
	}{
		{
			name:       "GET 요청 성공",
			method:     http.MethodGet,
			wantCode:   http.StatusOK,
			wantStatus: "ok",
		},
		{
			name:     "POST 요청 실패 - 허용되지 않는 메서드",
			method:   http.MethodPost,
			wantCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, testClient.srv.URL+"/api/v1/system/health", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Errorf("failed to close response body: %v", err)
				}
			}()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("status code: got %d, want %d", resp.StatusCode, tt.wantCode)
			}

			if tt.wantStatus == "" {
				return
			}

			var apiResp struct {
				Result map[string]string `json:"result"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}
			body := apiResp.Result

			if body["status"] != tt.wantStatus {
				t.Errorf("status: got '%s', want '%s'", body["status"], tt.wantStatus)
			}
			if body["version"] == "" {
				t.Error("version: should not be empty")
			}
			if body["git_hash"] == "" {
				t.Error("git_hash: should not be empty")
			}
			if body["uptime"] == "" {
				t.Error("uptime: should not be empty")
			}
			if body["started_at"] == "" {
				t.Error("started_at: should not be empty")
			}
		})
	}
}
