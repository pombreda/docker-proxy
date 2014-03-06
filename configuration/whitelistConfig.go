// whitelistConfig
package configuration

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"log"
)

type WhitelistConfig struct {
	whitelist map[string]interface{}
}

func newWhitelistConfig(filePath string) *WhitelistConfig {
	//Read config
	conf, err := config.NewConfig("json", filePath)
	if err != nil {
		log.Println(err)
		conf = config.NewFakeConfig()
	}
	whitelistStringArray := conf.Strings("whitelist")
	whitelist := make(map[string]interface{})
	for _, i := range whitelistStringArray {
		whitelist[i] = 1
	}

	fmt.Println(whitelist)
	wc := &WhitelistConfig{
		whitelist: whitelist,
	}
	return wc
}

func (bc *WhitelistConfig) Contains(ip string) bool {
	_, _ok := bc.whitelist[ip]
	return _ok
}
