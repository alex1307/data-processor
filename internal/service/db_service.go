package service

import (
	"fmt"
	"log"
	"os"
	"sync"

	yaml "gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	service *DBService
	once    sync.Once
	connect = func(filename string) {
		service = NewDBService(filename)
	}
)

func GetDBService(filename string) *DBService {
	once.Do(func() {
		connect(filename)
	})
	return service
}

type DBService struct {
	db *gorm.DB
}

type DBConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
	Database string `yaml:"database"`
	Port     int    `yaml:"port"`
	SSLMode  string `yaml:"sslmode"`
}

type Content struct {
	Config DBConfig `yaml:"db"`
}

func (c *DBConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Host, c.User, c.Password, c.Database, c.Port, c.SSLMode)
}

func (c *DBConfig) GetTablePrefix() string {
	return fmt.Sprintf("%s.", c.Schema)
}

func NewDBService(filename string) *DBService {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	// Create a variable to hold the YAML data
	var content Content
	log.Println("File is", string(data))
	yaml.Unmarshal(data, &content)

	log.Println("Config", content)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}
	config := content.Config
	dsn := config.GetConnectionString()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			NoLowerCase:   true,
			TablePrefix:   config.GetTablePrefix(), // schema name
			SingularTable: false,
		}})
	if err != nil {
		log.Fatal(err)
	}

	return &DBService{db: db}
}

func (s *DBService) GetDB() *gorm.DB {
	return s.db
}
