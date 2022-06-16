package mongo_db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	City      string             `bson:"city"`
	Company   string             `bson:"company"`
	Position  string             `bson:"position"`
	Name      string             `bson:"name"`
	InnerTel  string             `bson:"innerTel"`
	MobileTel string             `bson:"mobileTel"`
	Skype     string             `bson:"skype"`
	Photo     string             `bson:"photo"`
	Mail      string             `bson:"mail"`
	IsDeleted bool               `bson:"isdeleted, omitempty"`
}

var collection *mongo.Collection
var ctx = context.TODO()

func CreateConnection(address string, user string, password string) bool {
	credential := options.Credential{
		//AuthMechanism: "SCRAM-SHA-256",
		//AuthSource:    "authenticationDb",
		Username: user,
		Password: password,
	}
	clientOpts := options.Client().ApplyURI(address).SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return false
	}

	/*x := true
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error")
		x = false
		return x
	}

	_ = client
	collection = client.Database("test").Collection("customers")
	return x*/

	return true
}

func GetAll() ([]*User, error) {
	// passing bson.D{{}} matches all documents in the collection
	filter := bson.M{"IsDeleted": false}
	return FilterUsers(filter)
}

func CreateUser(user *User) error {
	user.ID = primitive.NewObjectID()
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Println(err)
	}
	objectID := result.InsertedID.(primitive.ObjectID)
	fmt.Println(objectID)
	return err

}

func DeleteUser(user *User) error {
	opts := options.Update().SetUpsert(true)
	id := user.ID
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"isdeleted": true}}
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func UpdateUser(user *User) error {
	opts := options.Replace().SetUpsert(true)
	if primitive.ObjectID.IsZero(user.ID) {
		return errors.New("ID must be set")
	}
	id := user.ID
	filter := bson.M{"_id": id, "isdeleted": false}
	update := user
	result, err := collection.ReplaceOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}
	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return nil
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}
	return err
}

func FilterUsers(filter interface{}) ([]*User, error) {

	var users []*User

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return users, err
	}

	for cur.Next(ctx) {
		var t User
		err := cur.Decode(&t)
		if err != nil {
			return users, err
		}

		users = append(users, &t)
	}

	if err := cur.Err(); err != nil {
		return users, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(users) == 0 {
		return users, mongo.ErrNoDocuments
	}
	return users, nil
}
