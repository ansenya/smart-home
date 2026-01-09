package models

type ConfirmPairingRequest struct {
	Code       string `json:"code" binding:"required"`
	DeviceUID  string `json:"device_uid" binding:"required"`
	MacAddress string `json:"mac_address" binding:"required"`
}
