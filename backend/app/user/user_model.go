package user

import (
	"bytes"
	"context"
	"joosum-backend/pkg/db"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModel struct{}

// User 스키마 정의
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserId    string             `bson:"user_id"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Social    string             `bson:"social"`
	Gender    string             `bson:"gender"`
	Age       uint8              `bson:"age"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// FindUserByEmail은 주어진 이메일을 가진 사용자를 찾아 반환합니다.
func (*UserModel) FindUserByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": email}
	user := &User{}

	err := db.UserCollection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (*UserModel) CreatUser(userInfo User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uniqueKey := make(chan string)

	// TO DO : uid generater 만들기
	go func() {
		s := "User-"
		buf := bytes.NewBufferString(s)
		uid := uuid.New()
		buf.WriteString(uid.String())
		uniqueKey <- buf.String()
	}()

	uid := <-uniqueKey
	newUserInfo := &User{
		UserId:    uid,
		Name:      userInfo.Name,
		Email:     userInfo.Email,
		Social:    userInfo.Social,
		Gender:    userInfo.Gender,
		Age:       userInfo.Age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := db.UserCollection.InsertOne(ctx, newUserInfo)
	if err != nil {
		return nil, err
	}

	newUserInfo.ID = result.InsertedID.(primitive.ObjectID)

	return newUserInfo, nil
}

func (*UserModel) FindUserById(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": id}
	user := &User{}

	err := db.UserCollection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (*UserModel) FindUsers() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var users []*User

	cursor, err := db.UserCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		users = append(users, &user)
	}

	return users, nil
}
