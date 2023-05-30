package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type DeliveryProduct struct {
	ID            uint64        `gorm:"primary_key;auto_increment" json:"id"`
	Quantity      int           `gorm:"size:255;not null;" json:"quantity"`
	Price         int           `gorm:"not null" json:"price"`
	Delivery      Delivery      `json:"delivery"`
	DeliveryID    uint32        `gorm:"not null" json:"delivery_id"`
	ProductType   []ProductType `json:"product_type"`
	ProductTypeID uint32        `gorm:"not null" json:"product_type_id"`
}

func (deliveryProduct *DeliveryProduct) Prepare() {
	deliveryProduct.ID = 0
	deliveryProduct.Quantity = 0
	deliveryProduct.Price = 0
	deliveryProduct.Delivery = Delivery{}
	deliveryProduct.ProductType = []ProductType{}
}

func (deliveryProduct *DeliveryProduct) Validate() error {

	if deliveryProduct.Quantity <= 0 {
		return errors.New("the quantity is required")
	}
	if deliveryProduct.Price <= 0 {
		return errors.New("the price is required")
	}

	return nil
}

func (deliveryProduct *DeliveryProduct) SaveDeliveryProduct(db *gorm.DB) (*DeliveryProduct, error) {

	err := db.Debug().Model(&Delivery{}).Create(&deliveryProduct).Error
	if err != nil {
		return &DeliveryProduct{}, err
	}

	return deliveryProduct, nil
}

func (deliveryProduct *DeliveryProduct) FindAllSaveDeliveryProducts(db *gorm.DB) (*[]DeliveryProduct, error) {
	deliveriesProducts := []DeliveryProduct{}
	err := db.Debug().Model(&Delivery{}).Limit(100).Find(&deliveriesProducts).Error
	if err != nil {
		return &[]DeliveryProduct{}, err
	}
	return &deliveriesProducts, nil
}

func (deliveryProduct *DeliveryProduct) FindDeliveryProductByID(db *gorm.DB, id uint64) (*DeliveryProduct, error) {
	err := db.Debug().Model(&Delivery{}).Where("id = ?", id).Take(&deliveryProduct).Error
	if err != nil {
		return &DeliveryProduct{}, err
	}

	return deliveryProduct, nil
}

func (deliveryProduct *DeliveryProduct) UpdateDeliveryProduct(db *gorm.DB) (*DeliveryProduct, error) {

	err := db.Debug().Model(&DeliveryProduct{}).Where("id = ?", deliveryProduct.ID).Updates(DeliveryProduct{
		Quantity:      deliveryProduct.Quantity,
		Price:         deliveryProduct.Price,
		DeliveryID:    deliveryProduct.DeliveryID,
		ProductTypeID: deliveryProduct.ProductTypeID,
	}).Error
	if err != nil {
		return &DeliveryProduct{}, err
	}

	return deliveryProduct, nil
}

func (deliveryProduct *DeliveryProduct) DeleteDeliveryProduct(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&DeliveryProduct{}).Where("id = ? ", id).Take(&DeliveryProduct{}).Delete(&DeliveryProduct{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Delivery not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
