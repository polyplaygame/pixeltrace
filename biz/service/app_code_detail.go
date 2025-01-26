package service

import (
	"context"

	"pixeltrace/biz/dal/repo"
	common "pixeltrace/hertz_gen/common"
	pixel "pixeltrace/hertz_gen/pixel"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type AppCodeDetailService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAppCodeDetailService(Context context.Context, RequestContext *app.RequestContext) *AppCodeDetailService {
	return &AppCodeDetailService{RequestContext: RequestContext, Context: Context}
}

func (h *AppCodeDetailService) Run(req *pixel.GetAppCodeDetailRequest) (resp *pixel.GetAppCodeDetailResponse, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	appCode, err := repo.GetAppCode(h.Context, req.Code)
	if err != nil {
		hlog.CtxErrorf(h.Context, "get app code error: %v", err)
		resp = &pixel.GetAppCodeDetailResponse{
			Status: common.Status_Error,
			Msg:    err.Error(),
		}
		return
	}
	if appCode == nil {
		resp = &pixel.GetAppCodeDetailResponse{
			Status: common.Status_NotFound,
			Msg:    common.Status_NotFound.String(),
		}
		return
	}
	resp = &pixel.GetAppCodeDetailResponse{
		Status: common.Status_Success,
		Msg:    common.Status_Success.String(),
		Data: &pixel.AppCode{
			Code:        appCode.Code,
			Description: appCode.Description,
			TimeZone:    appCode.TimeZone,
		},
	}
	return
}
