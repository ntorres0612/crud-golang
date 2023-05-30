package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type ProductType struct {
	ID   uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Type string `gorm:"size:255;not null;unique" json:"type"`
}

func (productType *ProductType) Prepare() {
	productType.ID = 0
	productType.Type = html.EscapeString(strings.TrimSpace(productType.Type))
}

func (productType *ProductType) Validate() error {

	if productType.Type == "" {
		return errors.New("required Type")
	}
	return nil
}

func (productType *ProductType) SaveProductType(db *gorm.DB) (*ProductType, error) {
	err := db.Debug().Model(&ProductType{}).Create(&productType).Error
	if err != nil {
		return &ProductType{}, err
	}

	return productType, nil
}

func (productType *ProductType) FindAllProductType(db *gorm.DB) (*[]ProductType, error) {
	drugs := []ProductType{}
	err := db.Debug().Model(&ProductType{}).Limit(100).Find(&drugs).Error
	if err != nil {
		return &[]ProductType{}, err
	}

	return &drugs, nil
}

func (productType *ProductType) FindProductTypeByID(db *gorm.DB, id uint64) (*ProductType, error) {
	err := db.Debug().Model(&ProductType{}).Where("id = ?", id).Take(&productType).Error
	if err != nil {
		return &ProductType{}, err
	}

	return productType, nil
}

func (productType *ProductType) CheckExistCreate(db *gorm.DB) bool {
	err := db.Debug().Model(&ProductType{}).Where("name = ?", productType.Type).Take(&productType).Error
	return err == nil
}

func (productType *ProductType) CheckExistUpdate(db *gorm.DB) bool {
	err := db.Debug().Model(&ProductType{}).Where("name = ? AND id <> ?", productType.Type, productType.ID).Take(&productType).Error
	return err == nil
}

func (productType *ProductType) UpdateProductType(db *gorm.DB) (*ProductType, error) {

	err := db.Debug().Model(&ProductType{}).Where("id = ?", productType.ID).Updates(ProductType{
		Type: productType.Type,
	}).Error
	if err != nil {
		return &ProductType{}, err
	}

	return productType, nil
}

func (productType *ProductType) DeleteProductType(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&ProductType{}).Where("id = ?", id).Take(&ProductType{}).Delete(&ProductType{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("ProductType not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
