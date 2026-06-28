package handler

import (
	"net/http"
	"strconv"

	"web3proof/backend/internal/pkg/response"
	"web3proof/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type WorkHandler struct {
	svc *service.WorkService
}

func NewWorkHandler(svc *service.WorkService) *WorkHandler {
	return &WorkHandler{svc: svc}
}

type createWorkRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	CategoryID  *uint  `json:"category_id"`
	ExternalURL string `json:"external_url"`
	Visibility  string `json:"visibility"`
}

func (h *WorkHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req createWorkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid request")
		return
	}
	work, err := h.svc.Create(userID.(uint64), service.CreateWorkInput{
		Title: req.Title, Description: req.Description, CategoryID: req.CategoryID, ExternalURL: req.ExternalURL, Visibility: req.Visibility,
	})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, work)
}

func (h *WorkHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	works, err := h.svc.ListByUser(userID.(uint64))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, works)
}

func (h *WorkHandler) Get(c *gin.Context) {
	userID, _ := c.Get("user_id")
	workID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid work id")
		return
	}
	work, err := h.svc.GetByID(userID.(uint64), workID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 40401, "work not found")
		return
	}
	files, _ := h.svc.ListFiles(work.ID)
	evidences, _ := h.svc.ListEvidenceByWork(work.ID)
	certificates, _ := h.svc.ListCertificatesByWork(work.ID)
	response.OK(c, gin.H{"work": work, "files": files, "evidences": evidences, "certificates": certificates})
}

func (h *WorkHandler) ListEvidence(c *gin.Context) {
	userID, _ := c.Get("user_id")
	records, err := h.svc.ListEvidenceByUser(userID.(uint64))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, records)
}

func (h *WorkHandler) ListCertificates(c *gin.Context) {
	userID, _ := c.Get("user_id")
	certs, err := h.svc.ListCertificatesByUser(userID.(uint64))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, certs)
}

func (h *WorkHandler) UploadFile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	workID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid work id")
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "file required")
		return
	}
	defer file.Close()
	record, err := h.svc.UploadFile(userID.(uint64), workID, file, header)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, record)
}

func (h *WorkHandler) PrepareEvidence(c *gin.Context) {
	userID, _ := c.Get("user_id")
	workID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid work id")
		return
	}
	result, err := h.svc.PrepareEvidence(userID.(uint64), workID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, err.Error())
		return
	}
	response.OK(c, result)
}

type confirmEvidenceRequest struct {
	TxHash          string `json:"tx_hash" binding:"required"`
	ChainEvidenceID uint64 `json:"chain_evidence_id"`
}

func (h *WorkHandler) ConfirmEvidence(c *gin.Context) {
	userID, _ := c.Get("user_id")
	workID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid work id")
		return
	}
	var req confirmEvidenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid request")
		return
	}
	record, err := h.svc.ConfirmEvidence(userID.(uint64), workID, req.TxHash, req.ChainEvidenceID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, record)
}

func (h *WorkHandler) GenerateCertificate(c *gin.Context) {
	var req struct {
		EvidenceID uint64 `json:"evidence_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid request")
		return
	}
	cert, err := h.svc.GenerateCertificate(req.EvidenceID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, err.Error())
		return
	}
	response.OK(c, cert)
}
