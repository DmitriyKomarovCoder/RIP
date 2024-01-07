package handler

import (
	_ "RIP/docs"
	"RIP/internal/app/config"
	"RIP/internal/app/pkg/hash"
	"RIP/internal/app/redis"
	"RIP/internal/app/repository"
	"RIP/internal/app/role"
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"

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
	Config     *config.Config
	Redis      *redis.Client
	//TokenManager auth.TokenManager
	Hasher hash.PasswordHasher
}

func NewHandler(
	l *logrus.Logger,
	r *repository.Repository,
	m *minio.Client,
	conf *config.Config,
	red *redis.Client,
	// tokenManager auth.TokenManager,
) *Handler {
	return &Handler{
		Logger:     l,
		Repository: r,
		Minio:      m,
		Config:     conf,
		Redis:      red,
		//TokenManager: tokenManager,
		Hasher: hash.NewSHA256Hasher(os.Getenv("SALT")),
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	// услуги
	api.GET("/companies", h.CompaniesList)      // ?
	api.GET("/companies/:id", h.GetCompanyById) // ?
	api.POST("/companies", h.AddCompany)
	api.PUT("/companies/:id", h.WithAuthCheck(role.Moderator, role.Admin), h.UpdateCompany)
	api.DELETE("/companies/:id", h.WithAuthCheck(role.Moderator, role.Admin), h.DeleteCompany)
	api.POST("/companies/request", h.WithAuthCheck(role.Buyer, role.Moderator, role.Admin), h.AddCompanyToRequest)
	api.Use(cors.Default()).DELETE("/companies/delete/:id", h.DeleteCompany)

	// заявки
	api.GET("/tenders", h.WithAuthCheck(role.Buyer, role.Moderator, role.Admin), h.TenderList)
	api.GET("/tenders/:id", h.WithAuthCheck(role.Buyer, role.Moderator, role.Admin), h.GetTenderById)
	// api.POST("/tenders/", h.CreateDraft)
	api.PUT("/tenders", h.WithAuthCheck(role.Buyer, role.Moderator, role.Admin), h.UpdateTender)

	// статусы
	api.PUT("/tenders/form", h.WithAuthCheck(role.Buyer, role.Moderator, role.Admin), h.FormTenderRequest)
	api.PUT("/tenders/updateStatus", h.WithAuthCheck(role.Moderator, role.Admin), h.UpdateStatusTenderRequest)
	//api.PUT("/tenders/finish/:id", h.WithAuthCheck([]ds.Role{ds.Admin}), h.FinishTenderRequest)

	api.DELETE("/tenders", h.WithAuthCheck(role.Buyer, role.Moderator, role.Admin), h.DeleteTender)

	// m-m
	api.DELETE("/tender-request-company", h.WithoutJWTError(role.Buyer, role.Moderator, role.Admin), h.DeleteCompanyFromRequest)
	api.PUT("/tender-request-company", h.WithoutJWTError(role.Buyer, role.Moderator, role.Admin), h.UpdateTenderCompany)
	registerStatic(router)

	// auth && reg
	api.POST("/user/signIn", h.Login)
	api.POST("/user/signUp", h.Register)
	api.POST("/user/logout", h.Logout)

	// асинхронный сервис
	api.PUT("/tenders/user-form-start", h.WithoutJWTError(role.Buyer), h.UserRequest) // обращение к асинхронному сервису
	api.PUT("/tenders/user-form-finish", h.FinishUserRequest)                         // обращение к асинхронному сервису

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
