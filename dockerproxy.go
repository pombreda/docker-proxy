package main

import (
	/*导入自定义包:configuration*/
	"github.com/alan7yg/docker-proxy/configuration"
	"github.com/elazarl/goproxy"
	"github.com/fzzy/radix/redis"
	/*导入日志包*/
	"log"
	/*导入网络http包*/
	"net/http"
	/*导入正则表达式包*/
	"regexp"
	/*导入时间包*/
	"time"
)

/*threshold：门槛，入口，开始
最大阻塞次数常量*/
const BLOCK_TRESHOLD int = 10

/*客户端ip缓冲列表？*/
var ipListCacheClient *redis.Client

func main() {
	/*goproxy:https://github.com/elazarl/goproxy
	  An HTTP proxy library for Go http://ripper234.com/p/introducing-goproxy-light-http-proxy/
	  建立proxy server*/
	proxy := goproxy.NewProxyHttpServer()

	/*完成初始化工作*/
	initConfig()

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			/*匹配，判断URL是否为docker.io*/
			match, _ := regexp.MatchString("^*.docker.io$", r.URL.Host)
			/*得到ip字符串*/
			ip := getIP(r)
			/*处于白名单列表，写入日志*/
			if configuration.GetWhitelistConfig().Contains(ip) {
				log.Printf("IP in whitelist: %s\n", ip)
				return r, nil
			/*处于黑名单列表，写入日志*/
			} else if configuration.GetBlacklistConfig().Contains(ip) {
				/*remote:远程的，
				写入日志：远程ip在黑名单内*/
				log.Println("The remote ip is in blacklist.")
				return r, goproxy.NewResponse(r, goproxy.ContentTypeText, http.StatusForbidden,
					"The remote ip is in blacklist.")
			/*非docker URL，记录次数*/
			} else if match == false {
				blockedTimes, err := ipListCacheClient.Cmd("incr", ip).Int()
				if err != nil {
					log.Println(err)
				}
				/*block：阻塞
				写入日志-被阻塞的ip，被阻塞的次数*/
				log.Printf("Request blocked. IP: %s, Block times:%d", ip, blockedTimes)
				/*当当前阻塞次数达到最大次数时，达到条件，加入黑名单*/
				if blockedTimes >= BLOCK_TRESHOLD {
					/*当前ip加入黑名单，实时保存黑名单*/
					configuration.GetBlacklistConfig().Add(ip)
					configuration.SaveBlacklistConfig(configuration.GetBlacklistConfig())
				}
				return r, goproxy.NewResponse(r, goproxy.ContentTypeText, http.StatusForbidden,
					/*proxy是被用来单纯的用了下载docker镜像用的，请不要为了任何目的而滥用它*/
					"The proxy is used exclusively to download docker image, please don't abuse it for any purpose.")
			} else {
				return r, nil
			}
		}
		)
	proxy.Verbose = false
	/*创建tcp端口监听*/
	http.ListenAndServe(":8384", proxy)
}
 /*获得string类型的ip*/
func getIP(r *http.Request) string {
	return regexp.MustCompile(`:\d+`).ReplaceAllString(r.RemoteAddr, "")
}

func initConfig() {
	/*创建文件列表、句柄列表*/
	filelist, handlerlist := make([]string, 0, 10), make([]configuration.FileEventHandler, 0, 10)
	/*添加黑白名单文件目录到表*/
	filelist = append(filelist, configuration.BlacklistFilePath)
	filelist = append(filelist, configuration.WhitelistFilePath)
	/*添加黑白名单文件目录句柄到表*/
	handlerlist = append(handlerlist, configuration.GetBlacklistFileHandler())
	handlerlist = append(handlerlist, configuration.GetWhitelistFileHandler())

	/*属性文件监视？实时监测黑白名单变动？*/
	cw := configuration.NewConfigFileWatcher(filelist, handlerlist)
	/*开启监视进程*/
	go cw.StartWatch()

	/*创建临时客户端*/
	tmpClient, err := redis.DialTimeout("tcp", "42.96.142.222:6379", time.Duration(100)*time.Second)
	/*错误处理*/
	if err != nil {
		log.Fatal(err)
	}
	/*将正确的客户端传递赋值*/
	ipListCacheClient = tmpClient
}
