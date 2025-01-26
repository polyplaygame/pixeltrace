package repo

import (
	"context"
	"encoding/json"
	"errors"
	"pixeltrace/biz/dal/cache"
	"pixeltrace/biz/dal/model"
	"pixeltrace/biz/dal/query"

	"gorm.io/gorm"
)

// SaveAppCode 保存应用代码
func SaveAppCode(ctx context.Context, appCode *model.AppCode) error {
	// 校验应用代码
	if err := appCode.Validate(); err != nil {
		return err
	}
	// 先从缓存中删除
	cache.Cache.Delete(cache.GetAppCodeKey(appCode.Code))
	// 保存到数据库
	return query.AppCode.WithContext(ctx).Save(appCode)
}

// GetAppCodes 获取所有应用代码
func GetAppCodes(ctx context.Context) ([]*model.AppCode, error) {
	// 使用 query.AppCode 获取查询构建器
	q := query.AppCode

	// 执行查询并获取所有记录
	appCodes, err := q.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}
	return appCodes, nil
}

// CheckAppCode 校验应用代码
func CheckAppCode(ctx context.Context, code string) (bool, error) {
	c, err := GetAppCode(ctx, code)
	if err != nil {
		return false, err
	}
	if c == nil {
		return false, nil
	}
	return true, nil
}

// GetAppCode 获取应用代码
func GetAppCode(ctx context.Context, code string) (*model.AppCode, error) {
	appCode, err := GetAppCodeNilErr(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return appCode, nil
}

// GetAppCodeNilErr 获取应用代码
func GetAppCodeNilErr(ctx context.Context, code string) (*model.AppCode, error) {
	// 先从缓存中查询
	if data, err := cache.Cache.Get(cache.GetAppCodeKey(code)); err == nil {
		var appCode model.AppCode
		if err := json.Unmarshal([]byte(data), &appCode); err == nil {
			return &appCode, nil
		}
	}

	// 缓存未命中,从数据库查询
	q := query.AppCode
	appCode, err := q.WithContext(ctx).Where(q.Code.Eq(code)).First()
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if data, err := json.Marshal(appCode); err == nil {
		cache.Cache.Set(cache.GetAppCodeKey(code), data)
	}

	return appCode, nil
}

// UpdateAppCode 更新应用代码
func UpdateAppCode(ctx context.Context, appCode *model.AppCode) error {
	// 校验应用代码
	if err := appCode.Validate(); err != nil {
		return err
	}
	// 更新数据库
	err := query.AppCode.WithContext(ctx).Save(appCode)
	if err != nil {
		return err
	}

	// 更新缓存
	data, err := json.Marshal(appCode)
	if err != nil {
		return err
	}
	return cache.Cache.Set(cache.GetAppCodeKey(appCode.Code), data)
}

// SearchAppCode 搜索应用代码
func SearchAppCode(ctx context.Context, page int64, pageSize int64) ([]*model.AppCode, int64, error) {
	q := query.AppCode
	qtx := q.WithContext(ctx)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return qtx.FindByPage(int((page-1)*pageSize), int(pageSize))
}
