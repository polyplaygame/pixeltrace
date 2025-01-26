package pixel

import (
	"context"

	"pixeltrace/biz/service"
	"pixeltrace/biz/utils"
	common "pixeltrace/hertz_gen/common"
	pixel "pixeltrace/hertz_gen/pixel"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const (
	blankImg = "image/png"
)

var blankImgData = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

// HandleImgPixelEvent .
// @router /pixel [GET]
func HandleImgPixelEvent(ctx context.Context, c *app.RequestContext) {
	var err error
	var req pixel.ImgPixelEvent
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	// 获取请求头中的ip
	ip := c.ClientIP()
	req.Ip = ip
	_, err = service.NewHandleImgPixelEventService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	// 返回一个空白图片，43byte
	c.Data(consts.StatusOK, blankImg, blankImgData)
}

// HandleBeaconPixelEvent .
// @router /pixel [POST]
func HandleBeaconPixelEvent(ctx context.Context, c *app.RequestContext) {
	var err error
	var req pixel.BeaconPixelEvent
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	// 获取请求头中的ip
	ip := c.ClientIP()
	req.Ip = ip

	resp := &common.BaseResponse{}
	resp, err = service.NewHandleBeaconPixelEventService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// AppCodeDetail .
// @router /app_code/{code} [GET]
func AppCodeDetail(ctx context.Context, c *app.RequestContext) {
	var err error
	var req pixel.GetAppCodeDetailRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &pixel.GetAppCodeDetailResponse{}
	resp, err = service.NewAppCodeDetailService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// AppCodeCreate .
// @router /app_code [POST]
func AppCodeCreate(ctx context.Context, c *app.RequestContext) {
	var err error
	var req pixel.CreateAppCodeRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &pixel.CreateAppCodeResponse{}
	resp, err = service.NewAppCodeCreateService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// AppCodeSearch .
// @router /app_code/search [GET]
func AppCodeSearch(ctx context.Context, c *app.RequestContext) {
	var err error
	var req pixel.SearchAppCodeRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &pixel.SearchAppCodeResponse{}
	resp, err = service.NewAppCodeSearchService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
