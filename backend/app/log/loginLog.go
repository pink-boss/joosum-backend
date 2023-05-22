package applog

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// LoginLog 스키마 정의
type LoginLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserId    string             `bson:"user_id"`
	CreatedAt time.Time          `bson:"created_at"`
}

// loginLogCollection은 LoginLog 모델의 컬렉션 인스턴스를 저장합니다.
var loginLogCollection *mongo.Collection

// 현재 driver 에서는 지원하지 않아서 shell 에서 직접 만듬 ㅠ
// db.createCollection("loginLog", { timeseries: { timeField: "created_at" } })

// InitLoginLogCollection은 전달된 클라이언트 인스턴스를 사용하여 loginLogCollection 변수를 설정합니다.
func InitLoginLogCollection(client *mongo.Client, dbName string) {
	loginLogCollection = client.Database(dbName).Collection("loginLog")
}

// InsertLoginLog는 로그인 로그를 삽입합니다.
func InsertLoginLog(log LoginLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := loginLogCollection.InsertOne(ctx, log)
	return err
}
