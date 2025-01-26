package geoip

import (
	"context"
	"log"

	"github.com/fsnotify/fsnotify"
)

// watchForChanges 监听文件变化
func watchForChanges(ctx context.Context, filePath string, reloadFunc func(string) error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("watchForChanges panic: %v", r)
		}
	}()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("创建文件监听器失败: %v", err)
		return
	}
	defer watcher.Close()

	done := make(chan bool)
	defer close(done)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("watchForChanges goroutine panic: %v", r)
			}
			done <- true
		}()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Println("文件监听器已关闭")
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if err := reloadFunc(filePath); err != nil {
						log.Printf("重载数据库文件失败: %v", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Println("文件监听器已关闭")
					return
				}
				log.Printf("文件监听错误: %v", err)
			case <-ctx.Done():
				log.Println("上下文已取消")
				return
			}
		}
	}()

	if err := watcher.Add(filePath); err != nil {
		log.Printf("添加文件到监听器失败: %v", err)
		return
	}
	<-done
}
