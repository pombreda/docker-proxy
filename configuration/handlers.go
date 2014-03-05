package configuration

import (
	"github.com/howeyc/fsnotify"
	"log"
)

type FileEventHandlerConstructor interface {
	new(filePath string) FileEventHandler
}

type FileEventHandler interface {
	Handle(event *fsnotify.FileEvent)
}

// func GetHandler(name, filePath string) FileEventHandler {
// 	switch strings.ToLower(name) {
// 	case "blacklisthandler":
// 		return &BlacklistHandler{
// 			filePath: strings.NewReplacer("/", "\\").Replace(filePath),
// 		}
// 	case "whitelisthandler":
// 		return &WhitelistHandler{
// 			filePath: strings.NewReplacer("/", "\\").Replace(filePath),
// 		}
// 	}
// 	return nil
// }

type BlacklistFileHandlerConstructor struct{}

func (bhc BlacklistFileHandlerConstructor) new(filePath string) FileEventHandler {
	return &BlacklistHandler{
		filePath: filePath,
	}
}

type BlacklistHandler struct {
	filePath string
}

func (h BlacklistHandler) Handle(event *fsnotify.FileEvent) {
	log.Println(event)
	log.Println(h.filePath)
	if event == nil {
		log.Print("blacklisthandler nil.")
	} else if event.Name == h.filePath && event.IsModify() {
		log.Println("blacklisthandler....")
		initBlacklistConfig()
	}
}

func (h BlacklistHandler) String() string {
	return "black\t" + h.filePath
}

type WhitelistFileHandlerConstructor struct{}

func (whc WhitelistFileHandlerConstructor) new(filePath string) FileEventHandler {
	return &WhitelistHandler{
		filePath: filePath,
	}
}

type WhitelistHandler struct {
	filePath string
}

func (h WhitelistHandler) Handle(event *fsnotify.FileEvent) {
	if event == nil {
		log.Print("whitelisthandler nil.")
	} else if event.Name == h.filePath && event.IsModify() {
		initWhitelistConfig()
	}
}

func (h WhitelistHandler) String() string {
	return "white\t" + h.filePath
}
