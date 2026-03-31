package handler

import (
	"log/slog"
	"net/http"

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
		//TODO:add err name  already exists
		b.log.Info("failed to add brigade", "err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, Response{Message: messages.MsgAddBrigdaesCorrect})
}
