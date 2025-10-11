package services

//
//import (
//	"errors"
//	"github.com/gin-gonic/gin"
//	"github.com/go-redis/redis/v8"
//	"github.com/google/uuid"
//	"golang.org/x/crypto/bcrypt"
//	"gorm.io/gorm"
//	"log"
//	"net/http"
//)
//
//func (s *AuthService) ResetPasswordHandler(c *gin.Context) {
//	type ResetPasswordRequest struct {
//		Email string `json:"email" binding:"required"`
//	}
//
//	var req ResetPasswordRequest
//	if err := c.ShouldBind(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error": "Invalid request",
//		})
//		return
//	}
//
//	user, err := s.repo.GetUserByEmail(req.Email)
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			c.JSON(http.StatusForbidden, gin.H{
//				"error": "User with this email not found",
//			})
//			return
//		}
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": "cannot reset password",
//		})
//		log.Printf("failed to get user from storage: %v\n", err)
//		return
//	}
//
//	// generate uuidV7Code code for resetting password
//	uuidV7Code, err := uuid.NewV7()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": "cannot reset password",
//		})
//		log.Printf("failed to generate uuidV7Code: %v\n", err)
//		return
//	}
//
//	uuidCode := uuidV7Code.String()
//
//	err = s.repo.SetResetPasswordCode(uuidCode, user.Email)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": "cannot reset password",
//		})
//		log.Printf("failed to save code in redis: %v\n", err)
//		return
//	}
//
//	err = s.sendResetPasswordCode(uuidCode, user.Email)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": "cannot reset password",
//		})
//		log.Printf("failed to send link on email: %v\n", err)
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"email": user.Email,
//	})
//}
//
//func (s *AuthService) ResetChangePasswordHandler(c *gin.Context) {
//	type ConfirmResetPasswordRequest struct {
//		Code        string `json:"code"`
//		NewPassword string `json:"password" binding:"required"`
//	}
//	var req ConfirmResetPasswordRequest
//	if err := c.ShouldBind(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error": "Invalid request",
//		})
//		return
//	}
//
//	// check if code is valid (present in redis)
//	emailInDb, err := s.repo.GetEmailByResetPasswordCode(req.Code)
//	if err != nil {
//		if errors.Is(err, redis.Nil) {
//			c.JSON(http.StatusForbidden, gin.H{
//				"error": "Reset password code not found",
//			})
//			return
//		}
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": "cannot reset password",
//		})
//		log.Printf("failed to get email by reset password code: %v\n", err)
//		return
//	}
//
//	user, err := s.repo.GetUserByEmail(emailInDb)
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			c.JSON(http.StatusForbidden, gin.H{
//				"error": "User with this email not found",
//			})
//			return
//		}
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": "cannot reset password",
//		})
//		log.Printf("failed to get user from storage: %v\n", err)
//		return
//	}
//
//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.Dex.faultCost)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": "cannot reset password",
//		})
//		log.Printf("failed to hash password: %v\n", err)
//		return
//	}
//
//	user.Password = string(hashedPassword)
//	err = s.repo.UpdateUserPassword(user.ID, user.Password)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": "cannot reset password",
//		})
//		log.Printf("failed to update user in storage: %v\n", err)
//		return
//	}
//
//	err = s.repo.DeleteResetPasswordCode(req.Code)
//	if err != nil {
//		log.Printf("failed to delete reset password code: %v\n", err)
//	}
//	c.JSON(http.StatusOK, gin.H{})
//}
