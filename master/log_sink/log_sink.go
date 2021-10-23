package log_sink

import (
	"context"
	"cronTab/common/cron_job"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// GetLogBatch 获取 log
func (logSink *LogSink) GetLogBatch(name string) ([]*cron_job.JobLog, error) {
	var err error

	ctx, cancel := context.WithTimeout(context.TODO(), 500*time.Millisecond)
	defer cancel()

	// 分页
	//findOptions := &options.FindOptions{}
	//findOptions.SetSkip(0)
	//findOptions.SetLimit(5)

	// 查找
	cursor, err := logSink.logCollection.Find(ctx, bson.M{"job_name": name})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor) {
		if err := cursor.Close(context.TODO()); err != nil {
			glog.Warning(err)
		}
	}(cursor)

	// 遍历结果
	var jobLogs []*cron_job.JobLog
	for cursor.Next(ctx) {
		log := &cron_job.JobLog{}
		// 反序列化
		if err = cursor.Decode(log); err != nil {
			return nil, err
		}
		jobLogs = append(jobLogs, log)
	}

	return jobLogs, nil
}

// DelJobLog 删除任务日志
func (logSink *LogSink) DelJobLog(name string) (int64, error) {
	var err error

	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()

	// 删除
	deleteResult, err := logSink.logCollection.DeleteMany(ctx, bson.M{"job_name": name})
	if err != nil {
		return 0, err
	}

	return deleteResult.DeletedCount, nil
}
