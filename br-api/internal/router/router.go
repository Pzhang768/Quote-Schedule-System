package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/melfish/br-api/internal/handler"
	"github.com/melfish/br-api/internal/logger"
	"github.com/melfish/br-api/internal/service"
	"github.com/melfish/br-api/internal/store"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func New(db *gorm.DB) *gin.Engine {
	quoteStore := store.NewQuoteStore(db)
	technicianStore := store.NewTechnicianStore(db)
	jobStore := store.NewJobStore(db)
	notificationStore := store.NewNotificationStore(db)
	managerStore := store.NewManagerStore(db)

	quoteSvc := service.NewQuoteService(quoteStore)
	technicianSvc := service.NewTechnicianService(technicianStore, jobStore)
	jobSvc := service.NewJobService(db, jobStore, quoteStore, notificationStore)
	notificationSvc := service.NewNotificationService(notificationStore)
	managerSvc := service.NewManagerService(managerStore)

	quotes := handler.NewQuoteHandler(quoteSvc)
	technicians := handler.NewTechnicianHandler(technicianSvc)
	jobs := handler.NewJobHandler(jobSvc)
	notifications := handler.NewNotificationHandler(notificationSvc)
	managers := handler.NewManagerHandler(managerSvc)

	r := gin.New()
	r.Use(requestLogger())
	r.Use(corsMiddleware())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1.GET("/quotes", quotes.List)
	v1.POST("/quotes", quotes.Create)

	v1.GET("/managers", managers.List)
	v1.GET("/technicians", technicians.List)
	v1.GET("/technicians/:id/jobs", technicians.GetSchedule)

	v1.GET("/jobs/:id", jobs.Get)
	v1.POST("/jobs", jobs.Assign)
	v1.PATCH("/jobs/:id/complete", jobs.Complete)

	v1.GET("/notifications", notifications.List)
	v1.GET("/notifications/stream", notifications.Stream)
	v1.PATCH("/notifications/:id/read", notifications.Read)

	return r
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Log.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start).String(),
		)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods: []string{"GET", "POST", "PATCH"},
		AllowHeaders: []string{"Content-Type"},
	})
}
