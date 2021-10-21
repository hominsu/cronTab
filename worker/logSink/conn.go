package logSink

import (
	"cronTab/common"
	"cronTab/worker/mongodbOps"
	"go.mongodb.org/mongo-driver/mongo"
)

// LogSink 日志池
type LogSink struct {
	cli            *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *common.JobLog
	autoCommitChan chan *common.LogBatch
}

var (
	GLogSink *LogSink
)

// InitLogSink 初始化日志池
func InitLogSink() error {
	// 选择 db 和 collection
	GLogSink = &LogSink{
		cli:            mongodbOps.MongodbCli,
		logCollection:  mongodbOps.MongodbCli.Database("cron").Collection("log"),
		logChan:        make(chan *common.JobLog, 1000),
		autoCommitChan: make(chan *common.LogBatch, 1000),
	}

	// 启动 mongodb 处理协程
	go GLogSink.writeLoop()

	return nil
}
