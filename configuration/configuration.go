package configuration

import (
	"log"
)

var configMap map[string]interface{} = make(map[string]interface{})
var configHandlerConstructorMap map[string]FileEventHandlerConstructor = make(map[string]FileEventHandlerConstructor)

const BlacklistFilePath string = "E:\\GoWorkspace\\conf\\blacklist.json"
const WhitelistFilePath string = "E:\\GoWorkspace\\conf\\whitelist.json"
const BlacklistConfigKey string = "blacklistConfig"
const WhitelistConfigKey string = "whitelistConfig"
const BlacklistHandlerKey string = "blacklistHandler"
const WhitelistHandlerKey string = "whitelistHandler"

func initBlacklistConfig() {
	// TODO: lock configMap?
	log.Println("Initalizing BlacklistConfig...")
	bc := newBlacklistConfig(BlacklistFilePath)
	configMap[BlacklistConfigKey] = bc
}

func GetBlacklistConfig() *BlacklistConfig {
	_, ok := configMap[BlacklistConfigKey]
	if !ok {
		initBlacklistConfig()
	}
	return configMap[BlacklistConfigKey].(*BlacklistConfig)
}

func SaveBlacklistConfig(bc *BlacklistConfig) {
	bc.SaveToFile(BlacklistFilePath)
}

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

func GetBlacklistFileHandler() FileEventHandler {
	return configHandlerConstructorMap[BlacklistHandlerKey].new(BlacklistFilePath)
}

func GetWhitelistFileHandler() FileEventHandler {
	return configHandlerConstructorMap[WhitelistHandlerKey].new(WhitelistFilePath)
}

func init() {
	initBlacklistConfig()
	initWhitelistConfig()
	configHandlerConstructorMap[BlacklistHandlerKey] = BlacklistFileHandlerConstructor{}
	configHandlerConstructorMap[WhitelistHandlerKey] = WhitelistFileHandlerConstructor{}
}
