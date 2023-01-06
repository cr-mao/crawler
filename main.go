package main

import (
	"fmt"
	"github.com/cr-mao/crawler/parse/doubangroup"
	"go.uber.org/zap"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/cr-mao/crawler/collect"
	"github.com/cr-mao/crawler/log"
)

func main() {

	// log
	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)
	logger.Info("log init end")
	//proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8888"}
	//p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	//if err != nil {
	//	fmt.Println("RoundRobinProxySwitcher failed")
	//}

	cookie := "bid=4P0TsgefDh8; gr_user_id=cd3d1a1d-cc88-4b35-af91-d6c065435978; __utmz=30149280.1668244404.2.2.utmcsr=amistyrain.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __gads=ID=449d581df79612af-22a024be50d80011:T=1668244403:RT=1668244403:S=ALNI_MYfPqNoEIbJY7n7PvOEb5VSeOj3Yw; viewed=\"1007305_25869486\"; __gpi=UID=00000b7a40e5f0bd:T=1668244403:RT=1672817979:S=ALNI_MYkYPaRpQVcRDw9ANWUL7HcJy7asQ; _pk_ses.100001.8cb4=*; ap_v=0,6.0; __yadk_uid=UfzazX5u8dmztBoQdXLKsU15Clixq4ik; __utma=30149280.345195326.1667959970.1672817978.1672974149.4; __utmc=30149280; __utmt=1; _pk_id.100001.8cb4=32346144b0d0c430.1672974146.1.1672974279.1672974146.; __utmb=30149280.12.5.1672974280368"

	var worklist []*collect.Request
	for i := 0; i <= 0; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		worklist = append(worklist, &collect.Request{
			Url:       str,
			Cookie:    cookie,
			ParseFunc: doubangroup.ParseURL,
		})
	}

	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		//Proxy:   p,
	}

	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			body, err := f.Get(item)
			time.Sleep(1 * time.Second)
			if err != nil {
				logger.Error("read content failed",
					zap.Error(err),
				)
				continue
			}
			res := item.ParseFunc(body, item)
			for _, item1 := range res.Items {
				logger.Info("result",
					zap.String("get url:", item1.(string)))
			}

			//广度优先
			worklist = append(worklist, res.Requesrts...)
		}
	}
}
