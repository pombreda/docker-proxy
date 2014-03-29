
package configuration

import (
	/*日志包*/
	"log"
)

/*创建string->空接口的map变量*/
var configMap map[string]interface{} = make(map[string]interface{})
/*constructor:构造器
构造器句柄配置*/
var configHandlerConstructorMap map[string]FileEventHandlerConstructor = make(map[string]FileEventHandlerConstructor)

/*黑白名单列表路径
将原绝对路径修改为通用相对目录*/
const BlacklistFilePath string = "./blacklist.json"
const WhitelistFilePath string = "./whitelist.json"
/**/
const BlacklistConfigKey string = "blacklistConfig"
const WhitelistConfigKey string = "whitelistConfig"
const BlacklistHandlerKey string = "blacklistHandler"
const WhitelistHandlerKey string = "whitelistHandler"

/*初始化黑名单配置*/
func initBlacklistConfig() {
	// TODO: lock configMap?
	/*日志记录*/
	log.Println("Initalizing BlacklistConfig...")
	/*打开黑名单配置文件*/
	bc := newBlacklistConfig(BlacklistFilePath)
	/*加入到map*/
	configMap[BlacklistConfigKey] = bc
}

/*获取黑名单配置*/
func GetBlacklistConfig() *BlacklistConfig {
	/*判断是否正确初始化黑名单*/
	_, ok := configMap[BlacklistConfigKey]
	/*未初始化*/
	if !ok {
		initBlacklistConfig()
	}
	/*返回黑名单配置*/
	return configMap[BlacklistConfigKey].(*BlacklistConfig)
}

/*保存黑名单表配置*/
func SaveBlacklistConfig(bc *BlacklistConfig) {
	bc.SaveToFile(BlacklistFilePath)
}

/*初始化白名单配置，同黑名单*/
func initWhitelistConfig() {
	// TODO: lock configMap?
	log.Println("Initalizing WhitelistConfig...")
	bc := newWhitelistConfig(WhitelistFilePath)
	configMap[WhitelistConfigKey] = bc
}

func GetWhitelistConfig() *WhitelistConfig {
	_, ok := configMap[WhitelistConfigKey]
	if !ok {
		initWhitelistConfig()
	}
	return configMap[WhitelistConfigKey].(*WhitelistConfig)
}

/*获得黑名单列表文件句柄*/
func GetBlacklistFileHandler() FileEventHandler {
	return configHandlerConstructorMap[BlacklistHandlerKey].new(BlacklistFilePath)
}

func GetWhitelistFileHandler() FileEventHandler {
	return configHandlerConstructorMap[WhitelistHandlerKey].new(WhitelistFilePath)
}

/*初始化*/
func init() {
	/*初始化黑白配置*/
	initBlacklistConfig()
	initWhitelistConfig()
	configHandlerConstructorMap[BlacklistHandlerKey] = BlacklistFileHandlerConstructor{}
	configHandlerConstructorMap[WhitelistHandlerKey] = WhitelistFileHandlerConstructor{}
}
