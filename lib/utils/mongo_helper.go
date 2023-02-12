package utils

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBClientInstance *mongo.Client

func CreateMongoDBClient() (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	return client, err
}

func FindByRange(collection *mongo.Collection, field string, min int, max int) ([]bson.M, error) {
	filter := bson.M{field: bson.M{"$gte": min, "$lte": max}}
	documentCursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer documentCursor.Close(context.TODO())

	var results []bson.M

	for documentCursor.Next(context.TODO()) {
		var result bson.M
		err := documentCursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	err = documentCursor.Err()
	if err != nil {
		return nil, err
	}
	return results, nil
}

func FindById(id interface{}, result *bson.M, collection *mongo.Collection) error {
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(result)
	return err
}

func Insert(document bson.M, collection *mongo.Collection) error {
	_, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return err
	}
	return nil
}

func InsertOrUpdate(document bson.M, collection *mongo.Collection) error {
	val, ok := document["_id"]
	if !ok {
		return fmt.Errorf("[mongoDBHelper][InsertOrUpdate] must contain _id in input")
	}
	var result bson.M
	err := collection.FindOne(context.TODO(), bson.M{"_id": val}).Decode(&result)
	if err == nil {
		update := bson.M{"$set": document}
		_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": val}, update)
		if err != nil {
			return err
		}
	} else if err == mongo.ErrNoDocuments {
		_, err := collection.InsertOne(context.TODO(), document)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
