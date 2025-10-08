package services

import (
	"crypto/rsa"
	"devices-api/models"
	"devices-api/repository"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DevicesService struct {
	deviceRepo *repository.DevicesRepo

	accessPublicKey *rsa.PublicKey
}

func NewDevicesService(deviceRepo *repository.DevicesRepo) (*DevicesService, error) {
	accessPublicKey, err := loadPublicKey("keys/access_public.pem")
	if err != nil {
		return nil, err
	}

	return &DevicesService{deviceRepo: deviceRepo, accessPublicKey: accessPublicKey}, nil
}

func (s *DevicesService) AliveStatus(c *gin.Context) {
	c.Status(200)
}

func (s *DevicesService) Ping(c *gin.Context) {
	c.JSON(200, "pong")
}

func (s *DevicesService) UnlinkAccount(c *gin.Context) {
	// idk what to do here, so
	c.Status(200)
}

func (s *DevicesService) ListDevices(c *gin.Context) {
	xRequestId := c.GetHeader("X-Request-ID")

	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}

	devices, err := s.deviceRepo.GetDevicesByUserId(fmt.Sprint(userId))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot list devices",
		})
		log.Printf("cannot list devices for user %s: %v", userId, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"request_id": xRequestId,
		"payload": gin.H{
			"user_id": userId,
			"devices": devices,
		},
	})
}

func (s *DevicesService) QueryDevices(c *gin.Context) {
	xRequestId := c.GetHeader("X-Request-ID")

	var devicesQuery models.DevicesQuery
	if err := c.ShouldBindJSON(&devicesQuery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	var deviceIds []string
	for _, device := range devicesQuery.Devices {
		deviceIds = append(deviceIds, fmt.Sprint(device.ID))
	}

	devices, err := s.deviceRepo.GetDevicesByIds(deviceIds)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot fetch devices",
		})
		log.Printf("cannot fetch devices for user %s: %v", deviceIds, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"request_id": xRequestId,
		"payload": gin.H{
			"devices": devices,
		},
	})
}

func (s *DevicesService) ActionDevices(c *gin.Context) {
	xRequestId := c.GetHeader("X-Request-ID")

	var devicesActionRequest models.DevicesActionRequest
	if err := c.ShouldBindJSON(&devicesActionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	var deviceActionResponse models.DevicesActionResponse
	deviceActionResponse.RequestID = xRequestId

	for idx := range devicesActionRequest.Payload.Devices {
		device := &devicesActionRequest.Payload.Devices[idx]

		deviceInDb, err := s.deviceRepo.GetDeviceById(device.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "cannot fetch devices",
			})
			log.Printf("cannot fetch devices for user %s: %v", device.ID, err)
			return
		}

		deviceActionResponse.Payload.Devices = append(deviceActionResponse.Payload.Devices, models.DeviceActionResponse{
			ID:           device.ID,
			CustomData:   device.CustomData,
			Capabilities: make([]models.CapabilityWithAction, 0),
		})

		capMap := make(map[string]*models.Capability)
		for i := range deviceInDb.Capabilities {
			capability := &deviceInDb.Capabilities[i]
			capMap[capability.Type] = capability
		}

		for _, capability := range device.Capabilities {
			dbCap, ok := capMap[capability.Type]
			if !ok {
				continue
			}

			if !reflect.DeepEqual(capability.State, dbCap.State) {
				// todo: update device physically
				//

				// update device's capability state in storage
				jsonIncomingState, err := json.Marshal(capability.State)
				if err != nil {
					log.Printf("Failed to marshal incoming state: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "failed to marshal incoming state",
					})
					return
				}

				err = s.deviceRepo.UpdateCapabilityStateByDeviceIdAndType(device.ID, capability.Type, jsonIncomingState)
				if err != nil {
					log.Printf("Failed to update state: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "failed to update state",
					})
					return
				}
			}

			var state map[string]string
			err = json.Unmarshal(capability.State, &state)

			newState := models.State{
				Instance: state["instance"],
				ActionResult: models.ActionResult{
					Status: "DONE",
				},
			}
			deviceActionResponse.Payload.Devices[idx].Capabilities = append(deviceActionResponse.Payload.Devices[idx].Capabilities, models.CapabilityWithAction{
				Type:  capability.Type,
				State: newState,
			})
		}
	}

	c.JSON(http.StatusOK, deviceActionResponse)
}
