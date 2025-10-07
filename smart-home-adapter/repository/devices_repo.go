package repository

import (
	"devices-api/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type devicesRepo struct {
	db *gorm.DB
}

func (r *devicesRepo) GetByID(id string) (*models.Device, error) {
	var device models.Device

	if err := r.db.Preload("Capabilities").
		Preload("Properties").
		First(&device, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &device, nil
}

func (r *devicesRepo) GetByIDs(ids []string) ([]models.Device, error) {
	if len(ids) == 0 {
		return []models.Device{}, nil
	}

	var devices []models.Device

	if err := r.db.
		Preload("Capabilities").
		Preload("Properties").
		Where("id IN ?", ids).
		Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

func (r *devicesRepo) GetByUserID(userID string) ([]models.Device, error) {
	var devices []models.Device

	if err := r.db.
		Preload("Capabilities").
		Preload("Properties").
		Where("user_id = ?", userID).
		Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

func (r *devicesRepo) Save(device *models.Device) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(device).Error; err != nil {
			return fmt.Errorf("failed to save device: %w", err)
		}

		for i := range device.Capabilities {
			device.Capabilities[i].DeviceID = device.ID
		}
		if len(device.Capabilities) > 0 {
			if err := tx.Create(&device.Capabilities).Error; err != nil {
				return fmt.Errorf("failed to save capabilities: %w", err)
			}
		}

		for i := range device.Properties {
			device.Properties[i].DeviceID = device.ID
		}
		if len(device.Properties) > 0 {
			if err := tx.Create(&device.Properties).Error; err != nil {
				return fmt.Errorf("failed to save properties: %w", err)
			}
		}

		return nil
	})
}

func (r *devicesRepo) Update(device *models.Device) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Device{}).Where("id = ?", device.ID).Updates(device).Error; err != nil {
			return err
		}

		//if err := tx.Model(&models.DeviceInfo{}).
		//	Where("device_id = ?", device.ID).
		//	Updates(&device.DeviceInfo).Error; err != nil {
		//	return err
		//}

		if err := tx.Where("device_id = ?", device.ID).Delete(&models.Capability{}).Error; err != nil {
			return err
		}
		for i := range device.Capabilities {
			device.Capabilities[i].DeviceID = device.ID
		}
		if len(device.Capabilities) > 0 {
			if err := tx.Create(&device.Capabilities).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("device_id = ?", device.ID).Delete(&models.Property{}).Error; err != nil {
			return err
		}
		for i := range device.Properties {
			device.Properties[i].DeviceID = device.ID
		}
		if len(device.Properties) > 0 {
			if err := tx.Create(&device.Properties).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *devicesRepo) Upsert(device *models.Device) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Upsert main device
		if err := tx.Clauses(clause.OnConflict{
			UpdateAll: true, // update all fields if conflict
		}).Create(device).Error; err != nil {
			return err
		}

		// Upsert DeviceInfo
		//device.DeviceInfo.DeviceID = device.ID
		//if err := tx.Clauses(clause.OnConflict{
		//	Columns:   []clause.Column{{Name: "device_id"}},
		//	UpdateAll: true,
		//}).Create(&device.DeviceInfo).Error; err != nil {
		//	return err
		//}

		// Replace Capabilities and Properties
		if err := tx.Where("device_id = ?", device.ID).Delete(&models.Capability{}).Error; err != nil {
			return err
		}
		for i := range device.Capabilities {
			device.Capabilities[i].DeviceID = device.ID
		}
		if len(device.Capabilities) > 0 {
			if err := tx.Create(&device.Capabilities).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("device_id = ?", device.ID).Delete(&models.Property{}).Error; err != nil {
			return err
		}
		for i := range device.Properties {
			device.Properties[i].DeviceID = device.ID
		}
		if len(device.Properties) > 0 {
			if err := tx.Create(&device.Properties).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *devicesRepo) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func NewDevicesRepo(db *gorm.DB) DevicesRepository {
	return &devicesRepo{db: db}
}

//
//func (r *DevicesRepo) GetDevicesByUserId(userID string) ([]models.Device, error) {
//	var devices []models.Device
//
//	if err := r.db.
//		Preload("Capabilities").
//		Preload("Properties").
//		Where("user_id = ?", userID).
//		Find(&devices).Error; err != nil {
//		return nil, err
//	}
//
//	return devices, nil
//}
//
//func (r *DevicesRepo) GetDeviceById(id string) (*models.Device, error) {
//	var device models.Device
//
//	if err := r.db.
//		Preload("Capabilities").
//		Preload("Properties").
//		Where("id = ?", id).
//		First(&device).Error; err != nil {
//		return nil, err
//	}
//
//	return &device, nil
//}
//
//func (r *DevicesRepo) GetDevicesByIds(ids []string) ([]models.Device, error) {
//	var devices []models.Device
//
//	if err := r.db.
//		Preload("Capabilities").
//		Preload("Properties").
//		Where("id IN ?", ids).
//		Find(&devices).Error; err != nil {
//		return nil, err
//	}
//
//	return devices, nil
//}
//
//func (r *DevicesRepo) SaveOrUpdateDevice(device *models.Device) error {
//	return r.db.Transaction(func(tx *gorm.DB) error {
//		if err := tx.Save(device).Error; err != nil {
//			return err
//		}
//
//		if err := tx.Model(device).Association("Capabilities").Replace(device.Capabilities); err != nil {
//			return err
//		}
//
//		return nil
//	})
//	//
//	//if err := r.storage.Save(device).Error; err != nil {
//	//	return err
//	//}
//	//
//	//if err := r.storage.Model(device).Association("Capabilities").Replace(device.Capabilities); err != nil {
//	//	return err
//	//}
//
//	//for _, capability := range device.Capabilities {
//	//	if err := r.storage.Save(capability).Error; err != nil {
//	//		return err
//	//	}
//	//}
//}
//
//func (r *DevicesRepo) UpdateCapabilityStateByDeviceIdAndType(deviceId string, capabilityType string, newState []byte) error {
//	return r.db.
//		Model(&models.Capability{}).
//		Where("device_id = ? AND type = ?", deviceId, capabilityType).
//		Update("state", newState).
//		Error
//}
