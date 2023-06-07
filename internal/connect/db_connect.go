package connect

import (
	dbmodel "data-processor/internal/model/db"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	yaml "gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	service Connect
	once    sync.Once
	connect = func(config Config) {
		service = ConnectToDatabase(config)
	}
)

func GetDBService(config Config) Connect {
	once.Do(func() {
		connect(config)
	})
	return service
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
	Database string `yaml:"database"`
	Port     int    `yaml:"port"`
	SSLMode  string `yaml:"sslmode"`
}

type InMemoryDBConfig struct {
	dsn         string
	tablePrefix string
}

func (c InMemoryDBConfig) GetConnectionString() string {
	return c.dsn
}

func (c InMemoryDBConfig) GetTablePrefix() string {
	return c.tablePrefix
}

func (c InMemoryDBConfig) GetDatabase() string {
	return "sqlite"
}

type Connect interface {
	Connect() *gorm.DB
}

type Config interface {
	GetConnectionString() string
	GetTablePrefix() string
	GetDatabase() string
}

type Database struct {
	db *gorm.DB
}

func (p Database) Connect() *gorm.DB {
	return p.db
}

func ConnectToDatabase(config Config) Connect {
	var db *gorm.DB
	var err error
	connectionString := config.GetConnectionString()
	database := config.GetDatabase()
	log.Println("Connecting to database...", connectionString)
	switch database {
	case "postgres":
		db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				NoLowerCase:   true,
				TablePrefix:   config.GetTablePrefix(), // schema name
				SingularTable: false,
			}})

	case "sqlite":

		db, err = gorm.Open(sqlite.Open(connectionString), &gorm.Config{})

		// Migrate the database schema

	default:
		err = fmt.Errorf("database %s not supported", config.GetDatabase())
	}
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&dbmodel.NotFound{},
		&dbmodel.Status{},
		&dbmodel.Vehicle{},
		&dbmodel.Equipment{})
	if err != nil {
		log.Panicln("Failed to migrate database", err)
		panic("failed to sync with database")
	}

	return Database{db}
}

type Content struct {
	Config PostgresConfig `yaml:"db"`
}

func (c PostgresConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Host, c.User, c.Password, c.Database, c.Port, c.SSLMode)
}

func (c PostgresConfig) GetTablePrefix() string {
	return fmt.Sprintf("%s.", c.Schema)
}

func (c PostgresConfig) GetDatabase() string {
	return "postgres"
}

func GetPosgresConfig(filename string) Config {
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

	return config
}

func GetInMemoryConfig() Config {
	return InMemoryDBConfig{
		dsn:         "file::memory:?cache=shared",
		tablePrefix: "",
	}
}
