package handlers

//
//import (
//	"devices-api/repository"
//	"devices-api/services"
//	"github.com/gin-gonic/gin"
//	"gorm.io/gorm"
//)
//
//func RegisterYandexRoutes(r *gin.Engine, db *gorm.DB) {
//	devicesRepo := repository.NewDevicesRepo(db)
//
//	svc, err := services.NewDevicesService(devicesRepo)
//	if err != nil {
//		panic(err)
//	}
//
//	r.HEAD("/yandex/v1.0", svc.AliveStatus)
//	r.POST("/ping", svc.Ping)
//
//	protected := r.Group("/")
//	protected.Use(svc.EmailConfirmedAuthMiddleware())
//	{
//		protected.POST("/yandex/v1.0/user/unlink", svc.UnlinkAccount)
//		protected.GET("/yandex/v1.0/user/devices", svc.ListDevices)
//		protected.POST("/yandex/v1.0/user/devices/query", svc.QueryDevices)
//		protected.POST("/yandex/v1.0/user/devices/action", svc.ActionDevices)
//	}
//}
