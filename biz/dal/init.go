package dal

import (
	"context"
	"log"
	"pixeltrace/biz/dal/cache"
	"pixeltrace/biz/dal/mysql"
	"pixeltrace/biz/dal/query"
)

func Init(ctx context.Context) {
	// redis.Init()
	cache.Init(ctx)
	mysql.Init()
	query.SetDefault(mysql.DB)
	log.Println("init dal success")
}
