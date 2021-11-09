package mongodb_ops

import (
	"context"
	"cronTab/worker/config"
	terrors "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	mongoCli *mongo.Client
)

// InitMongodbConn 初始化 mongodb 连接
func InitMongodbConn() error {
	var err error

	ctx, cancel := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	defer cancel()

	// 建立连接
	mongoCli, err = mongo.Connect(context.TODO(),
		options.Client().
			ApplyURI(config.GConfig.MongodbUri).
			SetConnectTimeout(time.Duration(config.GConfig.MongodbConnectTimeout)*time.Millisecond))
	if err != nil {
		return terrors.Wrap(err, "create mongodb connection failed")
	}

	// 测试 mongodb 连接
	if err = mongoCli.Ping(ctx, nil); err != nil {
		return terrors.Wrap(err, "test mongodb connection failed")
	}

	return nil
}

// CloseMongodbConn 关闭 mongodb 连接
func CloseMongodbConn() error {
	ctx, cancel := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	defer cancel()

	if err := mongoCli.Disconnect(ctx); err != nil {
		return terrors.Wrap(err, "disconnect mongodb failed")
	}
	return nil
}

func MongoCli() *mongo.Client {
	return mongoCli
}
