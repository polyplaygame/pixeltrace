package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	common "pixeltrace/hertz_gen/common"
	pixel "pixeltrace/hertz_gen/pixel"
)

type HandleImgPixelEventService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewHandleImgPixelEventService(Context context.Context, RequestContext *app.RequestContext) *HandleImgPixelEventService {
	return &HandleImgPixelEventService{RequestContext: RequestContext, Context: Context}
}

func (h *HandleImgPixelEventService) Run(req *pixel.ImgPixelEvent) (resp *common.BaseResponse, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	return
}
