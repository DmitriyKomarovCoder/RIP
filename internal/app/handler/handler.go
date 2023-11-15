package handler

import (
	"RIP/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
)

const (
	creatorID   = 1
	moderatorID = 1
)

type Handler struct {
	Logger     *logrus.Logger
	Repository *repository.Repository
	Minio      *minio.Client
}

func NewHandler(l *logrus.Logger, r *repository.Repository, m *minio.Client) *Handler {
	return &Handler{
		Logger:     l,
		Repository: r,
		Minio:      m,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api")
	// услуги
	api.GET("/companies", h.CompaniesList)

	api.GET("/companies/:id", h.GetCompanyById)
	api.POST("/companies/", h.AddCompany)
	api.PUT("/companies/", h.UpdateCompany)
	api.DELETE("/companies/:id", h.DeleteCompany)
	api.POST("/companies/request", h.AddCompanyToRequest)

	// заявки
	api.GET("/tenders", h.TenderList)
	api.GET("/tenders/:id", h.GetTenderById)
	api.POST("/tenders/", h.CreateDraft)
	api.PUT("/tenders/", h.UpdateTender)
	api.GET("/tenders/form/:id", h.FormTenderRequest)
	api.GET("tender/reject/:id", h.RejectTenderRequest)
	api.GET("tender/finish/:id", h.FinishTenderRequest)
	api.DELETE("/tenders/:id", h.DeleteTender)

	//m-m
	api.DELETE("/transacion-request-company/company/:id", h.DeleteCompanyFromRequest)

	registerStatic(router)
}

func registerStatic(router *gin.Engine) {
	router.LoadHTMLGlob("static/templates/*")
	router.Static("/static", "./static")
	router.Static("/css", "./static")
	router.Static("/img", "./static")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	h.Logger.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      errorStatusCode,
		"description": err.Error(),
	})
}

func (h *Handler) successHandler(ctx *gin.Context, key string, status int, data interface{}) {
	ctx.JSON(status, gin.H{
		"status": status,
		key:      data,
	})
}
