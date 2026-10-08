package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/uniflou/uniflou-lms/backend/internal/admin"
	"github.com/uniflou/uniflou-lms/backend/internal/ai"
	"github.com/uniflou/uniflou-lms/backend/internal/auth"
	"github.com/uniflou/uniflou-lms/backend/internal/config"
	"github.com/uniflou/uniflou-lms/backend/internal/parent"
	"github.com/uniflou/uniflou-lms/backend/internal/plan"
	"github.com/uniflou/uniflou-lms/backend/internal/progress"
	"github.com/uniflou/uniflou-lms/backend/internal/server/middleware"
	"github.com/uniflou/uniflou-lms/backend/internal/task"
	"github.com/uniflou/uniflou-lms/backend/internal/tutor"
	"github.com/uniflou/uniflou-lms/backend/internal/user"
	"gorm.io/gorm"
)

type Server struct {
	engine *gin.Engine
	port   string
}

func New(cfg *config.Config, db *gorm.DB, log zerolog.Logger) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.Recovery(log))
	r.Use(middleware.Logger(log))
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS(cfg.CORSOrigins))
	r.Use(middleware.MaxBody(1 << 20))

	jwt := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	if cfg.SMTPHost == "" {
		log.Warn().Msg("SMTP is not configured; password reset emails are disabled")
	}
	mailer := auth.NewSMTPMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.AppURL)
	authSvc := auth.NewService(db, jwt, log, mailer)
	if cfg.AdminBootstrapEnabled {
		if err := authSvc.EnsureBootstrapAdmin(context.Background(), cfg.AdminEmail, cfg.AdminPassword); err != nil {
			log.Warn().Err(err).Msg("admin bootstrap failed")
		}
	}
	aiClient := ai.NewClient(cfg)

	h := Handlers{
		Auth:     auth.NewHandler(authSvc),
		User:     user.NewHandler(user.NewService(db)),
		Plan:     plan.NewHandler(plan.NewService(db)),
		Task:     task.NewHandler(task.NewService(db)),
		Progress: progress.NewHandler(progress.NewService(db)),
		AI:       ai.NewHandler(ai.NewService(db, aiClient, log, cfg.RateLimitAI)),
		Parent:   parent.NewHandler(parent.NewService(db)),
		Tutor:    tutor.NewHandler(tutor.NewService(db)),
		Admin:    admin.NewHandler(admin.NewService(db)),
	}

	registerRoutes(r, h, jwt, cfg.RateLimitLogin, 3)

	return &Server{engine: r, port: cfg.Port}
}

func (s *Server) Run() error {
	return s.engine.Run(fmt.Sprintf(":%s", s.port))
}
