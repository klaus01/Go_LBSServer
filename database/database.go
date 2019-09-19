package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/klaus01/Go_LBSServer/config"
	"github.com/klaus01/Go_LBSServer/models"
)

var db *mongo.Database

// Init 初始化
func Init() {
	c := config.GetConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.GetString("mongo.uri")))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	db = client.Database(c.GetString("mongo.dbName"))

	checkSmsCodes()
	checkOrders()
}

// Context 获取队列
func Context() context.Context {
	return context.Background()
}

// GetDB 获取数据库
func GetDB() *mongo.Database {
	return db
}

func checkSmsCodes() {
	collectionName := models.TableNameSmsCode
	collection := db.Collection(collectionName)
	indexes := collection.Indexes()
	cur, err := indexes.List(Context())
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(Context())
	createAtIndexExist, phoneNumberIndexExist := false, false
	for cur.Next(Context()) {
		index := bson.D{}
		err := cur.Decode(&index)
		if err != nil {
			log.Fatal(err)
		}
		var name interface{}
		for _, item := range index {
			if item.Key == "name" {
				name = item.Value
				break
			}
		}
		if name == "createAt_1" {
			createAtIndexExist = true
			log.Println(collectionName, "索引 createAt 存在")
		} else if name == "phoneNumber_1" {
			phoneNumberIndexExist = true
			log.Println(collectionName, "索引 phoneNumber 存在")
		}
	}

	if !createAtIndexExist {
		log.Println(collectionName, "创建 createAt 索引")
		expireAfterSeconds := config.GetConfig().GetInt32("smsCodeExpireAfterSeconds")
		_, err := indexes.CreateOne(Context(), mongo.IndexModel{Keys: bsonx.Doc{{Key: "createAt", Value: bsonx.Int32(1)}}, Options: &options.IndexOptions{ExpireAfterSeconds: &expireAfterSeconds}})
		if err != nil {
			log.Fatal(err)
		}
	}
	if !phoneNumberIndexExist {
		log.Println(collectionName, "创建 phoneNumber 索引")
		unique := true
		_, err := indexes.CreateOne(Context(), mongo.IndexModel{Keys: bsonx.Doc{{Key: "phoneNumber", Value: bsonx.Int32(1)}}, Options: &options.IndexOptions{Unique: &unique}})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func checkOrders() {
	collectionName := "orders"
	collection := db.Collection(collectionName)
	indexes := collection.Indexes()
	cur, err := indexes.List(Context())
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(Context())
	orderIDIndexExist := false
	for cur.Next(Context()) {
		index := bson.D{}
		err := cur.Decode(&index)
		if err != nil {
			log.Fatal(err)
		}
		var name interface{}
		for _, item := range index {
			if item.Key == "name" {
				name = item.Value
				break
			}
		}
		if name == "orderId_1" {
			orderIDIndexExist = true
			log.Println(collectionName, "索引 orderId 存在")
		}
	}

	if !orderIDIndexExist {
		log.Println(collectionName, "创建 orderId 索引")
		unique := true
		_, err := indexes.CreateOne(Context(), mongo.IndexModel{Keys: bsonx.Doc{{Key: "orderId", Value: bsonx.Int32(1)}}, Options: &options.IndexOptions{Unique: &unique}})
		if err != nil {
			log.Fatal(err)
		}
	}
}
