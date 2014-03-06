// blacklistConfig
package configuration

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"log"
	"os"
	"regexp"
	"strings"
)

type BlacklistConfig struct {
	blacklist map[string]interface{}
}

func newBlacklistConfig(filePath string) *BlacklistConfig {
	//Read config
	conf, err := config.NewConfig("json", filePath)
	if err != nil {
		log.Println(err)
		conf = config.NewFakeConfig()
	}
	tmp_blacklist, err := conf.DIY("blacklist")
	if tmp_blacklist == nil {
		tmp_blacklist = make(map[string]interface{})
	}
	blacklist, _ok := tmp_blacklist.(map[string]interface{})
	if !_ok {
		log.Fatal("Failed to load config.")
	}
	fmt.Println(blacklist)
	bc := &BlacklistConfig{
		blacklist: blacklist,
	}
	return bc
}

func (bc *BlacklistConfig) Contains(ip string) bool {
	_, _ok := bc.blacklist[ip]
	return _ok
}

func (bc *BlacklistConfig) Add(ip string) {
	if _, _ok := bc.blacklist[ip]; !_ok {
		bc.blacklist[ip] = 1
	} else {
		log.Printf("Already blocked %s.\n", ip)
	}
}

func (bc *BlacklistConfig) SaveToFile(filePath string) {
	configTemplate := `{"blacklist": {${blacklist}}}`
	reg := regexp.MustCompile(`([\d\.]+)(:\d)`)
	template := `"$1"$2`

	//Save config
	//fileBakPath := filePath + fmt.Sprint(time.Now()) + "_bak"
	fileBakPath := filePath + "_bak"
	os.Remove(fileBakPath)
	err := os.Rename(filePath, fileBakPath)
	if err != nil {
		log.Println(err)
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
	}
	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsoncontext := reg.ReplaceAllString(fmt.Sprint(bc.blacklist), template)
	jsoncontext = strings.NewReplacer("map[", "", "]", "", " ", ",").Replace(jsoncontext)
	jsoncontext = strings.NewReplacer("${blacklist}", jsoncontext).Replace(configTemplate)
	_, err = f.WriteString(jsoncontext)
	if err != nil {
		log.Fatal(err)
	}
}
