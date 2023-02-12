package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"price-tracker/controller"
	"price-tracker/repository"
	"price-tracker/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type APITestSuite struct {
	suite.Suite
	t *testing.T
}

var rec *httptest.ResponseRecorder
var router *gin.Engine
var repo *mockRepository

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, &APITestSuite{t: t})
}

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) GetLatestPrice(pairTag string) (float64, int, error) {
	args := m.Called(pairTag)
	return args.Get(0).(float64), args.Int(1), args.Error(2)
}

func (m *mockRepository) GetPriceByTimestamp(pairTag string, timestamp int) (float64, int, error) {
	args := m.Called(pairTag, timestamp)
	return args.Get(0).(float64), args.Int(1), args.Error(2)
}

func (m *mockRepository) GetPriceInTimestampRange(pairTag string, timestampFrom int, timestampTo int) ([]bson.M, int, error) {
	args := m.Called(pairTag, timestampFrom, timestampTo)
	return args.Get(0).([]bson.M), args.Int(1), args.Error(2)
}

func (as *APITestSuite) SetupSuite() {
	log.Println("Setup Suite")

	router = gin.Default()

	router.GET("/api/last_price", controller.GetCurrentPrice)
	router.GET("/api/price_by_timestamp", controller.GetPriceByTimestamp)
	router.GET("/api/average_price_in_range", controller.GetAveragePrice)

	repo = new(mockRepository)
	service.PricePairServiceInstance = service.CreatePricePairService(repo).(*service.PricePairService)
}

func (as *APITestSuite) BeforeTest(suitName, testName string) {
	log.Printf("Before suite %s test %s\n", suitName, testName)
	repository.LastPrice = map[string]float64{
		"btcusd": 1,
	}
	repository.TmpPrice = map[string]map[int]float64{"btcusd": {1234567890: 2}}
	rec = httptest.NewRecorder()
}

func (as *APITestSuite) AfterTest(suitName, testName string) {
	log.Printf("After suite %s test %s\n", suitName, testName)
	repo.AssertExpectations(as.t)
	repo.ExpectedCalls = nil
}

// Test /api/last_price
func (as *APITestSuite) TestGetLatestPriceSuccessCase() {
	req, _ := http.NewRequest("GET", "/api/last_price?pair_tag=btcusd", nil)
	router.ServeHTTP(rec, req)

	as.Equal(200, rec.Code)
	as.Equal("{\"data\":1}", rec.Body.String())
}

func (as *APITestSuite) TestGetLatestPriceNoLatestPrice() {
	repository.LastPrice = map[string]float64{}
	req, _ := http.NewRequest("GET", "/api/last_price?pair_tag=btcusd", nil)
	router.ServeHTTP(rec, req)

	as.Equal(http.StatusInternalServerError, rec.Code)
	as.Equal("{\"message\":\"latest price not updated yet\"}", rec.Body.String())
}

func (as *APITestSuite) TestGetLatestPriceNoPriceTag() {
	req, _ := http.NewRequest("GET", "/api/last_price", nil)
	router.ServeHTTP(rec, req)

	as.Equal(http.StatusBadRequest, rec.Code)
	as.Equal("{\"message\":\"only btcusd is supported now\"}", rec.Body.String())
}

// Test /api/price_by_timestamp
func (as *APITestSuite) TestGetPriceByTimestampMemoryStorageSuccessCase() {
	req, _ := http.NewRequest("GET", "/api/price_by_timestamp?pair_tag=btcusd&timestamp=1234567890", nil)
	router.ServeHTTP(rec, req)

	as.Equal(200, rec.Code)
	as.Equal("{\"data\":2}", rec.Body.String())
}

func (as *APITestSuite) TestGetPriceByTimestampMongoDBSuccessCase() {
	repository.TmpPrice = map[string]map[int]float64{"btcusd": {}}
	repo.On("GetPriceByTimestamp", "btcusd", 1234567890).Return(float64(2), 200, nil)
	req, _ := http.NewRequest("GET", "/api/price_by_timestamp?pair_tag=btcusd&timestamp=1234567890", nil)
	router.ServeHTTP(rec, req)

	as.Equal(200, rec.Code)
	as.Equal("{\"data\":2}", rec.Body.String())
}

func (as *APITestSuite) TestGetPriceByTimestampApiSuccessCase() {
	repository.TmpPrice = map[string]map[int]float64{"btcusd": {}}
	repo.On("GetPriceByTimestamp", "btcusd", 1676118169).Return(float64(-1), 0, nil)
	req, _ := http.NewRequest("GET", "/api/price_by_timestamp?pair_tag=btcusd&timestamp=1676118169", nil)
	router.ServeHTTP(rec, req)
	as.Equal(200, rec.Code)
	as.Equal("{\"data\":21725.71830815818}", rec.Body.String())
}

func (as *APITestSuite) TestGetPriceByTimestampFailIfNoRecord() {
	repository.TmpPrice = map[string]map[int]float64{"btcusd": {}}
	repo.On("GetPriceByTimestamp", "btcusd", 2676118169).Return(float64(-1), 0, nil)
	req, _ := http.NewRequest("GET", "/api/price_by_timestamp?pair_tag=btcusd&timestamp=2676118169", nil)
	router.ServeHTTP(rec, req)
	as.Equal(http.StatusBadRequest, rec.Code)
	as.Equal("{\"message\":\"no price can be found at this timestamp\"}", rec.Body.String())
}

// Test /api/average_price_in_range
func (as *APITestSuite) TestGetAveragePriceSuccessCase() {
	repository.TmpPrice = map[string]map[int]float64{"btcusd": {1234567800: 1, 1234567860: 2, 1234567920: 3}}
	repo.On("GetPriceInTimestampRange", "btcusd", 1234567800, 1234567920).Return([]bson.M{}, 200, nil)
	req, _ := http.NewRequest("GET", "/api/average_price_in_range?pair_tag=btcusd&timestamp_from=1234567800&timestamp_to=1234567920", nil)
	router.ServeHTTP(rec, req)

	as.Equal(200, rec.Code)
	as.Equal("{\"data\":2}", rec.Body.String())
}

func (as *APITestSuite) TestGetAveragePriceSuccessCase2() {
	repository.TmpPrice = map[string]map[int]float64{"btcusd": {}}
	expectedResult := []bson.M{
		map[string]interface{}{
			"timestamp": 1234567800,
			"price":     1.0,
		},
		map[string]interface{}{
			"timestamp": 1234567860,
			"price":     1.1,
		},
		map[string]interface{}{
			"timestamp": 1234567920,
			"price":     1.2,
		},
	}
	repo.On("GetPriceInTimestampRange", "btcusd", 1234567800, 1234567920).Return(expectedResult, 0, nil)
	req, _ := http.NewRequest("GET", "/api/average_price_in_range?pair_tag=btcusd&timestamp_from=1234567800&timestamp_to=1234567920", nil)
	router.ServeHTTP(rec, req)

	as.Equal(200, rec.Code)
	as.Equal("{\"data\":1.0999999999999999}", rec.Body.String())
}

func (as *APITestSuite) TestGetAveragePriceSuccessCase3() {
	repository.TmpPrice = map[string]map[int]float64{"btcusd": {1234567860: 1.1, 1234567920: 1.2}}
	expectedResult := []bson.M{
		map[string]interface{}{
			"timestamp": 1234567800,
			"price":     1.0,
		},
	}
	repo.On("GetPriceInTimestampRange", "btcusd", 1234567800, 1234567920).Return(expectedResult, 0, nil)
	req, _ := http.NewRequest("GET", "/api/average_price_in_range?pair_tag=btcusd&timestamp_from=1234567800&timestamp_to=1234567920", nil)
	router.ServeHTTP(rec, req)

	as.Equal(200, rec.Code)
	as.Equal("{\"data\":1.0999999999999999}", rec.Body.String())
}

func (as *APITestSuite) TestGetAveragePriceRangeNotSatisfiable() {
	repository.TmpPrice = map[string]map[int]float64{"btcusd": {1234567860: 1.1}}
	expectedResult := []bson.M{}
	repo.On("GetPriceInTimestampRange", "btcusd", 1234567800, 1234567920).Return(expectedResult, 0, nil)
	req, _ := http.NewRequest("GET", "/api/average_price_in_range?pair_tag=btcusd&timestamp_from=1234567800&timestamp_to=1234567920", nil)
	router.ServeHTTP(rec, req)

	as.Equal(416, rec.Code)
	as.Equal("{\"message\":\"input range not satisfiable\"}", rec.Body.String())
}
