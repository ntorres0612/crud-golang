package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Truck struct {
	ID           uint64 `gorm:"primary_key;auto_increment" json:"id"`
	LicensePlate string `gorm:"size:255;not null;unique" json:"license_plate"`
}

func (truck *Truck) Prepare() {
	truck.ID = 0
	truck.LicensePlate = html.EscapeString(strings.TrimSpace(truck.LicensePlate))
}

func (truck *Truck) Validate() error {

	if truck.LicensePlate == "" {
		return errors.New("required LicensePlate")
	}
	return nil
}

func (truck *Truck) SaveTruck(db *gorm.DB) (*Truck, error) {
	err := db.Debug().Model(&Truck{}).Create(&truck).Error
	if err != nil {
		return &Truck{}, err
	}

	return truck, nil
}

func (truck *Truck) FindAllTruck(db *gorm.DB) (*[]Truck, error) {
	drugs := []Truck{}
	err := db.Debug().Model(&Truck{}).Limit(100).Find(&drugs).Error
	if err != nil {
		return &[]Truck{}, err
	}

	return &drugs, nil
}

func (truck *Truck) FindTruckByID(db *gorm.DB, id uint64) (*Truck, error) {
	err := db.Debug().Model(&Truck{}).Where("id = ?", id).Take(&truck).Error
	if err != nil {
		return &Truck{}, err
	}

	return truck, nil
}

func (truck *Truck) CheckExistCreate(db *gorm.DB) bool {
	err := db.Debug().Model(&Truck{}).Where("name = ?", truck.LicensePlate).Take(&truck).Error
	return err == nil
}

func (truck *Truck) CheckExistUpdate(db *gorm.DB) bool {
	err := db.Debug().Model(&Truck{}).Where("name = ? AND id <> ?", truck.LicensePlate, truck.ID).Take(&truck).Error
	return err == nil
}

func (truck *Truck) UpdateTruck(db *gorm.DB) (*Truck, error) {

	err := db.Debug().Model(&Truck{}).Where("id = ?", truck.ID).Updates(Truck{
		LicensePlate: truck.LicensePlate,
	}).Error
	if err != nil {
		return &Truck{}, err
	}

	return truck, nil
}

func (truck *Truck) DeleteTruck(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&Truck{}).Where("id = ?", id).Take(&Truck{}).Delete(&Truck{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Truck not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
