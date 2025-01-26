package pipeline

import (
	"context"
	"errors"
	"fmt"
	"log"
	"pixeltrace/conf"
	"pixeltrace/internal/pixel/events"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/reugn/go-streams"
	ext "github.com/reugn/go-streams/extension"
	"github.com/reugn/go-streams/flow"
)

// Stream 流处理
type Stream struct {
	ctx   context.Context
	inlet chan any
}

// NewStream 创建一个流处理
func NewStream(ctx context.Context) *Stream {
	s := &Stream{
		ctx:   ctx,
		inlet: make(chan any, 1),
	}
	tz, err := time.LoadLocation(conf.GetConf().Pixel.TimeZone)
	if err != nil {
		log.Fatalf("[stream]load timezone error: %v", err)
	}
	sink, err := s.NewSink(tz)
	if err != nil {
		log.Fatalf("[stream]init sink error: %v", err)
	}
	if err := s.NewInlet(ctx, sink); err != nil {
		log.Fatalf("[stream]init inlet error: %v", err)
	}
	return s
}

// Parse 解析事件
func (s *Stream) Parse(e *events.Event) {
	s.inlet <- e
}

// NewInlet 创建inlet
func (s *Stream) NewInlet(ctx context.Context, sink streams.Sink) error {
	if sink == nil {
		return errors.New("sink cannot be nil")
	}

	source := ext.NewChanSource(s.inlet)
	f := source.
		Via(flow.NewMap(addIPGeoExtra, 1)).
		Via(flow.NewMap(unmarshalEventData, 1))
	go func() {
		defer func() {
			hlog.Infof("[stream]inlet close")
			if r := recover(); r != nil {
				hlog.Errorf("[stream]inlet panic: %v", r)
			}
		}()
		hlog.Infof("[stream]inlet start")
		defer close(s.inlet)

		go func() {
			hlog.Infof("[stream]inlet to sink")
			f.To(sink)
		}()

		select {
		case <-ctx.Done():
			hlog.Infof("[stream]inlet context done")
			return
		}
	}()

	return nil
}

// NewSink 创建一个sink，用于将数据写入文件
func (s *Stream) NewSink(loc *time.Location) (streams.Sink, error) {
	_, offset := time.Now().In(loc).Zone()
	hourOffset := float64(offset) / 3600
	dir := fmt.Sprintf("%s+%v", originFileDestDir, hourOffset)
	fileName := originFileDestName
	sc := ""
	cnf := conf.GetConf()
	if cnf != nil {
		sc = cnf.Pixel.ServerCode
	}
	if sc != "" {
		fileName = fmt.Sprintf("%s_%s", sc, fileName)
	}

	sink, err := NewFileDest(
		WithRootDir(dir),
		WithFileName(fileName),
		WithFileExt(originFileDestExt),
		WithMaxSize(originFileDestMaxSize),
		WithLoc(loc))
	if err != nil {
		log.Fatalf("[stream]init file dest error: %v", err)
	}
	return sink, nil
}
