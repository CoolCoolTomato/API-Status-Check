package api

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *Handler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.POST("/api/auth/login", LoginHandler)

	admin := r.Group("/api/admin")
	admin.Use(JWTMiddleware())
	{
		admin.POST("/apis", handler.CreateAPI)
		admin.GET("/apis", handler.GetAPIs)
		admin.GET("/apis/:id", handler.GetAPI)
		admin.PUT("/apis/:id", handler.UpdateAPI)
		admin.PATCH("/apis/:id", handler.UpdateAPI)
		admin.DELETE("/apis/:id", handler.DeleteAPI)
	}

	checks := r.Group("/api/checks")
	{
		checks.GET("/history", handler.GetHistory)
		checks.GET("/recent", handler.GetRecent)
		checks.POST("/run", handler.RunCheck)
	}

	// Serve frontend static files if web/dist exists
	if _, err := os.Stat("web/dist"); err == nil {
		r.Static("/assets", "web/dist/assets")
		r.StaticFile("/favicon.ico", "web/dist/favicon.ico")
		r.NoRoute(func(c *gin.Context) {
			c.File("web/dist/index.html")
		})
	}

	return r
}
