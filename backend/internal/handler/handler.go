package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"web3proof/backend/internal/middleware"
	"web3proof/backend/internal/model"
	ethcrypto "web3proof/backend/internal/pkg/crypto"
	"web3proof/backend/internal/pkg/eth"
	"web3proof/backend/internal/pkg/response"
	"web3proof/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Handler struct {
	db          *gorm.DB
	jwtSecret   string
	authService *service.AuthService
	repService  *service.ReputationService
	ethClient   *eth.EthClient
}

func New(db *gorm.DB, jwtSecret string, authService *service.AuthService, repService *service.ReputationService, ethClient *eth.EthClient) *Handler {
	return &Handler{db: db, jwtSecret: jwtSecret, authService: authService, repService: repService, ethClient: ethClient}
}

func (h *Handler) Health(c *gin.Context) {
	response.OK(c, gin.H{"status": "up", "service": "web3proof-backend"})
}

func (h *Handler) GetNonce(c *gin.Context) {
	address := strings.ToLower(c.Query("address"))
	if address == "" {
		response.Fail(c, http.StatusBadRequest, 40001, "address required")
		return
	}
	_, message, err := h.authService.GenerateNonce(c.Request.Context(), address)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, gin.H{"message": message})
}

type loginRequest struct {
	Address   string `json:"address" binding:"required"`
	Signature string `json:"signature" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid request")
		return
	}
	address := strings.ToLower(req.Address)
	if _, err := h.authService.ValidateSIWEMessage(c.Request.Context(), address, req.Message); err != nil {
		response.Fail(c, http.StatusUnauthorized, 40101, err.Error())
		return
	}
	if err := ethcrypto.VerifyPersonalSign(req.Message, req.Signature, address); err != nil {
		response.Fail(c, http.StatusUnauthorized, 40101, "invalid signature")
		return
	}
	if err := h.authService.ConsumeNonce(c.Request.Context(), address); err != nil {
		response.Fail(c, http.StatusUnauthorized, 40101, err.Error())
		return
	}
	user, err := h.authService.FindOrCreateUser(address)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	token, err := h.signToken(user)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	payload, _ := h.userPayload(user.ID)
	response.OK(c, gin.H{"token": token, "user": payload})
}

func (h *Handler) signToken(user *model.User) (string, error) {
	claims := middleware.Claims{
		UserID:        user.ID,
		WalletAddress: user.WalletAddress,
		ActiveRole:    user.ActiveRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

func (h *Handler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	payload, err := h.userPayload(userID.(uint64))
	if err != nil {
		response.Fail(c, http.StatusNotFound, 40401, "user not found")
		return
	}
	response.OK(c, payload)
}

func (h *Handler) DashboardSummary(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := userID.(uint64)
	var workCount, evidenceCount, certificateCount, applicationCount, credentialCount int64
	h.db.Model(&model.Work{}).Where("user_id = ?", id).Count(&workCount)
	h.db.Model(&model.EvidenceRecord{}).Where("user_id = ?", id).Count(&evidenceCount)
	h.db.Model(&model.Certificate{}).Where("user_id = ?", id).Count(&certificateCount)
	h.db.Model(&model.CertificationApplication{}).Where("user_id = ?", id).Count(&applicationCount)
	h.db.Model(&model.SBTCredential{}).Where("user_id = ? AND status = ?", id, "active").Count(&credentialCount)

	var recentWorks []model.Work
	h.db.Preload("Category").Where("user_id = ?", id).Order("created_at desc").Limit(5).Find(&recentWorks)
	var recentEvidences []model.EvidenceRecord
	h.db.Where("user_id = ?", id).Order("created_at desc").Limit(5).Find(&recentEvidences)
	score, _ := h.repService.GetByUserID(id)

	response.OK(c, gin.H{
		"counts": gin.H{
			"works":        workCount,
			"evidences":    evidenceCount,
			"certificates": certificateCount,
			"applications": applicationCount,
			"credentials":  credentialCount,
		},
		"reputation":       score,
		"recent_works":     recentWorks,
		"recent_evidences": recentEvidences,
	})
}

func (h *Handler) userPayload(userID uint64) (gin.H, error) {
	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	var roles []model.UserRole
	h.db.Where("user_id = ? AND enabled = ?", userID, true).Order("role_code asc").Find(&roles)
	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		roleCodes = append(roleCodes, role.RoleCode)
	}
	return gin.H{
		"id":                user.ID,
		"wallet_address":    user.WalletAddress,
		"did":               user.DID,
		"nickname":          user.Nickname,
		"avatar_url":        user.AvatarURL,
		"bio":               user.Bio,
		"email":             user.Email,
		"active_role":       user.ActiveRole,
		"roles":             roleCodes,
		"is_did_registered": user.IsDIDRegistered,
	}, nil
}

type updateProfileRequest struct {
	Nickname  *string `json:"nickname"`
	Bio       *string `json:"bio"`
	AvatarURL *string `json:"avatar_url"`
	Email     *string `json:"email"`
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid request")
		return
	}
	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		response.Fail(c, http.StatusNotFound, 40401, "user not found")
		return
	}
	if req.Nickname != nil {
		user.Nickname = req.Nickname
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	if req.Email != nil {
		user.Email = req.Email
	}
	if err := h.db.Save(&user).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	_, _ = h.repService.Recalculate(user.ID)
	payload, _ := h.userPayload(user.ID)
	response.OK(c, payload)
}

func (h *Handler) ListRoles(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var roles []model.UserRole
	if err := h.db.Where("user_id = ? AND enabled = ?", userID.(uint64), true).Find(&roles).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, roles)
}

type switchRoleRequest struct {
	RoleCode string `json:"role_code" binding:"required"`
}

func (h *Handler) SwitchRole(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req switchRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid request")
		return
	}
	var role model.UserRole
	if err := h.db.Where("user_id = ? AND role_code = ? AND enabled = ?", userID.(uint64), req.RoleCode, true).First(&role).Error; err != nil {
		response.Fail(c, http.StatusForbidden, 40301, "role not available")
		return
	}
	var user model.User
	if err := h.db.First(&user, userID.(uint64)).Error; err != nil {
		response.Fail(c, http.StatusNotFound, 40401, "user not found")
		return
	}
	user.ActiveRole = req.RoleCode
	if err := h.db.Save(&user).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	token, _ := h.signToken(&user)
	payload, _ := h.userPayload(user.ID)
	response.OK(c, gin.H{"token": token, "user": payload})
}

func (h *Handler) RequestVerifierRole(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		OrgName      string `json:"org_name" binding:"required"`
		Industry     string `json:"industry"`
		ContactEmail string `json:"contact_email"`
		Website      string `json:"website"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid request")
		return
	}
	profile := model.VerifierProfile{UserID: userID.(uint64), OrgName: req.OrgName, Status: "approved"}
	if req.Industry != "" {
		profile.Industry = &req.Industry
	}
	if req.ContactEmail != "" {
		profile.ContactEmail = &req.ContactEmail
	}
	if req.Website != "" {
		profile.Website = &req.Website
	}
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID.(uint64)).Assign(profile).FirstOrCreate(&profile).Error; err != nil {
			return err
		}
		role := model.UserRole{UserID: userID.(uint64), RoleCode: "verifier", Enabled: true}
		return tx.Where("user_id = ? AND role_code = ?", userID.(uint64), "verifier").FirstOrCreate(&role).Error
	}); err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	payload, _ := h.userPayload(userID.(uint64))
	response.OK(c, payload)
}

func (h *Handler) ListCategories(c *gin.Context) {
	var categories []model.Category
	if err := h.db.Where("is_active = ?", true).Find(&categories).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, categories)
}

func (h *Handler) Portfolio(c *gin.Context) {
	address := strings.ToLower(c.Param("address"))
	var user model.User
	if err := h.db.Where("wallet_address = ?", address).First(&user).Error; err != nil {
		response.Fail(c, http.StatusNotFound, 40401, "user not found")
		return
	}
	var works []model.Work
	h.db.Preload("Category").Where("user_id = ? AND visibility = ?", user.ID, "public").Order("created_at desc").Find(&works)
	var evidences []model.EvidenceRecord
	h.db.Where("user_id = ? AND status = ?", user.ID, "confirmed").Order("created_at desc").Find(&evidences)
	var credentials []model.SBTCredential
	h.db.Where("user_id = ? AND status = ?", user.ID, "active").Find(&credentials)
	score, _ := h.repService.GetByUserID(user.ID)
	response.OK(c, gin.H{"user": user, "works": works, "evidences": evidences, "credentials": credentials, "reputation": score})
}

func (h *Handler) AdminUsers(c *gin.Context) {
	var users []model.User
	if err := h.db.Preload("Roles").Order("created_at desc").Limit(100).Find(&users).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, users)
}

func (h *Handler) AdminChains(c *gin.Context) {
	var chains []model.ChainNetwork
	h.db.Order("id asc").Find(&chains)
	var contracts []model.ContractConfig
	h.db.Order("id asc").Find(&contracts)
	response.OK(c, gin.H{"chains": chains, "contracts": contracts})
}

func (h *Handler) AdminStatistics(c *gin.Context) {
	var userCount, workCount, evidenceCount, credentialCount int64
	h.db.Model(&model.User{}).Count(&userCount)
	h.db.Model(&model.Work{}).Count(&workCount)
	h.db.Model(&model.EvidenceRecord{}).Where("status = ?", "confirmed").Count(&evidenceCount)
	h.db.Model(&model.SBTCredential{}).Where("status = ?", "active").Count(&credentialCount)
	response.OK(c, gin.H{"users": userCount, "works": workCount, "confirmed_evidences": evidenceCount, "credentials": credentialCount})
}

type updateUserRolesRequest struct {
	Roles      []string `json:"roles" binding:"required"`
	ActiveRole string   `json:"active_role"`
}

func (h *Handler) AdminUpdateUserRoles(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid user id")
		return
	}
	var req updateUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Roles) == 0 {
		response.Fail(c, http.StatusBadRequest, 40001, "roles required")
		return
	}
	allowed := map[string]bool{"creator": true, "verifier": true, "auditor": true, "admin": true}
	activeRole := req.ActiveRole
	if activeRole == "" {
		activeRole = req.Roles[0]
	}
	if !allowed[activeRole] {
		response.Fail(c, http.StatusBadRequest, 40001, "invalid active role")
		return
	}
	for _, role := range req.Roles {
		if !allowed[role] {
			response.Fail(c, http.StatusBadRequest, 40001, "invalid role: "+role)
			return
		}
	}

	adminID, _ := c.Get("user_id")
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.UserRole{}).Where("user_id = ?", targetID).Update("enabled", false).Error; err != nil {
			return err
		}
		for _, role := range req.Roles {
			record := model.UserRole{UserID: targetID, RoleCode: role, Enabled: true}
			if err := tx.Where("user_id = ? AND role_code = ?", targetID, role).Assign(record).FirstOrCreate(&record).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(&model.User{}).Where("id = ?", targetID).Update("active_role", activeRole).Error; err != nil {
			return err
		}
		raw, _ := json.Marshal(req)
		userID := adminID.(uint64)
		return tx.Create(&model.AuditLog{UserID: &userID, Action: "admin.update_roles", Resource: "users", Detail: string(raw)}).Error
	}); err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, gin.H{"updated": true})
}

func (h *Handler) AdminAuditLogs(c *gin.Context) {
	var logs []model.AuditLog
	if err := h.db.Order("created_at desc").Limit(100).Find(&logs).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, err.Error())
		return
	}
	response.OK(c, logs)
}

func (h *Handler) Logout(c *gin.Context) {
	response.OK(c, gin.H{"message": "logged out"})
}
