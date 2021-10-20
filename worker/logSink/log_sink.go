package logSink

import (
	"context"
	"cronTab/common"
	"cronTab/worker/config"
	"github.com/golang/glog"
	"time"
)

// Append 发送日志
func (logSink *LogSink) Append(jobLog *common.JobLog) {
	select {
	case logSink.logChan <- jobLog:
	default:
		// 队列满了就丢弃
	}
}

// 日志存储协程
func (logSink *LogSink) writeLoop() {
	var logBatch *common.LogBatch
	var commitTimer *time.Timer

	for {
		select {
		case log := <-logSink.logChan:
			// 初始化日志批次
			if logBatch == nil {
				logBatch = &common.LogBatch{}
				// 让这个批次超时自动提交
				commitTimer = time.AfterFunc(time.Duration(config.GConfig.JobLogCommitTimeout)*time.Millisecond,
					func(batch *common.LogBatch) func() {
						return func() {
							// 发出超时通知
							logSink.autoCommitChan <- batch
						}
					}(logBatch))
			}

			// 新的日志追加到批次
			logBatch.Logs = append(logBatch.Logs, log)

			// 如果批次满了就立即发送
			if len(logBatch.Logs) >= config.GConfig.JobLogBatchSize {
				// 发送函数
				if err := logSink.saveLogs(logBatch); err != nil {
					glog.Warning(err)
				}
				// 清空 logBatch
				logBatch = nil
				// 取消定时器
				commitTimer.Stop()
			}
		case timeoutBatch := <-logSink.autoCommitChan: // 过期的批次
			// 判断过期批次是否仍是当前批次
			if timeoutBatch != logBatch {
				continue // 跳过已经被提交的批次
			}
			// 把这个批次写入到 mongodb
			if err := logSink.saveLogs(timeoutBatch); err != nil {
				glog.Warning(err)
			}
			// 清空 logBatch
			logBatch = nil
		}
	}
}

// 批量写入日志到 mongodb
func (logSink *LogSink) saveLogs(batch *common.LogBatch) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	if _, err := logSink.logCollection.InsertMany(ctx, batch.Logs); err != nil {
		return err
	}
	return nil
}
