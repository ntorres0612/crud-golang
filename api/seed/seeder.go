package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/ntorres0612/ionix-crud/models"
)

var users = []models.User{
	{
		Name:     "Nelson Torres",
		Email:    "ntorres@gmail.com",
		Password: "nTorres",
	},
	{
		Name:     "Jesus",
		Email:    "jesus@gmail.com",
		Password: "Matiz",
	},
}

var logisticTypes = []models.LogisticType{
	{
		Type: "Portuaria",
		Tag:  "p",
	},
	{
		Type: "Marítimita",
		Tag:  "m",
	},
}

var productTypes = []models.ProductType{
	{
		Type: "Product 1",
	},
	{
		Type: "Product 2",
	},
	{
		Type: "Product 3",
	},
	{
		Type: "Product 4",
	},
	{
		Type: "Product 5",
	},
	{
		Type: "Product 6",
	},
	{
		Type: "Product 7",
	},
	{
		Type: "Product 8",
	},
	{
		Type: "Product 9",
	},
	{
		Type: "Product 10",
	},
	{
		Type: "Product 11",
	},
	{
		Type: "Product 12",
	},
}

var trucks = []models.Truck{
	{
		LicensePlate: "ABC123",
	},
}
var stores = []models.Store{
	{
		Name: "Store 1",
	},
}
var customers = []models.Customer{
	{
		Name:           "Customer 1",
		DocumentNumber: "1234567890",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(
		&models.Truck{},
		&models.Store{},
		&models.LogisticType{},
		&models.ProductType{},
		&models.Customer{},
		&models.DeliveryProduct{},
		&models.User{},
		&models.Delivery{},
	).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(
		&models.User{},
		&models.Store{},
		&models.Truck{},
		&models.Customer{},
		&models.LogisticType{},
		&models.ProductType{},
		&models.Delivery{},
		&models.DeliveryProduct{},
	).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Delivery{}).AddForeignKey("truck_id", "trucks(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
	err = db.Debug().Model(&models.Delivery{}).AddForeignKey("store_id", "stores(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
	err = db.Debug().Model(&models.Delivery{}).AddForeignKey("logistic_type_id", "logistic_types(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
	err = db.Debug().Model(&models.Delivery{}).AddForeignKey("customer_id", "customers(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
	err = db.Debug().Model(&models.DeliveryProduct{}).AddForeignKey("delivery_id", "deliveries(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
	err = db.Debug().Model(&models.DeliveryProduct{}).AddForeignKey("product_type_id", "product_types(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

	}
	for i := range logisticTypes {
		err = db.Debug().Model(&models.LogisticType{}).Create(&logisticTypes[i]).Error
		if err != nil {
			log.Fatalf("cannot seed logistic types table: %v", err)
		}

	}
	for i := range productTypes {
		err = db.Debug().Model(&models.ProductType{}).Create(&productTypes[i]).Error
		if err != nil {
			log.Fatalf("cannot seed products type table: %v", err)
		}

	}
	for i := range stores {
		err = db.Debug().Model(&models.Store{}).Create(&stores[i]).Error
		if err != nil {
			log.Fatalf("cannot seed store type table: %v", err)
		}

	}
	for i := range trucks {
		err = db.Debug().Model(&models.Truck{}).Create(&trucks[i]).Error
		if err != nil {
			log.Fatalf("cannot seed truck type table: %v", err)
		}

	}
	for i := range customers {
		err = db.Debug().Model(&models.Customer{}).Create(&customers[i]).Error
		if err != nil {
			log.Fatalf("cannot seed customers type table: %v", err)
		}

	}
}
