package handler

import (
	"net/http"
	"strconv"

	"web3proof/backend/internal/pkg/response"
	"web3proof/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type CertificationHandler struct {
	svc *service.CertificationService
}

func NewCertificationHandler(svc *service.CertificationService) *CertificationHandler {
	return &CertificationHandler{svc: svc}
}

func (h *CertificationHandler) Apply(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		WorkID        uint64 `json:"work_id" binding:"required"`
		EvidenceID    uint64 `json:"evidence_id" binding:"required"`
		MaterialsDesc string `json:"materials_desc"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid request")
		return
	}
	app, err := h.svc.Apply(userID.(uint64), service.ApplyInput{WorkID: req.WorkID, EvidenceID: req.EvidenceID, MaterialsDesc: req.MaterialsDesc})
	if err != nil {
		response.Fail(c, http.StatusConflict, 40901, err.Error())
		return
	}
	response.OK(c, app)
}

func (h *CertificationHandler) ListMy(c *gin.Context) {
	userID, _ := c.Get("user_id")
	apps, err := h.svc.ListMy(userID.(uint64))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, apps)
}

func (h *CertificationHandler) ListPending(c *gin.Context) {
	apps, err := h.svc.ListPending()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, apps)
}

func (h *CertificationHandler) Review(c *gin.Context) {
	reviewerID, _ := c.Get("user_id")
	appID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Status string `json:"status" binding:"required"`
		Note   string `json:"review_note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid request")
		return
	}
	app, err := h.svc.Review(appID, reviewerID.(uint64), req.Status, req.Note)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, err.Error())
		return
	}
	response.OK(c, app)
}

func (h *CertificationHandler) PrepareMint(c *gin.Context) {
	appID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	result, err := h.svc.PrepareMint(appID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *CertificationHandler) ConfirmMint(c *gin.Context) {
	appID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		TxHash   string `json:"tx_hash" binding:"required"`
		TokenID  uint64 `json:"token_id"`
		TokenURI string `json:"token_uri" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid request")
		return
	}
	record, err := h.svc.ConfirmMint(appID, req.TxHash, req.TokenID, req.TokenURI)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, record)
}
