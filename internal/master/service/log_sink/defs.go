package log_sink

import (
	"cronTab/internal/master/data/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

// LogSink 日志池
type LogSink struct {
	cli           *mongo.Client
	logCollection *mongo.Collection
}

var (
	GLogSink *LogSink
)

// InitLogSink 初始化日志池
func InitLogSink() error {
	// 选择 db 和 collection
	GLogSink = &LogSink{
		cli:           mongodb.Cli,
		logCollection: mongodb.Cli.Database("cron").Collection("log"),
	}

	// 启动 mongodb 处理协程
	//go GLogSink.writeLoop()

	return nil
}
