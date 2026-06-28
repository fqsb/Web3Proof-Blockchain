package handler

import (
	"net/http"

	"web3proof/backend/internal/pkg/response"
	"web3proof/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type VerificationHandler struct {
	svc *service.VerificationService
}

func NewVerificationHandler(svc *service.VerificationService) *VerificationHandler {
	return &VerificationHandler{svc: svc}
}

func (h *VerificationHandler) VerifyFile(c *gin.Context) {
	viewerID := optionalUserID(c)
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "file required")
		return
	}
	defer file.Close()
	report, err := h.svc.VerifyFile(viewerID, file)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, report)
}

func (h *VerificationHandler) VerifyEvidence(c *gin.Context) {
	report, err := h.svc.VerifyEvidenceNo(optionalUserID(c), c.Param("no"))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, report)
}

func (h *VerificationHandler) VerifyCertificate(c *gin.Context) {
	report, err := h.svc.VerifyCertificateNo(optionalUserID(c), c.Param("no"))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, report)
}

func (h *VerificationHandler) VerifyWallet(c *gin.Context) {
	report, err := h.svc.VerifyWallet(optionalUserID(c), c.Param("address"))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, report)
}

func (h *VerificationHandler) ListReports(c *gin.Context) {
	userID, _ := c.Get("user_id")
	reports, err := h.svc.ListReports(userID.(uint64))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, reports)
}

func optionalUserID(c *gin.Context) *uint64 {
	raw, ok := c.Get("user_id")
	if !ok {
		return nil
	}
	id, ok := raw.(uint64)
	if !ok {
		return nil
	}
	return &id
}
