package pixel

import (
	"context"
	"log"
	"pixeltrace/internal/pixel/pipeline"
)

var s *pipeline.Stream

func Init(ctx context.Context) {
	s = pipeline.NewStream(ctx)
	log.Println("init pixel success")
}
