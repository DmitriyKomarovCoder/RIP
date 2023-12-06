package handler

import (
	"RIP/internal/app/repository"
	"net/http"

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
	api.POST("/companies", h.AddCompany)
	api.PUT("/companies/:id", h.UpdateCompany)
	api.DELETE("/companies/:id", h.DeleteCompany)
	api.POST("/companies/request/:id", h.AddCompanyToRequest)

	// заявки
	api.GET("/tenders", h.TenderList)
	api.GET("/tenders/:id", h.GetTenderById)
	//api.POST("/tenders/", h.CreateDraft)
	api.PUT("/tenders/", h.UpdateTender)
	api.PUT("/tenders/form/:id", h.FormTenderRequest)
	api.PUT("tenders/reject/:id", h.RejectTenderRequest)
	api.PUT("tenders/finish/:id", h.FinishTenderRequest)
	api.DELETE("/tenders/:id", h.DeleteTender)

	//m-m
	api.DELETE("/tender-request-company", h.DeleteCompanyFromRequest)
	api.PUT("/tender-request-company/", h.UpdateTenderCompany)
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

func (h *Handler) successHandler(ctx *gin.Context, key string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		key:      data,
	})
}

func (h *Handler) successAddHandler(ctx *gin.Context, key string, data interface{}) {
	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		key:      data,
	})
}
