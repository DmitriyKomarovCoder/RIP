package main

import (
	"RIP/internal/app/config"
	"RIP/internal/app/dsn"
	"RIP/internal/app/handler"
	app "RIP/internal/app/pkg"
	"RIP/internal/app/pkg/auth"
	"RIP/internal/app/redis"
	"RIP/internal/app/repository"
	Minio "RIP/internal/app/s3/minio"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

// @title Tender App
// @version 1.0
// @description App for serving tender requests

// @host 127.0.0.1:8080
// @schemes http
// @BasePath /
func main() {
	logger := logrus.New()
	minioClient := Minio.NewMinioClient(logger)

	router := gin.Default()
	//router.Use(corsMiddleware())

	conf, err := config.NewConfig(logger)
	if err != nil {
		logger.Fatalf("Error with configuration reading: %s", err)
	}

	ctx := context.Background()
	redisClient, errRedis := redis.New(ctx, conf.Redis)
	if errRedis != nil {
		logger.Fatalf("Errof with redis connect: %s", err)
	}

	postgresString, errPost := dsn.FromEnv()
	if errPost != nil {
		logger.Fatalf("Error of reading postgres line: %s", errPost)
	}

	rep, errRep := repository.NewRepository(postgresString, logger)
	if errRep != nil {
		logger.Fatalf("Error from repository: %s", err)
	}

	tokenManager, err := auth.NewManager(os.Getenv("TOKEN_SECRET"))
	if err != nil {
		logger.Fatalln(err)
	}

	hand := handler.NewHandler(logger, rep, minioClient, conf, redisClient, tokenManager)
	application := app.NewApp(conf, router, logger, hand)
	application.RunApp()
}

//func corsMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
//		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
//		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
//		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
//
//		if c.Request.Method == "OPTIONS" {
//			c.AbortWithStatus(204)
//			return
//		}
//
//		c.Next()
//	}
//}
