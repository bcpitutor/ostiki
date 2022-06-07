package apiserver

import (
	"html/template"
	"time"

	"github.com/bcpitutor/ostiki/appconfig"
	"github.com/bcpitutor/ostiki/logger"
	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/models"
	"github.com/bcpitutor/ostiki/repositories"
	"github.com/bcpitutor/ostiki/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

type ServerParams struct {
	dig.In

	AppConfig            *appconfig.AppConfig
	TikiLogger           *logger.TikiLogger
	AWSService           *services.AWSService
	PermissionRepository *repositories.PermissionRepository
	SessionRepository    *repositories.SessionRepository
	DomainRepository     *repositories.DomainRepository
	BanRepository        *repositories.BanRepository
	GroupRepository      *repositories.GroupRepository
	TicketRepository     *repositories.TicketRepository
	DBLayer              models.DBLayer
	CacheLayer           models.CacheLayer `optional:"true"`
}
type Server struct {
	appConfig            *appconfig.AppConfig
	tikiLogger           *logger.TikiLogger
	awsServices          *services.AWSService
	permissionRepository repositories.PermissionRepository
	sessionRepository    repositories.SessionRepository
	banRepository        repositories.BanRepository
	domainRepository     repositories.DomainRepository
	groupRepository      repositories.GroupRepository
	ticketRepository     repositories.TicketRepository
	dbLayer              models.DBLayer
	cacheLayer           models.CacheLayer
}

func ProvideServer(params ServerParams) *Server {
	return &Server{
		appConfig:            params.AppConfig,
		tikiLogger:           params.TikiLogger,
		awsServices:          params.AWSService,
		permissionRepository: *params.PermissionRepository,
		sessionRepository:    *params.SessionRepository,
		banRepository:        *params.BanRepository,
		domainRepository:     *params.DomainRepository,
		groupRepository:      *params.GroupRepository,
		ticketRepository:     *params.TicketRepository,
		dbLayer:              params.DBLayer,
		cacheLayer:           params.CacheLayer,
	}
}

func (s *Server) Run() {
	config := s.appConfig
	logger := s.tikiLogger.Logger
	sugar := logger.Sugar()

	switch config.Deployment {
	case "dev", "development", "debug":
		gin.SetMode("debug")
	case "test", "preprod":
		gin.SetMode("test")
	default:
		gin.SetMode("release")
	}

	ginEngine := gin.Default()
	html := template.Must(template.ParseFiles("html/pages/index.html"))
	ginEngine.SetHTMLTemplate(html)

	ginEngine.Static("/static", "./static")
	ginEngine.LoadHTMLGlob("html/pages/*")

	ginEngine.Use(
		cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
	)

	addMiscHandlers(ginEngine, middleware.GinHandlerVars{
		Logger:            logger,
		AppConfig:         config,
		SessionRepository: &s.sessionRepository,
		BanRepository:     &s.banRepository,
	})

	addDomainHandlers(ginEngine,
		middleware.GinHandlerVars{
			Logger:               logger,
			AppConfig:            config,
			DomainRepository:     &s.domainRepository,
			GroupRepository:      &s.groupRepository,
			TicketRepository:     &s.ticketRepository,
			PermissionRepository: &s.permissionRepository,
			AWSService:           s.awsServices,
		},
	)

	addTicketHandlers(ginEngine,
		middleware.GinHandlerVars{
			Logger:               logger,
			AppConfig:            config,
			GroupRepository:      &s.groupRepository,
			TicketRepository:     &s.ticketRepository,
			PermissionRepository: &s.permissionRepository,
			AWSService:           s.awsServices,
		},
	)

	addGroupHandlers(ginEngine,
		middleware.GinHandlerVars{
			Logger:               logger,
			AppConfig:            config,
			GroupRepository:      &s.groupRepository,
			TicketRepository:     &s.ticketRepository,
			PermissionRepository: &s.permissionRepository,
			AWSService:           s.awsServices,
		},
	)

	addBanHandlers(ginEngine,
		middleware.GinHandlerVars{
			Logger:        logger,
			AppConfig:     config,
			BanRepository: &s.banRepository,
		},
	)

	sugar.Infof(
		"Server started to listen address on %s:%s",
		config.ListenerHost,
		config.ListenerPort,
	)
	ginEngine.Run(config.ListenerHost + ":" + config.ListenerPort)
}
