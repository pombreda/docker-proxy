// configFileWatcher.go
package configuration

import (
	"github.com/howeyc/fsnotify"
	"log"
)

/*新建文件监听*/
func NewConfigFileWatcher(filePaths []string, handlers []FileEventHandler) *ConfigWatcher {
	return &ConfigWatcher{
		filePaths: filePaths,
		handlers:  handlers,
	}
}

/*监听配置结构体*/
type ConfigWatcher struct {
	filePaths []string
	handlers  []FileEventHandler
}

/*启动监听*/
func (cw *ConfigWatcher) StartWatch() {
	/*创建新监听*/
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	/*创建通讯channel*/
	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event: ", ev)
				//Add Handler
				for _, h := range cw.handlers {
					h.Handle(ev)
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	for _, filePath := range cw.filePaths {
		err = watcher.Watch(filePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	<-done

	/* ... do stuff ... */
	defer watcher.Close()
}
