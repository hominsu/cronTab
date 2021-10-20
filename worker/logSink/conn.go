package logSink

import (
	"context"
	"cronTab/common"
	"cronTab/worker/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LogSink struct {
	cli            *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *common.JobLog
	autoCommitChan chan *common.LogBatch
}

var (
	GLogSink *LogSink
)

func InitLogSink() error {
	// 建立连接
	cli, err := mongo.Connect(context.TODO(),
		options.Client().
			ApplyURI(config.GConfig.MongodbUri).
			SetConnectTimeout(time.Duration(config.GConfig.MongodbConnectTimeout)*time.Millisecond))
	if err != nil {
		return err
	}

	// 选择 db 和 collection
	GLogSink = &LogSink{
		cli:            cli,
		logCollection:  cli.Database("cron").Collection("log"),
		logChan:        make(chan *common.JobLog, 1000),
		autoCommitChan: make(chan *common.LogBatch, 1000),
	}

	// 启动 mongodb 处理协程
	go GLogSink.writeLoop()

	return nil
}
