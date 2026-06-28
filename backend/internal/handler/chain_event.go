package handler

import (
	"net/http"
	"strconv"

	"web3proof/backend/internal/pkg/response"
	"web3proof/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ChainEventHandler struct {
	svc *service.ChainEventService
}

func NewChainEventHandler(svc *service.ChainEventService) *ChainEventHandler {
	return &ChainEventHandler{svc: svc}
}

func (h *ChainEventHandler) Sync(c *gin.Context) {
	lookback := uint64(5000)
	if raw := c.Query("lookback"); raw != "" {
		value, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, 40001, "invalid lookback")
			return
		}
		lookback = value
	}
	result, err := h.svc.SyncRecent(c.Request.Context(), lookback)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *ChainEventHandler) List(c *gin.Context) {
	limit := 50
	if raw := c.Query("limit"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, 40001, "invalid limit")
			return
		}
		limit = value
	}
	events, err := h.svc.ListRecent(limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, events)
}
