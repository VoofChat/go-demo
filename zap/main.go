package main

import (
	"go.uber.org/zap"
	"net/http"
)


func main() {
	InitZap1()
	InitZap2()
	defer sugarLogger1.Sync()
	defer sugarLogger2.Sync()
	for i := 0; i < 1; i++ {
		simpleHttpGet(sugarLogger1, "www.google.com")
		simpleHttpGet(sugarLogger1, "http://www.baidu.com")

		simpleHttpGet(sugarLogger2, "www.google.com")
		simpleHttpGet(sugarLogger2, "http://www.baidu.com")
	}
}

func simpleHttpGet(sugarLogger *zap.SugaredLogger, url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching url %s.. Error : %s", url, err)
		sugarLogger.Errorf("Error fetching url %s.. Error : %s", url, err)
	} else {
		sugarLogger.Infof("Success..! statusCode = %s for url %s", resp.Status, url)
		sugarLogger.Infof("Success..! statusCode = %s for url %s", resp.Status, url)
		_ = resp.Body.Close()
	}
}

