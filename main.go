package main

import (
	"log"
	"price-tracker/controller"
	"price-tracker/cronjob"
	"price-tracker/lib/utils"
	"price-tracker/repository"
	"price-tracker/service"
	"price-tracker/variable"
	"price-tracker/worker"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	var mu sync.Mutex

	//instantiate mongoDB client
	mongoDBClient, err := utils.CreateMongoDBClient()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	utils.MongoDBClientInstance = mongoDBClient

	//instantiate API interface
	pricePairRepository := repository.CreatePricePairMongoDBRepository(mongoDBClient)
	pricePairService := service.CreatePricePairService(pricePairRepository)
	service.PricePairServiceInstance = pricePairService.(*service.PricePairService)

	repository.TmpPrice = map[string]map[int]float64{"btcusd": {}}

	//start cronjob workers
	PriceTracker := worker.PriceTracker{
		Mu: &mu,
	}
	PriceTrackerRetryer := worker.PriceTrackRetryer{
		Mu: &mu,
	}
	PriceAllocator := worker.PriceAllocator{
		Mu: &mu,
	}
	PriceAllocateRetryer := worker.PriceAllocateRetryer{
		Mu: &mu,
	}
	cronjob.RunCronJob(&PriceTracker, time.Duration(variable.UpdateInterval)*time.Second)
	cronjob.RunCronJob(&PriceTrackerRetryer, time.Duration(variable.UpdateInterval)*time.Second)
	cronjob.RunCronJob(&PriceAllocator, time.Duration(variable.UpdateInterval*5)*time.Second)
	cronjob.RunCronJob(&PriceAllocateRetryer, time.Duration(variable.UpdateInterval)*time.Second)

	router := gin.Default()

	router.GET("/api/last_price", controller.GetCurrentPrice)
	router.GET("/api/price_by_timestamp", controller.GetPriceByTimestamp)
	router.GET("/api/average_price_in_range", controller.GetAveragePrice)

	router.Run(":8080")
}
