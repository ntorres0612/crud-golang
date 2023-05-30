package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Customer struct {
	ID             uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Name           string `gorm:"size:255;not null;" json:"name"`
	DocumentNumber string `gorm:"size:255;not null;unique" json:"document_number"`
}

func (customer *Customer) Prepare() {
	customer.ID = 0
	customer.Name = html.EscapeString(strings.TrimSpace(customer.Name))
}

func (customer *Customer) Validate() error {

	if customer.Name == "" {
		return errors.New("required Type")
	}
	return nil
}

func (customer *Customer) SaveCustomer(db *gorm.DB) (*Customer, error) {
	err := db.Debug().Model(&Customer{}).Create(&customer).Error
	if err != nil {
		return &Customer{}, err
	}

	return customer, nil
}

func (customer *Customer) FindAllCustomer(db *gorm.DB) (*[]Customer, error) {
	drugs := []Customer{}
	err := db.Debug().Model(&Customer{}).Limit(100).Find(&drugs).Error
	if err != nil {
		return &[]Customer{}, err
	}

	return &drugs, nil
}

func (customer *Customer) FindCustomerByID(db *gorm.DB, id uint64) (*Customer, error) {
	err := db.Debug().Model(&Customer{}).Where("id = ?", id).Take(&customer).Error
	if err != nil {
		return &Customer{}, err
	}

	return customer, nil
}

func (customer *Customer) CheckExistCreate(db *gorm.DB) bool {
	err := db.Debug().Model(&Customer{}).Where("document_number = ?", customer.DocumentNumber).Take(&customer).Error
	return err == nil
}

func (customer *Customer) CheckExistUpdate(db *gorm.DB) bool {
	err := db.Debug().Model(&Customer{}).Where("document_number = ? AND id <> ?", customer.DocumentNumber, customer.ID).Take(&customer).Error
	return err == nil
}

func (customer *Customer) UpdateCustomer(db *gorm.DB) (*Customer, error) {

	err := db.Debug().Model(&Customer{}).Where("id = ?", customer.ID).Updates(Customer{
		Name: customer.Name,
	}).Error
	if err != nil {
		return &Customer{}, err
	}

	return customer, nil
}

func (customer *Customer) DeleteCustomer(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&Customer{}).Where("id = ?", id).Take(&Customer{}).Delete(&Customer{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Customer not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
