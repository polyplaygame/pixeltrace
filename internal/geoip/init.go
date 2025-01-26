package geoip

import (
	"context"
	"log"
	"pixeltrace/conf"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
)

var dirPath = "/usr/share/GeoIP" // 默认目录
var defaultLang = "en"
var reader Reader

func Init(ctx context.Context) {
	d := conf.GetConf().Pixel.GetIPDir
	if d != "" {
		dirPath = d
	}
	r, err := NewMaxMindReader(ctx, dirPath)
	if err != nil {
		hlog.Error("[Init] geoip reader init error", err)
		panic(err)
	}
	reader = r
	log.Println("[Init] geoip init success")
}

func Lookup(ip, lang string) (*Result, error) {
	if reader == nil {
		return nil, errors.New("geoip reader is nil")
	}
	if lang == "" {
		lang = defaultLang
	}
	return reader.Lookup(ip, lang)
}
