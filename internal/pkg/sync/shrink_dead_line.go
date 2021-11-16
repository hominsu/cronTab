package sync

import (
	"context"
	"time"
)

// ShrinkDeadLine 重新计算 timeout
func ShrinkDeadLine(ctx context.Context, timeout time.Duration) time.Duration {
	timeoutTime := time.Now().Add(timeout)

	if deadline, ok := ctx.Deadline(); ok && timeoutTime.After(deadline) {
		return deadline.Sub(time.Time{})
	}

	return timeoutTime.Sub(time.Time{})
}
