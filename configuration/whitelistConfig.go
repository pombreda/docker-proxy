// whitelistConfig
package configuration

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"log"
)

/*白名单配置结构体*/
type WhitelistConfig struct {
	whitelist map[string]interface{}
}

/*新建白名单配置，完全等同于黑名单，详情见blacklistConfig.go*/
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

/*判断是否已存在*/
func (bc *WhitelistConfig) Contains(ip string) bool {
	_, _ok := bc.whitelist[ip]
	return _ok
}
