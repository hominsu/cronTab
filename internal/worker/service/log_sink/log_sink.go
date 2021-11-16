package log_sink

import (
	"context"
	"cronTab/configs/worker_conf"
	"cronTab/internal/pkg/cron_job"
	"cronTab/internal/pkg/xerrors"
	terrors "github.com/pkg/errors"
	"time"
)

// Append 发送日志
func (logSink *LogSink) Append(jobLog *cron_job.JobLog) {
	select {
	case logSink.logChan <- jobLog:
	default:
		// 队列满了就丢弃
	}
}

// WriteLoop 日志存储协程
func (logSink *LogSink) WriteLoop(stop <-chan struct{}) error {
	var logBatch *cron_job.LogBatch
	var commitTimer *time.Timer

	for {
		select {
		case log := <-logSink.logChan:
			// 初始化日志批次
			if logBatch == nil {
				logBatch = &cron_job.LogBatch{}
				// 让这个批次超时自动提交
				commitTimer = time.AfterFunc(time.Duration(worker_conf.GConfig.JobLogCommitTimeout)*time.Millisecond,
					func(batch *cron_job.LogBatch) func() {
						return func() {
							// 发出超时通知
							logSink.autoCommitChan <- batch
						}
					}(logBatch))
			}

			// 新的日志追加到批次
			logBatch.Logs = append(logBatch.Logs, log)

			// 如果批次满了就立即发送
			if len(logBatch.Logs) >= worker_conf.GConfig.JobLogBatchSize {
				if err := logSink.saveLogs(logBatch); err != nil {
					xerrors.ErrFmt(err)
				}
				logBatch = nil
				commitTimer.Stop()
			}

		case timeoutBatch := <-logSink.autoCommitChan: // 过期的批次
			// 判断过期批次是否仍是当前批次
			if timeoutBatch != logBatch {
				continue // 跳过已经被提交的批次
			}
			// 把这个批次写入到 mongodb
			if err := logSink.saveLogs(timeoutBatch); err != nil {
				xerrors.ErrFmt(err)
			}
			logBatch = nil

		case <-stop:
			// 平滑退出，退出前把日志写到 mongodb
			if logBatch != nil {
				if err := logSink.saveLogs(logBatch); err != nil {
					return err
				}
				logBatch = nil
			}
			commitTimer.Stop()
			return nil
		}
	}
}

// 批量写入日志到 mongodb
func (logSink *LogSink) saveLogs(batch *cron_job.LogBatch) error {
	if _, err := logSink.logCollection.InsertMany(context.Background(), batch.Logs); err != nil {
		return terrors.Wrap(err, "log sink insert job log failed")
	}
	return nil
}
