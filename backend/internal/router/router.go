package router

import (
	"log"

	"web3proof/backend/internal/config"
	"web3proof/backend/internal/handler"
	"web3proof/backend/internal/middleware"
	"web3proof/backend/internal/pkg/eth"
	"web3proof/backend/internal/pkg/storage"
	"web3proof/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB, redisClient *redis.Client) *gin.Engine {
	gin.SetMode(cfg.GinMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.CORSOrigin))

	ethClient, err := eth.NewEthClient(cfg)
	if err != nil {
		log.Printf("warning: eth client: %v", err)
	}
	store := storage.NewLocalStore(cfg.StorageRoot, cfg.PublicBaseURL)

	authService := service.NewAuthService(db, redisClient, cfg.ChainID, cfg.AdminWalletAddress, cfg.AuthDomain, cfg.AuthURI)
	repService := service.NewReputationService(db)
	workService := service.NewWorkService(db, store, cfg, ethClient, repService)
	certService := service.NewCertificationService(db, store, cfg, ethClient, repService)
	verifyService := service.NewVerificationService(db)
	chainEventService := service.NewChainEventService(db, ethClient)

	h := handler.New(db, cfg.JWTSecret, authService, repService, ethClient)
	wh := handler.NewWorkHandler(workService)
	ch := handler.NewCertificationHandler(certService)
	vh := handler.NewVerificationHandler(verifyService)
	ceh := handler.NewChainEventHandler(chainEventService)

	r.GET("/health", h.Health)
	r.Static("/storage", cfg.StorageRoot)

	api := r.Group("/api/v1")
	{
		api.GET("/auth/nonce", h.GetNonce)
		api.POST("/auth/login", h.Login)
		api.GET("/categories", h.ListCategories)
		api.GET("/portfolio/:address", h.Portfolio)
		verify := api.Group("/verify")
		verify.Use(middleware.OptionalJWTAuth(cfg.JWTSecret, db))
		{
			verify.POST("/file", vh.VerifyFile)
			verify.GET("/evidence/:no", vh.VerifyEvidence)
			verify.GET("/certificate/:no", vh.VerifyCertificate)
			verify.GET("/wallet/:address", vh.VerifyWallet)
		}

		auth := api.Group("")
		auth.Use(middleware.JWTAuth(cfg.JWTSecret, db))
		{
			auth.GET("/auth/me", h.Me)
			auth.POST("/auth/logout", h.Logout)
			auth.GET("/dashboard/summary", h.DashboardSummary)
			auth.PUT("/users/profile", h.UpdateProfile)
			auth.GET("/users/roles", h.ListRoles)
			auth.PUT("/users/current-role", h.SwitchRole)
			auth.POST("/users/verifier-profile", h.RequestVerifierRole)

			auth.POST("/works", middleware.RequireRole("creator", "admin"), wh.Create)
			auth.GET("/works", middleware.RequireRole("creator", "admin"), wh.List)
			auth.GET("/works/:id", middleware.RequireRole("creator", "admin"), wh.Get)
			auth.POST("/works/:id/files", middleware.RequireRole("creator", "admin"), wh.UploadFile)
			auth.POST("/works/:id/evidence/prepare", middleware.RequireRole("creator", "admin"), wh.PrepareEvidence)
			auth.POST("/works/:id/evidence/confirm", middleware.RequireRole("creator", "admin"), wh.ConfirmEvidence)
			auth.POST("/certificates/generate", middleware.RequireRole("creator", "admin"), wh.GenerateCertificate)
			auth.GET("/evidence/my", middleware.RequireRole("creator", "admin"), wh.ListEvidence)
			auth.GET("/certificates/my", middleware.RequireRole("creator", "admin"), wh.ListCertificates)

			auth.POST("/applications", middleware.RequireRole("creator", "admin"), ch.Apply)
			auth.GET("/applications/my", middleware.RequireRole("creator", "admin"), ch.ListMy)

			auth.GET("/verifier/reports", middleware.RequireRole("verifier", "admin"), vh.ListReports)

			auditor := auth.Group("/auditor")
			auditor.Use(middleware.RequireRole("auditor", "admin"))
			{
				auditor.GET("/applications", ch.ListPending)
				auditor.PUT("/applications/:id/review", ch.Review)
				auditor.POST("/applications/:id/sbt/prepare", ch.PrepareMint)
				auditor.POST("/applications/:id/sbt/mint", ch.ConfirmMint)
			}

			admin := auth.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/users", h.AdminUsers)
				admin.PUT("/users/:id/roles", h.AdminUpdateUserRoles)
				admin.GET("/chains", h.AdminChains)
				admin.GET("/statistics", h.AdminStatistics)
				admin.GET("/audit-logs", h.AdminAuditLogs)
				admin.POST("/chain-events/sync", ceh.Sync)
				admin.GET("/chain-events", ceh.List)
			}
		}
	}

	return r
}
