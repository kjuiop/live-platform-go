package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kjuiop/live-platform-go/config"
	rerror "github.com/kjuiop/live-platform-go/internal/room/adapter/in/http/error"
	"github.com/kjuiop/live-platform-go/internal/room/adapter/in/http/form"
	"github.com/kjuiop/live-platform-go/internal/room/domain"
	roomin "github.com/kjuiop/live-platform-go/internal/room/port/in"
	serrors "github.com/kjuiop/live-platform-go/internal/shared/errors"
	"github.com/kjuiop/live-platform-go/internal/shared/models"
)

type RoomHandler struct {
	cfg     config.Policy
	service roomin.RoomService
}

func NewRoomHandler(cfg config.Policy, service roomin.RoomService) *RoomHandler {
	return &RoomHandler{
		cfg:     cfg,
		service: service,
	}
}

func (r *RoomHandler) RegisterRoutes(router gin.IRouter) {
	group := router.Group("/rooms")
	group.POST("/", r.CreateChatRoom)
	group.DELETE("/:roomId", r.DeleteChatRoom)
}

func (r *RoomHandler) successResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, models.APIResponse{
		ErrorCode: serrors.NoError,
		Message:   serrors.GetCustomMessage(serrors.NoError),
		Result:    data,
	})
}

func (r *RoomHandler) failedResponse(c *gin.Context, err error, cause ...error) {
	m := rerror.GetMapping(err)
	logMessage := fmt.Sprintf("%s, err: %s", m.Msg, err.Error())
	if len(cause) > 0 && cause[0] != nil {
		logMessage = fmt.Sprintf("%s, cause: %s", logMessage, cause[0].Error())
	}
	c.Errors = append(c.Errors, &gin.Error{
		Err:  errors.New(logMessage),
		Type: gin.ErrorTypePrivate,
	})
	c.JSON(m.HttpStatus, models.APIResponse{
		ErrorCode: m.ErrorCode,
		Message:   m.Msg,
	})
}

func (r *RoomHandler) CreateChatRoom(c *gin.Context) {
	req := form.RoomRequest{}
	ctx := c.Request.Context()
	if err := c.ShouldBind(&req); err != nil {
		r.failedResponse(c, serrors.ErrInvalidRequest, err)
		return
	}

	roomInfo := domain.NewRoomInfo(req, r.cfg.Prefix)
	if err := r.service.CreateChatRoom(ctx, *roomInfo); err != nil {
		r.failedResponse(c, serrors.ErrRedisSave)
		return
	}

	roomRes := form.RoomResponse{
		RoomId:       roomInfo.RoomId,
		CustomerId:   roomInfo.CustomerId,
		ChannelKey:   roomInfo.ChannelKey,
		BroadcastKey: roomInfo.BroadcastKey,
		CreatedAt:    roomInfo.CreatedAt,
	}

	r.successResponse(c, http.StatusCreated, roomRes)
}

func (r *RoomHandler) DeleteChatRoom(c *gin.Context) {
	roomId := c.Param("roomId")
	if roomId == "" {
		r.failedResponse(c, serrors.ErrInvalidRequest)
		return
	}

	ctx := c.Request.Context()
	if err := r.service.DeleteChatRoom(ctx, roomId); err != nil {
		r.failedResponse(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
