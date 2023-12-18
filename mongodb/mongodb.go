package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClien(ctx context.Context, host, port, database string) (db *mongo.Database, err error) {
	mongoBDURL := fmt.Sprintf("mongodb://%s%s", host, port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoBDURL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongoDB %v", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping to mongoDB %v", err)
	}
	fmt.Println("mongoBD started")
	return client.Database(database), nil
}
