package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User 스키마 정의
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UID       string             `bson:"uid"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// userCollection은 User 모델의 컬렉션 인스턴스를 저장합니다.
var userCollection *mongo.Collection

// InitUserCollection은 전달된 클라이언트 인스턴스를 사용하여 userCollection 변수를 설정합니다.
func InitUserCollection(client *mongo.Client, dbName string) {
	userCollection = client.Database(dbName).Collection("users")
	EnsureIndexes(userCollection)
}

// TO DO
// Index 생성, 본인의 Collection 인스턴스 변수, 해당 collection을 init 하는 함수는
// 공통으로 쓰일 것 같으니 패턴화 해서 분리해두는 것이 좋을 것 같습니다.

// email에 대한 인덱스 생성
func EnsureIndexes(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

// FindUserByEmail은 주어진 이메일을 가진 사용자를 찾아 반환합니다.
func FindUserByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": email}
	user := &User{}

	err := userCollection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
