package handlers

//import (
//	"auth-server/repository"
//	"auth-server/services"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/go-redis/redis/v8"
//	"gorm.io/gorm"
//)
//
//func RegisterAuthRoutes(r *gin.Engine, db *gorm.DB, redisClient *redis.Client, smtpConfig services.SmtpConfig) {
//	repo := repository.NewUserRepository(db, redisClient)
//
//	svc, err := services.NewAuthService(repo, smtpConfig)
//	if err != nil {
//		panic(fmt.Errorf("failed initializing AuthService: %s", err))
//	}
//
//	r.POST("/auth/register", svc.RegisterHandler)
//	r.POST("/auth/login", svc.LoginHandler)
//	r.POST("/auth/token", svc.GetTokenHandler)
//	r.POST("/auth/refresh", svc.RefreshTokenHandler)
//
//	r.POST("/auth/reset-password", svc.ResetPasswordHandler)
//	r.POST("/auth/reset-password/change-password", svc.ResetChangePasswordHandler)
//
//	codeRequestGroup := r.Group("/", svc.JWTAuthMiddleware())
//	codeRequestGroup.POST("/auth/register/code/request", svc.CodeRequestHandler)
//	codeRequestGroup.POST("/auth/register/code/confirm", svc.CodeConfirmationHandler)
//}
