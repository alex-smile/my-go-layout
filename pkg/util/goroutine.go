package util

import (
	"context"
	"log"
)

// GoWithRecovery is a wrapper of goroutine that can recover panic
func GoWithRecovery(ctx context.Context, fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
			}
		}()

		fn()
	}()
}
