package service

import (
	"context"

	"pixeltrace/biz/dal/repo"
	common "pixeltrace/hertz_gen/common"
	pixel "pixeltrace/hertz_gen/pixel"
	pixelsvc "pixeltrace/internal/pixel"
	"pixeltrace/internal/pixel/decoder"
	"pixeltrace/internal/pixel/events"
	"pixeltrace/pkg/async"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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
	// 校验appCode
	ok, err := repo.CheckAppCode(h.Context, req.App)
	if err != nil {
		resp = &common.BaseResponse{
			Status: common.Status_Error,
			Msg:    "get app code error",
		}
		return resp, nil
	}
	if !ok {
		resp = &common.BaseResponse{
			Status: common.Status_NotFound,
			Msg:    common.Status_NotFound.String(),
		}
		return resp, nil
	}
	// 解析data
	data, err := decoder.DecodeImgData(req.Data)
	if err != nil {
		hlog.CtxErrorf(h.Context, "decode data error,request:%v,err:%v", req, err)
		resp = &common.BaseResponse{
			Status: common.Status_Error,
			Msg:    "decode data error",
		}
		return resp, err
	}
	e := events.New().
		SetApp(req.App).
		SetIP(req.Ip).
		SetData(data)
	// 处理事件
	async.CtxGo(h.Context, func() {
		pixelsvc.Parse(e)
	})
	resp = &common.BaseResponse{
		Status: common.Status_Success,
		Msg:    "success",
	}
	return
}
