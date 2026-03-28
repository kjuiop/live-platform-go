package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

// Recovery는 핸들러에서 발생한 panic을 복구하는 미들웨어다.
// panic 발생 시 slog.Error로 스택 트레이스를 기록하고 500을 반환한다.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// broken pipe / connection reset은 클라이언트가 연결을 끊은 것이므로
				// 500을 쓰려 해도 실패한다. 조용히 abort만 한다.
				if isBrokenPipe(err) {
					slog.Warn("broken pipe, client disconnected",
						slog.Any("error", err),
						slog.String("path", c.Request.URL.Path),
					)
					c.Abort()
					return
				}

				stack := debug.Stack()
				slog.Error("panic recovered",
					slog.Any("error", err),
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
					slog.String("stack", string(stack)),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// isBrokenPipe는 broken pipe / connection reset 에러인지 판단한다.
func isBrokenPipe(err any) bool {
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			errStr := strings.ToLower(se.Error())
			return strings.Contains(errStr, "broken pipe") ||
				strings.Contains(errStr, "connection reset by peer")
		}
	}
	return false
}
