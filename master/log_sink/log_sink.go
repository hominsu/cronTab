package log_sink

import (
	"context"
	"cronTab/common"
	"cronTab/common/cron_job"
	terrors "github.com/pkg/errors"
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
		return nil, terrors.Wrap(err, "log sink find job collection failed")
	}
	defer func(cursor *mongo.Cursor) {
		if err := cursor.Close(context.TODO()); err != nil {
			common.ErrFmt(err)
		}
	}(cursor)

	// 遍历结果
	var jobLogs []*cron_job.JobLog
	for cursor.Next(ctx) {
		log := &cron_job.JobLog{}
		// 反序列化
		if err = cursor.Decode(log); err != nil {
			return nil, terrors.Wrap(err, "log sink decode job log failed")
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
		return 0, terrors.Wrap(err, "log sink delete job collection failed")
	}

	return deleteResult.DeletedCount, nil
}
