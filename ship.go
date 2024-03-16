package main

import (
    "flag"
    "log"

    "github.com/raymondragon/golib"
)

var rawURL = flag.String("url", "", "http(s)://user:pass@host:port/path#dir")

func main() {
    flag.Parse()
    if *rawURL == nil {
        flag.Usage()
        log.Fatalf("[ERRO] %v", "Flag Missing")
    }
    parsedURL, err := golib.URLParse(*rawURL)
    if err != nil {
        log.Printf("[WARN] %v", err)
    }
    webdavHandler := golib.WebdavHandler(parsedURL.Fragment, parsedURL.Path)
    proxyHandler := golib.ProxyHandler(parsedURL.Hostname, parsedURL.Username, parsedURL.Password, webdavHandler)
    if parsedURL.Path == "" || parsedURL.Path == "/" {
        proxyHandler = golib.ProxyHandler(parsedURL.Hostname, parsedURL.Username, parsedURL.Password, nil)
    }
    switch parsedURL.Scheme {
    case "http":
        log.Printf("[INFO] %v", *rawURL)
        if err := golib.ServeHTTP(parsedURL.Hostname, parsedURL.Port, proxyHandler); err != nil {
            log.Fatalf("[ERRO] %v", err)
        }
    case "https":
        tlsConfig, err := golib.TLSConfigApplication(parsedURL.Hostname)
        if err != nil {
            log.Printf("[WARN] %v", err)
            tlsConfig, err = golib.TLSConfigGeneration(parsedURL.Hostname)
            if err != nil {
                log.Printf("[WARN] %v", err)
            }
        }
        log.Printf("[INFO] %v", *rawURL)
        if err := golib.ServeHTTPS(parsedURL.Hostname, parsedURL.Port, proxyHandler, tlsConfig); err != nil {
            log.Fatalf("[ERRO] %v", err)
        }
    default:
        log.Fatalf("[ERRO] %v", parsedURL.Scheme)
    }
}
