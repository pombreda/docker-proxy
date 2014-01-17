package main

import (
  "github.com/elazarl/goproxy"
  "net/http"
  "regexp"
)

func main() {
  proxy := goproxy.NewProxyHttpServer()
  proxy.OnRequest().DoFunc(
    func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
      match, _ := regexp.MatchString("^*.docker.io$", r.URL.Host)
      if match == false {
        return r, goproxy.NewResponse(r, goproxy.ContentTypeText, http.StatusForbidden, "Don't use it for other proxy!")
      } else {
        return r, nil
      }
    })
  proxy.Verbose = false
  http.ListenAndServe(":8384", proxy)
}
