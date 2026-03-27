package http

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kjuiop/live-platform-go/internal/system/domain"
)

var testClient *TestClient

type TestClient struct {
	srv           *httptest.Server
	systemHandler *SystemHandler
}

type mockSystemService struct {
	status domain.HealthStatus
}

func (m *mockSystemService) HealthCheck(_ context.Context) domain.HealthStatus {
	return m.status
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	svc := &mockSystemService{
		status: domain.HealthStatus{
			Status:    "ok",
			Version:   "1.0.0",
			GitHash:   "abc1234",
			Uptime:    5 * time.Second,
			StartedAt: time.Now().Add(-5 * time.Second),
		},
	}

	testClient = &TestClient{
		systemHandler: NewSystemHandler(svc),
	}
	testClient.srv = setupTestServer()

	exitCode := m.Run()
	testClient.srv.Close()
	os.Exit(exitCode)
}

func setupTestServer() *httptest.Server {
	r := gin.New()
	r.HandleMethodNotAllowed = true
	testClient.systemHandler.RegisterRoutes(r.Group("/api/v1/system"))
	return httptest.NewServer(r)
}
