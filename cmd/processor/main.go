package main

import (
	"context"
	"data-processor/internal/connect"
	KafkaConsumer "data-processor/internal/kafka"
	"data-processor/internal/logger"
	service "data-processor/internal/service/db"
	"data-processor/utils"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	// runCPUProfile()
	// defer runMemProfile() // You might want to check this call. It seems repetitive.

	// Setup context with cancellation
	sessionID := logger.GenerateSessionID()
	f, err := os.OpenFile("logs/processor.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
		panic(err)
	}
	defer f.Close()
	// Create a new instance of the custom formatter
	customFormatter := &logger.CustomTextFormatter{SessionID: sessionID}
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(io.MultiWriter(os.Stdout, f))

	logrus.SetReportCaller(true)
	// Set the custom formatter as the formatter for the logger
	logrus.SetFormatter(customFormatter)

	// Now the logger is configured for the entire application
	logrus.Info("Starting data processor...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure context cancellation when main exits

	// Database connection setup
	config := connect.GetPosgresConfig("resources/config/postgres_config.yml")
	dbConnection := connect.ConnectToDatabase(config)

	// Kafka consumer setup
	basicDataService := service.NewBasicDataService(dbConnection)
	priceDataService := service.NewPriceDataService(dbConnection)
	consumptionDataService := service.NewConsumptionDataService(dbConnection)
	detailsDataService := service.NewDetailsDataService(dbConnection)
	changeLogDataService := service.NewChangeLogDataService(dbConnection)
	idService := service.NewIDDataService(dbConnection)

	basicDataProcessor := KafkaConsumer.NewDataProcessor(basicDataService)
	priceDataProcessor := KafkaConsumer.NewDataProcessor(priceDataService)
	consumptionDataProcessor := KafkaConsumer.NewDataProcessor(consumptionDataService)
	detailedDataProcessor := KafkaConsumer.NewDataProcessor(detailsDataService)
	changeLogProcessor := KafkaConsumer.NewDataProcessor(changeLogDataService)
	idProcessor := KafkaConsumer.NewDataProcessor(idService)
	kafka_env := utils.GetEnv("KAFKA_BROKER", "127.0.0.1:9094")
	brokers := []string{kafka_env}

	basicDataConsumer := KafkaConsumer.NewKafkaConsumer("base_info", "base_info", brokers, basicDataProcessor)
	priceDataConsumer := KafkaConsumer.NewKafkaConsumer("price_info", "price_info", brokers, priceDataProcessor)
	consumptionDataConsumer := KafkaConsumer.NewKafkaConsumer("consumption_info", "consumption_info", brokers, consumptionDataProcessor)
	detailedDataConsumer := KafkaConsumer.NewKafkaConsumer("details_info", "details_info", brokers, detailedDataProcessor)
	changeLogConsumer := KafkaConsumer.NewKafkaConsumer("change_log", "change_log", brokers, changeLogProcessor)
	idConsumer := KafkaConsumer.NewKafkaConsumer("ids", "id_info", brokers, idProcessor)

	// Run consumer in a goroutine
	go basicDataConsumer.Consume(ctx)
	go priceDataConsumer.Consume(ctx)
	go consumptionDataConsumer.Consume(ctx)
	go detailedDataConsumer.Consume(ctx)
	go changeLogConsumer.Consume(ctx)
	go idConsumer.Consume(ctx)

	// Wait for interrupt signal to gracefully shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Perform shutdown operations
	logrus.Info("Shutting down gracefully...")
	cancel() // Cancel context to stop consumer
	// Other cleanup code if necessary
}

func runMemProfile() {
	f, err := os.Create("memprofile")
	if err != nil {
		logrus.Error("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	runtime.GC()    // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		logrus.Error("could not write memory profile: ", err)
	}
	f.Close()
}

func runCPUProfile() {
	f, err := os.Create("cpuprofile")
	if err != nil {
		logrus.Error("could not create CPU profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	if err := pprof.StartCPUProfile(f); err != nil {
		logrus.Error("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()
}
