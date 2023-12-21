package mongodb

import (
	"benchmark/logic"
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI = "mongodb://localhost:27017"
	dbName   = "your_db"
	collName = "iot_snapshots"
)

func Init() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to default database: %v\n", err)
		os.Exit(1)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to default database: %v\n", err)
		os.Exit(1)
	}

	collOpts := options.CreateCollection()
	err = client.Database(dbName).CreateCollection(context.Background(), collName, collOpts)

	indexModel := mongo.IndexModel{
		Keys: bson.M{"createdAt": -1},
	}

	client.Database(dbName).Collection(collName).Indexes().CreateOne(context.Background(), indexModel)

	indexModel = mongo.IndexModel{
		Keys: bson.M{"stationId": 1},
	}

	client.Database(dbName).Collection(collName).Indexes().CreateOne(context.Background(), indexModel)

	return client
}

func InsertAll(client *mongo.Client, data [][][]logic.Snapshot) error {
	coll := client.Database(dbName).Collection(collName)

	var documents []interface{}
	for _, station := range data {
		for _, simulation := range station {
			for _, snapshot := range simulation {
				doc := bson.D{
					{Key: "stationId", Value: snapshot.StationID},
					{Key: "createdAt", Value: snapshot.CreatedAt},
					{Key: "outsideTemperature", Value: snapshot.OutsideTemperature},
					{Key: "voltage", Value: snapshot.Voltage},
					{Key: "heatingTemperature", Value: snapshot.HeatingTemperature},
					{Key: "coolingTemperature", Value: snapshot.CoolingTemperature},
				}
				documents = append(documents, doc)
			}
		}
	}

	_, err := coll.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	return nil
}

func QueryLastDayData(client *mongo.Client) ([]bson.M, error) {
	coll := client.Database(dbName).Collection(collName)

	oneDayAgo := time.Now().Add(-24 * time.Hour)
	filter := bson.M{
		"stationId": 1,
		"createdAt": bson.M{"$gte": oneDayAgo},
	}

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func QueryAverageValues(client *mongo.Client) ([]bson.M, error) {
	coll := client.Database(dbName).Collection(collName)

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "stationId", Value: 1}}}}
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$stationId"},
			{Key: "avgOutsideTemperature", Value: bson.D{{Key: "$avg", Value: "$outsideTemperature"}}},
			{Key: "avgVoltage", Value: bson.D{{Key: "$avg", Value: "$voltage"}}},
			{Key: "avgHeatingTemperature", Value: bson.D{{Key: "$avg", Value: "$heatingTemperature"}}},
			{Key: "avgCoolingTemperature", Value: bson.D{{Key: "$avg", Value: "$coolingTemperature"}}},
		}},
	}

	cursor, err := coll.Aggregate(context.Background(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
