package mongodb_ops

import (
	"context"
	"cronTab/worker/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	MongodbCli *mongo.Client
)

// InitMongodbConn 初始化 mongodb 连接
func InitMongodbConn() error {
	var err error

	ctx, cancel := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	defer cancel()

	// 建立连接
	MongodbCli, err = mongo.Connect(context.TODO(),
		options.Client().
			ApplyURI(config.GConfig.MongodbUri).
			SetConnectTimeout(time.Duration(config.GConfig.MongodbConnectTimeout)*time.Millisecond))
	if err != nil {
		return err
	}

	// 测试 mongodb 连接
	if err = MongodbCli.Ping(ctx, nil); err != nil {
		return err
	}

	return nil
}

// CloseMongodbConn 关闭 mongodb 连接
func CloseMongodbConn() error {
	ctx, cancel := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	defer cancel()

	if err := MongodbCli.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}
