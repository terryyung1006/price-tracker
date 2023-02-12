package interfaces

import "go.mongodb.org/mongo-driver/bson"

type ILatestPriceGetter interface {
	GetLatestPrice(PairTag string) (float64, int, error)
}

type IByTimestampPriceGetter interface {
	GetPriceByTimestamp(pairTag string, timestamp int) (float64, int, error)
}

type IPriceInTimestampRangeGetter interface {
	GetPriceInTimestampRange(pairTag string, timestampFrom int, timestampTo int) ([]bson.M, int, error)
}

type IAveragePriceInTimestampRangeGetter interface {
	GetPriceInTimestampRange(pairTag string, timestampFrom int, timestampTo int) (float64, int, error)
}
