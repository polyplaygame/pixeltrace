package service

import (
	"context"

	"pixeltrace/biz/dal/repo"
	common "pixeltrace/hertz_gen/common"
	pixel "pixeltrace/hertz_gen/pixel"

	"github.com/cloudwego/hertz/pkg/app"
)

type AppCodeSearchService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAppCodeSearchService(Context context.Context, RequestContext *app.RequestContext) *AppCodeSearchService {
	return &AppCodeSearchService{RequestContext: RequestContext, Context: Context}
}

func (h *AppCodeSearchService) Run(req *pixel.SearchAppCodeRequest) (resp *pixel.SearchAppCodeResponse, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	appCodes, total, err := repo.SearchAppCode(h.Context, req.Page, req.PageSize)
	if err != nil {
		resp = &pixel.SearchAppCodeResponse{
			Status: common.Status_ServerError,
			Msg:    err.Error(),
		}
		return resp, nil
	}
	items := make([]*pixel.AppCode, 0)
	for _, appCode := range appCodes {
		items = append(items, &pixel.AppCode{
			Code:        appCode.Code,
			Description: appCode.Description,
			TimeZone:    appCode.TimeZone,
		})
	}
	resp = &pixel.SearchAppCodeResponse{
		Status: common.Status_Success,
		Data: &pixel.SearchAppCodeData{
			Total: total,
			Items: items,
		},
		Msg: common.Status_Success.String(),
	}
	return resp, nil
}
