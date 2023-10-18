package storage

import (
	"NatsStream_Service/internal/cashe"
	"NatsStream_Service/internal/model"
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Storage struct {
	db *gorm.DB
}

type Content struct {
	Order_uid string `gorm:"primary key"`
	Json      string
}

func NewDB(storagePath string) (Storage, error) {
	db, err := gorm.Open(postgres.Open(storagePath))
	if err != nil {
		return Storage{}, err
	}
	err = db.AutoMigrate(&Content{})
	if err != nil {
		return Storage{}, err
	}
	fmt.Println("Create db")
	return Storage{db}, nil
}

func (s *Storage) SaveOrder(orderToSave model.Order_client) error {
	fmt.Println("Inserting in db")
	jsonToSave, err := json.MarshalIndent(orderToSave, "", " ")
	if err != nil {
		return err
	}
	res := s.db.Create(Content{orderToSave.Order_uid, string(jsonToSave)})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *Storage) UploadCashe(cashe *cashe.Cashe) error {
	var allDataFromTable []Content
	res := s.db.Find(&allDataFromTable) // limit + skip pagination gorm
	if res.Error != nil {
		return res.Error
	}
	for _, row := range allDataFromTable {
		var jsonOrder model.Order_client
		err := json.Unmarshal([]byte(row.Json), &jsonOrder)
		if err != nil {
			log.Print(err)
		}
		cashe.InsertToCashe(jsonOrder)
	}
	return nil
}
