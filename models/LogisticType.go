package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type LogisticType struct {
	ID   uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Type string `gorm:"size:255;not null;unique" json:"type"`
	Tag  string `gorm:"size:255;not null;unique" json:"tag"`
}

func (logisticType *LogisticType) Prepare() {
	logisticType.ID = 0
	logisticType.Type = html.EscapeString(strings.TrimSpace(logisticType.Type))
}

func (logisticType *LogisticType) Validate() error {

	if logisticType.Type == "" {
		return errors.New("required Type")
	}
	return nil
}

func (logisticType *LogisticType) SaveLogisticType(db *gorm.DB) (*LogisticType, error) {
	err := db.Debug().Model(&LogisticType{}).Create(&logisticType).Error
	if err != nil {
		return &LogisticType{}, err
	}

	return logisticType, nil
}

func (logisticType *LogisticType) FindAllLogisticType(db *gorm.DB) (*[]LogisticType, error) {
	drugs := []LogisticType{}
	err := db.Debug().Model(&LogisticType{}).Limit(100).Find(&drugs).Error
	if err != nil {
		return &[]LogisticType{}, err
	}

	return &drugs, nil
}

func (logisticType *LogisticType) FindLogisticTypeByID(db *gorm.DB, id uint64) (*LogisticType, error) {
	err := db.Debug().Model(&LogisticType{}).Where("id = ?", id).Take(&logisticType).Error
	if err != nil {
		return &LogisticType{}, err
	}

	return logisticType, nil
}

func (logisticType *LogisticType) UpdateLogisticType(db *gorm.DB) (*LogisticType, error) {

	err := db.Debug().Model(&LogisticType{}).Where("id = ?", logisticType.ID).Updates(LogisticType{
		Type: logisticType.Type,
	}).Error
	if err != nil {
		return &LogisticType{}, err
	}

	return logisticType, nil
}

func (logisticType *LogisticType) DeleteLogisticType(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&LogisticType{}).Where("id = ?", id).Take(&LogisticType{}).Delete(&LogisticType{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("LogisticType not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
