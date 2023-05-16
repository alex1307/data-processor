package main

import (
	dbmodel "data-processor/internal/model/db"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Table struct {
	TableName string `gorm:"column:tablename"`
}

func listTables(db *gorm.DB) ([]Table, error) {
	var tables []Table
	err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'db_data'").Scan(&tables).Error
	if err != nil {
		return nil, err
	}
	return tables, nil
}

func main() {
	fmt.Println("Hello, World!") // print "Hello, World!"
	dsn := "host=localhost user=db_data password=db_data dbname=ayagasha port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			NoLowerCase:   true,
			TablePrefix:   "db_data.", // schema name
			SingularTable: false,
		}})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database", db.Name())
	tables, err := listTables(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, table := range tables {
		fmt.Println(table.TableName)
	}
	err = db.AutoMigrate(
		&dbmodel.Equipment{},
		&dbmodel.Details{},
		&dbmodel.Status{},
		&dbmodel.Equipment{})
	if err != nil {
		log.Fatal(err)
	}
	// Close the database connection
	dbSQL, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	dbSQL.Close()
}
