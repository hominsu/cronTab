package mongodb

import (
	"context"
	"cronTab/configs/worker_conf"
	"cronTab/internal/pkg/sync"
	terrors "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	Cli *mongo.Client
)

// InitMongodbConn 初始化 mongodb 连接
func InitMongodbConn(ctx context.Context) error {
	var err error

	// 建立连接
	Cli, err = mongo.Connect(ctx,
		options.Client().
			ApplyURI(worker_conf.GConfig.MongodbUri).
			SetConnectTimeout(sync.ShrinkDeadLine(ctx, time.Duration(worker_conf.GConfig.MongodbConnectTimeout)*time.Millisecond)))
	if err != nil {
		return terrors.Wrap(err, "create mongodb connection failed")
	}

	// 测试 mongodb 连接
	if err = Cli.Ping(ctx, nil); err != nil {
		return terrors.Wrap(err, "test mongodb connection failed")
	}

	return nil
}

// CloseMongodbConn 关闭 mongodb 连接
func CloseMongodbConn() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := Cli.Disconnect(ctx); err != nil {
		return terrors.Wrap(err, "disconnect mongodb failed")
	}
	return nil
}

func MongoCli() *mongo.Client {
	return Cli
}
