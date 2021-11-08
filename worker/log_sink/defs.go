package log_sink

import (
	"cronTab/common/cron_job"
	"cronTab/worker/mongodb_ops"
	"go.mongodb.org/mongo-driver/mongo"
)

// LogSink 日志池
type LogSink struct {
	cli            *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *cron_job.JobLog
	autoCommitChan chan *cron_job.LogBatch
}

var (
	GLogSink *LogSink
)

// InitLogSink 初始化日志池
func InitLogSink() error {
	// 选择 db 和 collection
	GLogSink = &LogSink{
		cli:            mongodb_ops.MongoCli(),
		logCollection:  mongodb_ops.MongoCli().Database("cron").Collection("log"),
		logChan:        make(chan *cron_job.JobLog, 1000),
		autoCommitChan: make(chan *cron_job.LogBatch, 1000),
	}

	return nil
}
