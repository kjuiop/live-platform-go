package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kjuiop/live-platform-go/config"
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
	group.POST("/", r.CreateRoom)
}

func (r *RoomHandler) successResponse(c *gin.Context, statusCode int, data interface{}) {

	c.JSON(statusCode, models.APIResponse{
		ErrorCode: models.NoError,
		Message:   models.GetCustomMessage(models.NoError),
		Result:    data,
	})
}

func (r *RoomHandler) failResponse(c *gin.Context, statusCode, errorCode int, err error) {
	logMessage := models.GetCustomErrMessage(errorCode, err.Error())
	c.Errors = append(c.Errors, &gin.Error{
		Err:  errors.New(logMessage),
		Type: gin.ErrorTypePrivate,
	})

	c.JSON(statusCode, models.APIResponse{
		ErrorCode: errorCode,
		Message:   models.GetCustomMessage(errorCode),
	})
}

func (r *RoomHandler) CreateRoom(c *gin.Context) {
	req := form.RoomRequest{}
	if err := c.ShouldBind(&req); err != nil {
		r.failResponse(c, http.StatusBadRequest, models.ErrParsing, fmt.Errorf("CreateRoom failed to parse request: %w", err))
		return
	}

	roomInfo := domain.NewRoomInfo(req, r.cfg.Prefix)
	if err := r.service.CreateChatRoom(c, *roomInfo); err != nil {
		r.failResponse(c, http.StatusInternalServerError, models.ErrRedisHMSETError, fmt.Errorf("CreateRoom HMSET err : %w", err))
		return
	}

	if err := r.service.RegisterRoomId(c, *roomInfo); err != nil {
		r.failResponse(c, http.StatusInternalServerError, models.ErrRedisHMSETError, fmt.Errorf("register room id HMSET err : %w", err))
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
