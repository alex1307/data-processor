package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Car struct {
	ID        string    `gorm:"column:id"`
	Make      string    `gorm:"column:make"`
	Model     string    `gorm:"column:model"`
	Engine    string    `gorm:"column:engine"`
	Gearbox   string    `gorm:"column:gearbox"`
	Millage   int       `gorm:"column:millage"`
	Year      int       `gorm:"column:year"`
	Source    string    `gorm:"column:source"`
	Promoted  bool      `gorm:"column:promoted"`
	Sold      bool      `gorm:"column:sold"`
	Contact   string    `gorm:"column:contact"`
	ScrapedOn time.Time `gorm:"column:scraped_on"`
	CreatedOn time.Time `gorm:"column:created_on"`
}

func main() {
	dsn := "host=localhost user=db_data password=db_data dbname=ayagasha port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var cars []Car
	result := db.Table("db_data.DATA_SNAPSHOT").Find(&cars)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	for _, car := range cars {
		fmt.Println("ID:", car)
		// ... Print other car fields
		fmt.Println("----------------------")
	}

	car := Car{
		ID:        "123457",
		Make:      "VW",
		Model:     "Golf",
		Engine:    "2.0L",
		Gearbox:   "Automatic",
		Millage:   50000,
		Year:      2020,
		Source:    "Dealer",
		Promoted:  false,
		Sold:      false,
		Contact:   "John Doe",
		ScrapedOn: time.Now(),
	}

	create_result := db.Table("db_data.DATA_SNAPSHOT").Create(&car)
	if create_result.Error != nil {
		log.Fatal(create_result.Error)
	}
	// Close the database connection
	dbSQL, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	dbSQL.Close()
}
