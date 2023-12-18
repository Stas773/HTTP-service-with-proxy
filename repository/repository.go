package repository

import (
	"D/Works/GO/31/models"
	"D/Works/GO/31/usecase"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepoStruct struct {
	collection *mongo.Collection
}

func NewStorage(database *mongo.Database, collection string) usecase.Storage {
	return &RepoStruct{
		collection: database.Collection(collection),
	}
}

func (rs *RepoStruct) CreateNewUser(ctx context.Context, mu models.User) error {
	stringId := uuid.NewString()
	mu.Id = stringId
	result, err := rs.collection.InsertOne(ctx, mu)
	if err != nil {
		return fmt.Errorf("failed to create user")
	}
	fmt.Println("User created with id:", result)
	return nil
}

func (rs *RepoStruct) MakeFriends(ctx context.Context, id1 string, id2 string) (error, error) {
	collection := rs.collection.Database().Collection("users")
	filter1 := bson.M{"_id": id1}
	filter2 := bson.M{"_id": id2}
	user1 := collection.FindOne(ctx, filter1)
	if user1.Err() != nil {
		return fmt.Errorf("user1 not found"), nil
	}
	user2 := collection.FindOne(ctx, filter2)
	if user2.Err() != nil {
		return nil, fmt.Errorf("user2 not found")
	}
	_, err := collection.UpdateOne(ctx, filter1, bson.M{"$addToSet": bson.M{"friends": id2}})
	if err != nil {
		return err, nil
	}
	_, err = collection.UpdateOne(ctx, filter2, bson.M{"$addToSet": bson.M{"friends": id1}})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (rs *RepoStruct) GetUser(ctx context.Context, id string) (mu models.User, err error) {
	filter := bson.M{"_id": id}

	result := rs.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return mu, fmt.Errorf("ErrEntityNotFound")
		}
		return mu, fmt.Errorf("failed to find user by id: %s due to error %v", id, err)
	}
	if err = result.Decode(&mu); err != nil {
		return mu, fmt.Errorf("failed to decode user (id: %s) from DB due to error %v", id, err)
	}
	return mu, nil
}

func (rs *RepoStruct) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	collection := rs.collection.Database().Collection("users")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("collection not find: %v", err)
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, fmt.Errorf("failed to decode user: %v", err)
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (rs *RepoStruct) GetUserFriends(ctx context.Context, id string) (friends []string, err error) {
	filter := bson.M{"_id": id}
	var mu models.User

	result := rs.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return mu.Friends, fmt.Errorf("ErrEntityNotFound")
		}
		return mu.Friends, fmt.Errorf("failed to find user by id: %s due to error %v", id, err)
	}
	if err = result.Decode(&mu); err != nil {
		return mu.Friends, fmt.Errorf("failed to decode user (id: %s) from DB due to error %v", id, err)
	}
	return mu.Friends, nil
}

func (rs *RepoStruct) UpdateUser(ctx context.Context, id string, mu models.User) error {
	collection := rs.collection.Database().Collection("users")
	filter := bson.M{"_id": id}

	user := collection.FindOne(ctx, filter)
	if user.Err() != nil {
		return fmt.Errorf("user not found")
	}
	_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
		"name": mu.Name,
		"age":  mu.Age,
	},
	})
	if err != nil {
		return fmt.Errorf("user not updated")
	}
	return nil
}

func (rs *RepoStruct) DeleteUser(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}

	result, err := rs.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query, error: %v", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}
	filter = bson.M{"friends": id}
	update := bson.M{"$pull": bson.M{"friends": id}}
	_, err = rs.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete user from friends")
	}
	return nil
}
