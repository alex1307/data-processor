package main

import (
	"context"
	"data-processor/internal/connect"
	KafkaConsumer "data-processor/internal/kafka"
	service "data-processor/internal/service/db"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
)

func main() {
	runCPUProfile()
	defer runMemProfile() // You might want to check this call. It seems repetitive.

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure context cancellation when main exits

	// Database connection setup
	config := connect.GetPosgresConfig("resources/config/postgres_config.yml")
	dbConnection := connect.ConnectToDatabase(config)

	// Kafka consumer setup
	basicDataService := service.NewDetailsDataService(dbConnection)
	priceDataService := service.NewPriceDataService(dbConnection)
	consumptionDataService := service.NewConsumptionDataService(dbConnection)

	basicDataProcessor := KafkaConsumer.NewDataProcessor("base_info", "test", basicDataService)
	priceDataProcessor := KafkaConsumer.NewDataProcessor("price_info", "test", priceDataService)
	consumptionDataProcessor := KafkaConsumer.NewDataProcessor("consumption_info", "test", consumptionDataService)

	basicDataConsumer := KafkaConsumer.NewKafkaConsumer("base_info", "basic_data_group2", []string{"127.0.0.1:9094"}, basicDataProcessor)
	priceDataConsumer := KafkaConsumer.NewKafkaConsumer("price_info", "price_data_group2", []string{"127.0.0.1:9094"}, priceDataProcessor)
	consumptionDataConsumer := KafkaConsumer.NewKafkaConsumer("consumption_info", "consumption_data_group2", []string{"127.0.0.1:9094"}, consumptionDataProcessor)

	// Run consumer in a goroutine
	go basicDataConsumer.Consume(ctx)
	go priceDataConsumer.Consume(ctx)
	go consumptionDataConsumer.Consume(ctx)

	// Wait for interrupt signal to gracefully shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Perform shutdown operations
	fmt.Println("Shutting down gracefully...")
	cancel() // Cancel context to stop consumer
	// Other cleanup code if necessary
}

func runMemProfile() {
	f, err := os.Create("memprofile")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	runtime.GC()    // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	f.Close()
}

func runCPUProfile() {
	f, err := os.Create("cpuprofile")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()
}
