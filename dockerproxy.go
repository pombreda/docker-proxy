package main

import (
	configuration "github.com/alan7yg/docker-proxy/configuration"
	"github.com/elazarl/goproxy"
	"github.com/fzzy/radix/redis"
	"log"
	"net/http"
	"regexp"
	"time"
)

const BLOCK_TRESHOLD int = 10

var ipListCacheClient *redis.Client

func main() {
	proxy := goproxy.NewProxyHttpServer()

	initConfig()

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			match, _ := regexp.MatchString("^*.docker.io$", r.URL.Host)
			ip := getIP(r)
			if configuration.GetWhitelistConfig().Contains(ip) {
				log.Printf("IP in whitelist: %s\n", ip)
				return r, nil
			} else if configuration.GetBlacklistConfig().Contains(ip) {
				log.Println("The remote ip is in blacklist.")
				return r, goproxy.NewResponse(r, goproxy.ContentTypeText, http.StatusForbidden,
					"The remote ip is in blacklist.")
			} else if match == false {
				blockedTimes, err := ipListCacheClient.Cmd("incr", ip).Int()
				if err != nil {
					log.Println(err)
				}
				log.Printf("Request blocked. IP: %s, Block times:%d", ip, blockedTimes)
				if blockedTimes >= BLOCK_TRESHOLD {
					configuration.GetBlacklistConfig().Add(ip)
					configuration.SaveBlacklistConfig(configuration.GetBlacklistConfig())
				}
				return r, goproxy.NewResponse(r, goproxy.ContentTypeText, http.StatusForbidden,
					"The proxy is used exclusively to download docker image, please don't abuse it for any purpose.")
			} else {
				return r, nil
			}
		})
	proxy.Verbose = false
	http.ListenAndServe(":8384", proxy)
}

func getIP(r *http.Request) string {
	return regexp.MustCompile(`:\d+`).ReplaceAllString(r.RemoteAddr, "")
}

func initConfig() {
	filelist, handlerlist := make([]string, 0, 10), make([]configuration.FileEventHandler, 0, 10)
	filelist = append(filelist, configuration.BlacklistFilePath)
	filelist = append(filelist, configuration.WhitelistFilePath)

	handlerlist = append(handlerlist, configuration.GetBlacklistFileHandler())
	handlerlist = append(handlerlist, configuration.GetWhitelistFileHandler())

	cw := configuration.NewConfigFileWatcher(filelist, handlerlist)
	go cw.StartWatch()

	tmpClient, err := redis.DialTimeout("tcp", "42.96.142.222:6379", time.Duration(100)*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	ipListCacheClient = tmpClient
}
