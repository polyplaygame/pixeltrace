package async

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func CtxGo(ctx context.Context, fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				hlog.CtxErrorf(ctx, "[async]panic: %v", err)
			}
		}()
		fn()
	}()
}

func Go(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				hlog.Errorf("[async]panic: %v", err)
			}
		}()
		fn()
	}()
}
