package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Store struct {
	ID   uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:255;not null;unique" json:"name"`
}

func (store *Store) Prepare() {
	store.ID = 0
	store.Name = html.EscapeString(strings.TrimSpace(store.Name))
}

func (store *Store) Validate() error {

	if store.Name == "" {
		return errors.New("required Type")
	}
	return nil
}

func (store *Store) SaveStore(db *gorm.DB) (*Store, error) {
	err := db.Debug().Model(&Store{}).Create(&store).Error
	if err != nil {
		return &Store{}, err
	}

	return store, nil
}

func (store *Store) FindAllStore(db *gorm.DB) (*[]Store, error) {
	drugs := []Store{}
	err := db.Debug().Model(&Store{}).Limit(100).Find(&drugs).Error
	if err != nil {
		return &[]Store{}, err
	}

	return &drugs, nil
}

func (store *Store) FindStoreByID(db *gorm.DB, id uint64) (*Store, error) {
	err := db.Debug().Model(&Store{}).Where("id = ?", id).Take(&store).Error
	if err != nil {
		return &Store{}, err
	}

	return store, nil
}

func (store *Store) CheckExistCreate(db *gorm.DB) bool {
	err := db.Debug().Model(&Store{}).Where("name = ?", store.Name).Take(&store).Error
	return err == nil
}

func (store *Store) CheckExistUpdate(db *gorm.DB) bool {
	err := db.Debug().Model(&Store{}).Where("name = ? AND id <> ?", store.Name, store.ID).Take(&store).Error
	return err == nil
}

func (store *Store) UpdateStore(db *gorm.DB) (*Store, error) {

	err := db.Debug().Model(&Store{}).Where("id = ?", store.ID).Updates(Store{
		Name: store.Name,
	}).Error
	if err != nil {
		return &Store{}, err
	}

	return store, nil
}

func (store *Store) DeleteStore(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&Store{}).Where("id = ?", id).Take(&Store{}).Delete(&Store{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Store not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
