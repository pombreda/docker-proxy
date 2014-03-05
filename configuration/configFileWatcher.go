// configFileWatcher.go
package configuration

import (
	"github.com/howeyc/fsnotify"
	"log"
)

func NewConfigFileWatcher(filePaths []string, handlers []FileEventHandler) *ConfigWatcher {
	return &ConfigWatcher{
		filePaths: filePaths,
		handlers:  handlers,
	}
}

type ConfigWatcher struct {
	filePaths []string
	handlers  []FileEventHandler
}

func (cw *ConfigWatcher) StartWatch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

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
