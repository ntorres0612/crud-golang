package models

import (
	"errors"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Delivery struct {
	ID              uint64            `gorm:"primary_key;auto_increment" json:"id"`
	Date            time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
	GuideNumber     int               `gorm:"not null" json:"guide_number"`
	FleetNumber     string            `gorm:"not null" json:"fleet_number"`
	Discount        float32           `gorm:"not null" json:"discount"`
	CreatedAt       time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Truck           Truck             `json:"truck"`
	TruckID         uint32            `gorm:"not null" json:"truck_id"`
	Store           Store             `json:"store"`
	StoreID         uint32            `gorm:"not null" json:"store_id"`
	LogisticType    LogisticType      `json:"logistic_type"`
	LogisticTypeID  uint32            `gorm:"not null" json:"logistic_type_id"`
	Customer        Customer          `json:"customer"`
	CustomerID      uint32            `gorm:"not null" json:"customer_id"`
	DeliveryProduct []DeliveryProduct `json:"products"`
}

func (delivery *Delivery) Prepare() {
	delivery.ID = 0
	delivery.FleetNumber = html.EscapeString(strings.TrimSpace(delivery.FleetNumber))
	delivery.Truck = Truck{}
	delivery.Store = Store{}
	delivery.LogisticType = LogisticType{}
	delivery.Customer = Customer{}
	delivery.CreatedAt = time.Now()
	delivery.UpdatedAt = time.Now()
}

func (delivery *Delivery) Validate() error {

	if delivery.GuideNumber <= 0 {
		return errors.New("the guide number is required")
	}
	if delivery.FleetNumber == "" {
		return errors.New("the fleet number is required")
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z]{3}\d{4}[a-zA-Z]$`, delivery.FleetNumber)
	if !matched {
		return errors.New("the fleet number is invalid")
	}

	if len(strconv.Itoa(delivery.GuideNumber)) != 10 {
		return errors.New("the length of guide number must be 10")
	}

	return nil
}

func (delivery *Delivery) SaveDelivery(db *gorm.DB) (*Delivery, error) {

	err := db.Debug().Model(&Delivery{}).Create(&delivery).Error
	if err != nil {
		return &Delivery{}, err
	}
	for i := range delivery.DeliveryProduct {
		fmt.Println("----------------------------------------- ")
		deliveryProduct := &DeliveryProduct{
			Quantity:      delivery.DeliveryProduct[i].Quantity,
			Price:         delivery.DeliveryProduct[i].Price,
			ProductTypeID: delivery.DeliveryProduct[i].ProductTypeID,
			DeliveryID:    uint32(delivery.ID),
		}
		err = db.Debug().Model(&DeliveryProduct{}).Create(&deliveryProduct).Error
		if err != nil {
			return &Delivery{}, err
		}

	}

	return delivery, nil
}

func (delivery *Delivery) FindAllDeliverys(db *gorm.DB) (*[]Delivery, error) {
	deliveries := []Delivery{}
	err := db.Debug().Model(&Delivery{}).Limit(100).Find(&deliveries).Error
	if err != nil {
		return &[]Delivery{}, err
	}

	return &deliveries, nil
}

func (delivery *Delivery) FindDeliveryByID(db *gorm.DB, id uint64) (*Delivery, error) {
	err := db.Debug().Model(&Delivery{}).Where("id = ?", id).Take(&delivery).Error
	if err != nil {
		return &Delivery{}, err
	}

	return delivery, nil
}

func (delivery *Delivery) UpdateDelivery(db *gorm.DB) (*Delivery, error) {

	err := db.Debug().Model(&Delivery{}).Where("id = ?", delivery.ID).Updates(Delivery{
		Date:           delivery.Date,
		GuideNumber:    delivery.GuideNumber,
		FleetNumber:    delivery.FleetNumber,
		TruckID:        delivery.TruckID,
		StoreID:        delivery.StoreID,
		LogisticTypeID: delivery.LogisticTypeID,
		CustomerID:     delivery.CustomerID,
		UpdatedAt:      time.Now()}).Error
	if err != nil {
		return &Delivery{}, err
	}

	return delivery, nil
}

func (delivery *Delivery) DeleteDelivery(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&Delivery{}).Where("id = ? ", id).Take(&Delivery{}).Delete(&Delivery{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Delivery not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
