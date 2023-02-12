package http

import (
	"fmt"
	"net/http"
	"price-tracker/lib/number"
	"price-tracker/lib/slice"
	"strconv"
)

func CoingeckoQuery(method string, httpMethod string, params map[string]string, result interface{}) error {
	baseURL := "https://api.coingecko.com/api/v3" + method
	req, err := http.NewRequest(httpMethod, baseURL, nil)
	if err != nil {
		return fmt.Errorf("[CoingeckoQuery] build request obj failed with error: %s", err.Error())
	}

	q := req.URL.Query()
	for key, val := range params {
		q.Set(key, val)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set("accept", "application/json")
	err = HttpRequest(req, result)
	if err != nil {
		return fmt.Errorf("[BlockchainHttpPost] method %s post request failed with err: %s", method, err.Error())
	}
	return nil
}

type SimplePriceResponse map[string]map[string]float64

type SimplePriceUSDLastUpdatedAtResponse struct {
	USD           float64 `json:"usd"`
	LastUpdatedAt int     `json:"last_updated_at"`
}

type SimplePriceBitcoinUSDLastUpdatedAtResponse struct {
	Bitcoin SimplePriceUSDLastUpdatedAtResponse `json:"bitcoin"`
}

type MarketChartRangeResponse struct {
	Prices [][]interface{} `json:"prices"`
}

func GetBTCUSDPriceByTimestamp(timestamp int) (int, float64, error) {
	from := strconv.Itoa(timestamp - 300)
	to := strconv.Itoa(timestamp + 300)
	payload := map[string]string{
		"vs_currency": "usd",
		"from":        from,
		"to":          to,
	}
	var resp MarketChartRangeResponse

	err := CoingeckoQuery("/coins/bitcoin/market_chart/range", "GET", payload, &resp)
	if err != nil {
		return -1, -1, nil
	}
	var resultTimestamp int
	var price float64
	if len(resp.Prices) == 0 {
		return -1, -1, nil
	} else if len(resp.Prices) > 1 {
		timeStampDiffSlice := make([]int, 0, len(resp.Prices))
		for i := 0; i < len(resp.Prices); i++ {
			iterateTimestamp, ok := resp.Prices[i][0].(float64)
			if !ok {
				return -1, -1, fmt.Errorf("[GetBTCUSDPriceByTimestamp] /market_chart/range return timestamp not int")
			}
			iterateTimestampInt := int(iterateTimestamp)
			iterateTimestampInt, err = number.MiliSecondToSecondTimestamp(iterateTimestampInt)
			if err != nil {
				return -1, -1, fmt.Errorf("[GetBTCUSDPriceByTimestamp] parse timestamp in second failed with error: %s", err.Error())
			}
			diff := timestamp - iterateTimestampInt
			if diff < 0 {
				diff *= -1
			}
			timeStampDiffSlice = append(timeStampDiffSlice, diff)
		}
		index, err := slice.IndexOfMinInt(timeStampDiffSlice)
		if err != nil {
			return -1, -1, err
		}
		resultTimestamp = timeStampDiffSlice[index]
		price = resp.Prices[index][1].(float64)
	} else {
		resultTimestamp, err = number.MiliSecondToSecondTimestamp(int(resp.Prices[0][0].(float64)))
		if err != nil {
			return -1, -1, fmt.Errorf("[GetBTCUSDPriceByTimestamp] parse timestamp in second failed with error: %s", err.Error())
		}
		price = resp.Prices[0][1].(float64)
	}

	return resultTimestamp, price, nil
}
