package service

import (
	"context"

	"pixeltrace/biz/dal/model"
	"pixeltrace/biz/dal/repo"
	"pixeltrace/conf"
	common "pixeltrace/hertz_gen/common"
	pixel "pixeltrace/hertz_gen/pixel"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type AppCodeCreateService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAppCodeCreateService(Context context.Context, RequestContext *app.RequestContext) *AppCodeCreateService {
	return &AppCodeCreateService{RequestContext: RequestContext, Context: Context}
}

func (h *AppCodeCreateService) Run(req *pixel.CreateAppCodeRequest) (resp *pixel.CreateAppCodeResponse, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	code, err := repo.GetAppCode(h.Context, req.Code)
	if err != nil {
		hlog.CtxErrorf(h.Context, "get app code error: %v", err)
		resp = &pixel.CreateAppCodeResponse{
			Status: common.Status_Error,
			Msg:    err.Error(),
		}
		return resp, nil
	}
	if code != nil {
		resp = &pixel.CreateAppCodeResponse{
			Status: common.Status_Error,
			Msg:    "code already exists",
		}
		return resp, nil
	}
	appCode := &model.AppCode{
		Code:        req.Code,
		Description: req.Description,
		TimeZone:    conf.GetConf().Pixel.TimeZone, // 时区
	}
	err = repo.SaveAppCode(h.Context, appCode)
	if err != nil {
		hlog.CtxErrorf(h.Context, "save app code error: %v", err)
		resp = &pixel.CreateAppCodeResponse{
			Status: common.Status_Error,
			Msg:    err.Error(),
		}
		return resp, nil
	}
	resp = &pixel.CreateAppCodeResponse{
		Status: common.Status_Success,
		Msg:    common.Status_Success.String(),
		Data: &pixel.AppCode{
			Code:        appCode.Code,
			Description: appCode.Description,
			TimeZone:    appCode.TimeZone,
		},
	}
	return resp, nil
}
