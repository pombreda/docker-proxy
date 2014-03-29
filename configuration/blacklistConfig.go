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

/*黑名单配置结构体*/
type BlacklistConfig struct {
	blacklist map[string]interface{}
}

/*由文件路径参数新建黑名单配置，返回黑名单配置指针*/
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

/*判断名单中是否包含此ip*/
func (bc *BlacklistConfig) Contains(ip string) bool {
	_, _ok := bc.blacklist[ip]
	return _ok
}

/*添加ip到名单*/
func (bc *BlacklistConfig) Add(ip string) {
	/*判断是否已存在，未存在的话便添加*/
	if _, _ok := bc.blacklist[ip]; !_ok {
		bc.blacklist[ip] = 1
	} else {
		/*已存在，信息写入日志*/
		log.Printf("Already blocked %s.\n", ip)
	}
}

/*保存到文件*/
func (bc *BlacklistConfig) SaveToFile(filePath string) {
	configTemplate := `{"blacklist": {${blacklist}}}`
	reg := regexp.MustCompile(`([\d\.]+)(:\d)`)
	template := `"$1"$2`

	//Save config
	//fileBakPath := filePath + fmt.Sprint(time.Now()) + "_bak"
	fileBakPath := filePath + "_bak"
	/*移除已有文件*/
	os.Remove(fileBakPath)
	/*重命名*/
	err := os.Rename(filePath, fileBakPath)
	if err != nil {
		/*重命名失败*/
		log.Println(err)
		/*失败时检测文件是否存在*/
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
	}
	/*创建文件*/
	f, err := os.Create(filePath)
	/*延迟关闭，清理工作*/
	defer f.Close()
	/*创建失败*/
	if err != nil {
		log.Fatal(err)
	}
	jsoncontext := reg.ReplaceAllString(fmt.Sprint(bc.blacklist), template)
	jsoncontext = strings.NewReplacer("map[", "", "]", "", " ", ",").Replace(jsoncontext)
	jsoncontext = strings.NewReplacer("${blacklist}", jsoncontext).Replace(configTemplate)
	/*将**（看不懂的内容）信息写入文件*/
	_, err = f.WriteString(jsoncontext)
	if err != nil {
		log.Fatal(err)
	}
}
