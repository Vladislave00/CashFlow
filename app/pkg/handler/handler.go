package handler

import (
	"github.com/Vladislave00/CashFlow/app/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Static("/static", "static")

	router.LoadHTMLGlob("static/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})

	router.GET("/accounting", func(c *gin.Context) {
		c.HTML(200, "accounting.html", nil)
	})

	router.GET("/transactions", func(c *gin.Context) {
		c.HTML(200, "transactions.html", nil)
	})

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		values := api.Group("/values")
		{
			values.GET("/:id", h.getValueById)
			values.GET("/getByName/:name", h.getValueByName)
		}
		accountings := api.Group("/accountings")
		{
			accountings.POST("/", h.createAccounting)
			accountings.GET("/", h.getAccountings)
			accountings.GET("/:id", h.getAccountingById)
			accountings.PUT("/:id", h.updateAccounting)
			accountings.DELETE("/:id", h.deleteAccounting)

			accounts := accountings.Group(":id/accounts")
			{
				accounts.POST("/", h.createAccount)
				accounts.GET("/", h.getAccounts)
				accounts.GET("/general", h.getGeneralAccount)

				transactions := accounts.Group("/transactions")
				{
					transactions.POST("/", h.createTransaction)
					transactions.GET("/", h.getTransactions)
					transactions.GET("/:accountId", h.getTransactionsByAccountId)
				}
			}
		}
		accounts := api.Group("/accounts")
		{
			accounts.GET("/:id", h.getAccountById)
			accounts.PUT("/:id", h.updateAccount)
			accounts.DELETE("/:id", h.deleteAccount)
		}
		transactions := api.Group("/transactions")
		{
			transactions.GET("/:id", h.getTransactionById)
			transactions.PUT("/:id", h.updateTransaction)
			transactions.DELETE("/:id", h.deleteTransaction)
		}
	}
	return router
}
