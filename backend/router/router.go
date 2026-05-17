package router

import (
	"io/fs"
	"net/http"
	"strings"
	"time"

	"posixfy-cloud/backend/config"
	"posixfy-cloud/backend/handler"
	"posixfy-cloud/backend/middleware"
	"posixfy-cloud/backend/service"
	"posixfy-cloud/backend/static"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config, userService *service.UserService, fsClient *service.FSClient) *gin.Engine {
	r := gin.Default()

	// CORS configuration from environment
	allowAllOrigins := len(cfg.CORSOrigins) == 1 && cfg.CORSOrigins[0] == "*"
	corsConfig := cors.Config{
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Authorization", "Content-Type", "X-Expected-MTime", "X-Expected-Size"},
		ExposeHeaders: []string{"Content-Disposition"},
	}
	if allowAllOrigins {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = cfg.CORSOrigins
		corsConfig.AllowCredentials = true
	}
	r.Use(cors.New(corsConfig))

	authHandler := handler.NewAuthHandler(userService, cfg)
	adminHandler := handler.NewAdminUserHandler(userService)
	fsHandler := handler.NewFSProxyHandler(fsClient)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			// Rate limit: 10 login attempts per minute per IP
			auth.POST("/login", middleware.RateLimit(10, time.Minute), authHandler.Login)
		}

		authed := api.Group("")
		authed.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			authed.GET("/auth/me", authHandler.Me)

			fsGroup := authed.Group("/fs")
			{
				fsGroup.GET("/mounts", fsHandler.Mounts)
				fsGroup.GET("/list", fsHandler.List)
				fsGroup.GET("/file", fsHandler.Download)
				fsGroup.POST("/upload", fsHandler.Upload)
				fsGroup.DELETE("/delete", fsHandler.Delete)
				fsGroup.POST("/mkdir", fsHandler.Mkdir)
				fsGroup.GET("/watch", fsHandler.Watch)
			}

			admin := authed.Group("/admin")
			admin.Use(middleware.RequireAdmin())
			{
				admin.GET("/users", adminHandler.List)
				admin.POST("/users", adminHandler.Create)
				admin.GET("/users/:id", adminHandler.Get)
				admin.PUT("/users/:id", adminHandler.Update)
				admin.DELETE("/users/:id", adminHandler.Delete)
			}
		}
	}

	// Serve embedded frontend for non-API routes
	frontendFS := static.DistFS()
	if frontendFS != nil {
		r.NoRoute(serveFrontend(frontendFS))
	}

	return r
}

// serveFrontend serves the embedded Vue SPA.
// For known static file extensions, serve the file directly.
// For everything else, serve index.html (SPA client-side routing).
func serveFrontend(frontendFS fs.FS) gin.HandlerFunc {
	fileServer := http.FileServer(http.FS(frontendFS))

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip API routes (should not reach here, but just in case)
		if strings.HasPrefix(path, "/api/") {
			c.Next()
			return
		}

		// Try to serve the exact file
		trimmed := strings.TrimPrefix(path, "/")
		if trimmed == "" {
			trimmed = "index.html"
		}

		if f, err := frontendFS.Open(trimmed); err == nil {
			f.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// Fallback to index.html for SPA routing
		c.Request.URL.Path = "/"
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}
