package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kjuiop/live-platform-go/config"
	"github.com/kjuiop/live-platform-go/internal/room/adapter/in/http/errror"
	"github.com/kjuiop/live-platform-go/internal/room/adapter/in/http/form"
	"github.com/kjuiop/live-platform-go/internal/room/domain"
	roomin "github.com/kjuiop/live-platform-go/internal/room/port/in"
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
		ErrorCode: models.NoError,
		Message:   models.GetCustomMessage(models.NoError),
		Result:    data,
	})
}

func (r *RoomHandler) failedResponse(c *gin.Context, err error) {
	m := errror.GetMapping(err)
	logMessage := fmt.Sprintf("%s, err: %s", m.Msg, err.Error())
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
		r.failedResponse(c, errror.ErrInvalidRequest)
		return
	}

	roomInfo := domain.NewRoomInfo(req, r.cfg.Prefix)
	if err := r.service.CreateChatRoom(ctx, *roomInfo); err != nil {
		r.failedResponse(c, errror.ErrRedisSave)
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
		r.failedResponse(c, errror.ErrInvalidRequest)
		return
	}

	ctx := c.Request.Context()
	if err := r.service.DeleteChatRoom(ctx, roomId); err != nil {
		r.failedResponse(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
