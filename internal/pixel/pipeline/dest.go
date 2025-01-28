package pipeline

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	jsoniter "github.com/json-iterator/go"
)

const (
	originFileDestDir     = "./raw"
	originFileDestName    = "raw"
	originFileDestExt     = ".log"
	originFileDestMaxSize = 1024 * 1024 * 100 // 100MB
)

// 换行符
var newLine = []byte{'\n'}
var newLineLen = len(newLine)

// FileDest 文件写入器
type FileDest struct {
	fileDir          string
	fileName         string
	ext              string
	maxSize          int64
	byteCount        *atomic.Int64
	withBufferWriter bool
	bufSize          int
	writer           io.Writer
	file             *os.File
	lastRotated      time.Time // 用于处理跨天逻辑
	ctx              context.Context
	cancel           context.CancelFunc
	inlet            chan any
	closed           atomic.Bool
	loc              *time.Location
	mu               sync.Mutex
}

type Option func(*FileDest)

func WithRootDir(d string) Option {
	return func(f *FileDest) {
		f.fileDir = d
	}
}

func WithFileName(n string) Option {
	return func(f *FileDest) {
		f.fileName = n
	}
}

func WithFileExt(e string) Option {
	return func(f *FileDest) {
		f.ext = e
	}
}

func WithMaxSize(s int64) Option {
	return func(f *FileDest) {
		f.maxSize = s
	}
}

func WithBufferWriter(size int) Option {
	return func(f *FileDest) {
		f.withBufferWriter = true
		f.bufSize = size
	}
}

func WithCtx(c context.Context) Option {
	return func(f *FileDest) {
		f.ctx, f.cancel = context.WithCancel(c)
	}
}

func WithLoc(l *time.Location) Option {
	return func(f *FileDest) {
		f.loc = l
	}
}

func NewFileDest(opts ...Option) (*FileDest, error) {
	h := &FileDest{
		inlet:     make(chan any, 1), // 初始化
		byteCount: new(atomic.Int64),
		closed:    atomic.Bool{},
		fileDir:   originFileDestDir,
		fileName:  originFileDestName,
		ext:       originFileDestExt,
		maxSize:   originFileDestMaxSize,
		loc:       time.Local,
	}
	for _, opt := range opts {
		opt(h)
	}
	if h.ctx == nil {
		h.ctx, h.cancel = context.WithCancel(context.Background())
	}
	h.closed.Store(false)
	h.byteCount.Store(0)
	h.init()
	return h, nil
}

func (n *FileDest) In() chan<- any {
	return n.inlet
}

func (n *FileDest) rotate() {
	defer func() {
		if err := recover(); err != nil {
			hlog.Errorf("pixel event pipeline:rotate file panic, err: %v", err)
			// 如果panic了，就不要再rotate了
			return
		}
		hlog.Infof("pixel event pipeline:rotate file success, fileName: %v", n.file.Name())
	}()

	n.mu.Lock()
	defer n.mu.Unlock()

	if !strings.HasSuffix(n.fileDir, "/") {
		n.fileDir += "/"
	}

	// !!! 注意：这里的时间是按照时区生成的
	day := time.Now().In(n.loc).Format(time.DateOnly) // 按照时区生成
	dayDir := fmt.Sprintf("dt=%s", day)               // 按天目录生成
	dir := filepath.Join(n.fileDir, dayDir)           //  第一层：raw，第二层:按天设置目录

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			hlog.Errorf("Failed to create directory %s: %v", dir, err)
			return
		}
	}

	// !! 注意：这里的时间是按照时区生成的
	now := time.Now().In(n.loc).Format("2006_01_02_15_04_05")
	filePath := filepath.Join(dir, n.fileName+"_"+now+n.ext)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		hlog.Errorf("Failed to open file %s: %v", filePath, err)
		return
	}

	if n.file != nil {
		n.closeFile()
	}

	n.file = file
	if n.withBufferWriter {
		n.writer = bufio.NewWriterSize(file, n.bufSize)
	} else {
		n.writer = file
	}
	// !! 注意：这里的时间是按照时区生成的
	n.lastRotated = time.Now().In(n.loc)
}

func (n *FileDest) closeFile() {
	defer func() {
		if err := recover(); err != nil {
			return
		}
		hlog.Infof("pixel event pipeline:close file: %v", n.file.Name())
	}()
	if n.file == nil {
		return
	}

	if n.withBufferWriter {
		w := n.writer.(*bufio.Writer)
		_ = w.Flush()
	}

	_ = n.file.Close()
}

func (n *FileDest) close() {
	if !n.closed.Load() {
		n.closed.Store(true)
		close(n.inlet)
		n.cancel()
		n.closeFile()
	}
}

func (n *FileDest) Close() {
	n.close()
}

func (n *FileDest) init() {
	go func() {

		defer func() {
			if err := recover(); err != nil {
				hlog.Errorf("pixel event pipeline:file dest init panic, err: %v", err)
				return
			}
			n.close()
		}()
		// rotate first file
		n.rotate()

		// write data
		for {
			select {
			case data, ok := <-n.inlet:
				if !ok {
					return
				}
				b, err := jsoniter.Marshal(data)
				if err != nil {
					hlog.Errorf("pixel event pipeline:marshal event error, data: %v, err: %v", data, err)
					continue
				}
				bl, err := n.writer.Write(b)
				if err != nil {
					hlog.Errorf("pixel event pipeline:write event error: %v", err)
					continue
				}
				_, err = n.writer.Write(newLine)
				if err != nil {
					hlog.Errorf("pixel event pipeline:write new line error: %v", err)
					continue
				}
				n.byteCount.Add(int64(bl + newLineLen))
				// 日志分割逻辑：
				// 跨天分割 && 大小分割
				// !!! 注意：这里的时间是按照时区生成的
				if n.byteCount.Load() >= n.maxSize {
					hlog.Infof("pixel event pipeline:rotate file by size, fileName: %v", n.file.Name())
					n.rotate()
					n.byteCount.Store(0)
				} else {
					// 避免竞态条件，获取当前时间的副本进行比较
					now := time.Now().In(n.loc)
					if now.Day() != n.lastRotated.Day() {
						hlog.Infof("pixel event pipeline:rotate file by day, fileName: %v", n.file.Name())
						n.rotate()
						n.byteCount.Store(0)
					}
				}
			case <-n.ctx.Done():
				return
			}
		}
	}()
}
