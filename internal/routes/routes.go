package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wallet/internal/handler"
	"github.com/wallet/internal/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(r *gin.Engine, h *handler.Handler, jwtSecret string) {
	r.Use(middleware.CORS())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", h.Health)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.RefreshToken)
		}

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{
			protected.GET("/auth/me", h.GetProfile)

			protected.GET("/wallet/balance", h.GetBalance)
			protected.POST("/wallet/credit", h.CreditWallet)
			protected.POST("/wallet/debit", h.DebitWallet)

			transactions := protected.Group("/transactions")
			{
				transactions.POST("/transfer", h.Transfer)
				transactions.GET("", h.GetTransactions)
			}

			protected.POST("/pix", h.CreatePixPayment)
			protected.POST("/boleto", h.GenerateBoleto)
		}
	}
}
