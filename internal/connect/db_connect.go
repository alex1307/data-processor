package connect

import (
	dbmodel "data-processor/internal/model/db"
	"data-processor/utils"
	"fmt"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	Port     string `yaml:"port"`
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
	logrus.Info("Connecting to database...", connectionString)
	switch database {
	case "postgres":
		db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				NoLowerCase:   true,
				TablePrefix:   config.GetTablePrefix(), // schema name
				SingularTable: false,
			}})

	case "sqlite":

		db, err = gorm.Open(sqlite.Open(connectionString), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

		// Migrate the database schema

	default:
		err = fmt.Errorf("database %s not supported", config.GetDatabase())
	}
	if err != nil {
		panic("failed to connect database")
	}
	if db.Migrator().HasTable(&dbmodel.BasicData{}) {
		logrus.Info("Table exists")
	} else {
		logrus.Info("Table does not exist")
		db.AutoMigrate(&dbmodel.BasicData{})
	}

	if db.Migrator().HasTable(&dbmodel.ConsumptionData{}) {
		logrus.Info("Table exists")
	} else {
		logrus.Info("Table does not exist")
		db.AutoMigrate(&dbmodel.ConsumptionData{})
	}

	if db.Migrator().HasTable(&dbmodel.DetailsData{}) {
		logrus.Info("Table exists")
	} else {
		logrus.Info("Table does not exist")
		db.AutoMigrate(&dbmodel.DetailsData{})
	}

	if db.Migrator().HasTable(&dbmodel.PriceData{}) {
		logrus.Info("Table exists")
	} else {
		logrus.Info("Table does not exist")
		db.AutoMigrate(&dbmodel.PriceData{})
	}

	if db.Migrator().HasTable(&dbmodel.ChangeLogData{}) {
		logrus.Info("Table exists")
	} else {
		logrus.Info("Table does not exist")
		db.AutoMigrate(&dbmodel.ChangeLogData{})
	}

	if db.Migrator().HasTable(&dbmodel.IDData{}) {
		logrus.Info("Table exists")
	} else {
		logrus.Info("Table does not exist")
		db.AutoMigrate(&dbmodel.IDData{})
	}
	if err != nil {
		logrus.Error("Failed to migrate database", err)
		panic("failed to sync with database")
	}

	return Database{db}
}

type Content struct {
	Config PostgresConfig `yaml:"db"`
}

func (c PostgresConfig) GetConnectionString() string {
	host := utils.GetEnv("DB_HOST", c.Host)
	port := utils.GetEnv("DB_PORT", c.Port)
	user := utils.GetEnv("DB_USER", c.User)
	password := utils.GetEnv("DB_PASSWORD", c.Password)
	dbname := utils.GetEnv("DB_NAME", c.Database)
	ssl := utils.GetEnv("SSL_MODE", c.SSLMode)
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, ssl)
}

func (c PostgresConfig) GetTablePrefix() string {
	schema := utils.GetEnv("DB_SCHEMA", c.Schema)
	return fmt.Sprintf("%s.", schema)
}

func (c PostgresConfig) GetDatabase() string {
	return "postgres"
}

func GetPosgresConfig(filename string) Config {
	data, err := os.ReadFile(filename)
	if err != nil {
		logrus.Error("Failed to read file", err)
	}
	// Create a variable to hold the YAML data
	var content Content
	logrus.Info("Filename is: ", filename)
	yaml.Unmarshal(data, &content)

	logrus.Debug("YAML content is: ", content)
	if err != nil {
		logrus.Error("Failed to unmarshal YAML: ", err)
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
