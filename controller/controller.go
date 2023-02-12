package controller

import (
	"fmt"
	netHttp "net/http"
	"price-tracker/lib/http"
	"price-tracker/lib/number"
	"price-tracker/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCurrentPrice(ctx *gin.Context) {
	pairTag := ctx.Query("pair_tag")
	if pairTag != "btcusd" {
		http.ResponseJson(ctx, nil, fmt.Errorf("only btcusd is supported now"), netHttp.StatusBadRequest)
		return
	}
	result, statusCode, err := service.PricePairServiceInstance.GetLatestPrice(pairTag)
	if err != nil {
		http.ResponseJson(ctx, -1, err, statusCode)
		return
	}
	http.ResponseJson(ctx, result, nil, 200)
}

func GetPriceByTimestamp(ctx *gin.Context) {
	pairTag := ctx.Query("pair_tag")
	if pairTag != "btcusd" {
		http.ResponseJson(ctx, nil, fmt.Errorf("only btcusd is supported now"), netHttp.StatusBadRequest)
		return
	}
	timestamp := ctx.Query("timestamp")
	var timestampInt int
	if timestampInt = checkTimestamp(ctx, timestamp); timestampInt < 0 {
		return
	}
	result, statusCode, err := service.PricePairServiceInstance.GetPriceByTimestamp(pairTag, timestampInt)
	if err != nil {
		http.ResponseJson(ctx, -1, err, statusCode)
		return
	}
	http.ResponseJson(ctx, result, nil, 200)
}

func GetAveragePrice(ctx *gin.Context) {
	pairTag := ctx.Query("pair_tag")
	if pairTag != "btcusd" {
		http.ResponseJson(ctx, nil, fmt.Errorf("only btcusd is supported now"), netHttp.StatusBadRequest)
		return
	}
	timestampFrom := ctx.Query("timestamp_from")
	var timestampFromInt int
	if timestampFromInt = checkTimestamp(ctx, timestampFrom); timestampFromInt < 0 {
		return
	}
	timestampTo := ctx.Query("timestamp_to")
	var timestampToInt int
	if timestampToInt = checkTimestamp(ctx, timestampTo); timestampToInt < 0 {
		return
	}

	if timestampFrom > timestampTo {
		http.ResponseJson(ctx, nil, fmt.Errorf("input timestamp_from > timestamp_to"), netHttp.StatusBadRequest)
		return
	}

	result, statusCode, err := service.PricePairServiceInstance.GetPriceInTimestampRange(pairTag, timestampFromInt, timestampToInt)
	if err != nil {
		http.ResponseJson(ctx, nil, err, statusCode)
		return
	}
	http.ResponseJson(ctx, result, nil, 200)
}

func checkTimestamp(ctx *gin.Context, timestamp string) int {
	if len(timestamp) != 10 && len(timestamp) != 13 {
		http.ResponseJson(ctx, nil, fmt.Errorf("timestamp [%s] length of digit invalid", timestamp), netHttp.StatusBadRequest)
		return -1
	}
	timestampInt, err := strconv.Atoi(timestamp)
	if err != nil {
		http.ResponseJson(ctx, nil, fmt.Errorf("timestamp [%s] is not valid", timestamp), netHttp.StatusBadRequest)
		return -1
	}
	if timestampInt < 0 {
		http.ResponseJson(ctx, nil, fmt.Errorf("timestamp [%s] < 0", timestamp), netHttp.StatusBadRequest)
		return -1
	}
	timestampInt, err = number.MiliSecondToSecondTimestamp(timestampInt)
	if err != nil {
		http.ResponseJson(ctx, nil, fmt.Errorf("timestamp [%s] cannot be parsed as timestamp in second", timestamp), netHttp.StatusBadRequest)
		return -1
	}
	return timestampInt
}
