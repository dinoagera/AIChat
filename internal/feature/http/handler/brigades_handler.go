package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dinoagera/AIChat/internal/domain"
	"github.com/dinoagera/AIChat/pkg/messages"
	"github.com/gin-gonic/gin"
)

type BrigadesHandler struct {
	log            *slog.Logger
	brigadeService BrigadeService
}

func NewBrigadesHandler(log *slog.Logger, brigadeService BrigadeService) *BrigadesHandler {
	return &BrigadesHandler{
		log:            log,
		brigadeService: brigadeService,
	}
}
func (b *BrigadesHandler) AddBrigade(c *gin.Context) {
	var req AddBrigadeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.log.Info("failed to decode json req", "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{Message: messages.MsgUncorrectDatesBrigdaes})
		return
	}
	if ok := req.Validate(); !ok {
		b.log.Info("invalid dates while add brigade")
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{Message: messages.MsgUncorrectDatesBrigdaes})
		return
	}
	if err := b.brigadeService.AddBrigade(c.Request.Context(), req.Name, req.Lat, req.Lon, req.Status); err != nil {
		if err == domain.ErrBrigadeAlreadyExists {
			b.log.Info("failed to add brigade", "err", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, Response{Message: messages.MsgNameBrigadeAlreadyExists})
			return
		}
		b.log.Info("failed to add brigade", "err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, Response{Message: messages.MsgAddBrigdaesCorrect})
}
func (b *BrigadesHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		b.log.Info("failed to parse str to int64", "err", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.log.Info("failed to decode json req", "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{Message: messages.MsgStatusUncorrect})
		return
	}
	if err := b.brigadeService.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		//add to BRIGADE IS EMPTY
		b.log.Info("failed to update status", "err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, Response{Message: messages.MsgStatusUpdateCorrect})
}
func (b *BrigadesHandler) HandlerEmergency(c *gin.Context) {
	var req EmergencyRequest
}
func (b *BrigadesHandler) SetupRoutes(router *gin.Engine) {
	brigade := router.Group("/brigade")
	{
		brigade.POST("/add", b.AddBrigade)
		brigade.POST("/:id/update", b.UpdateStatus)
	}
}
