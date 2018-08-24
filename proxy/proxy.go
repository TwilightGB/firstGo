package proxy

import (
	"firstGo/ping"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sync"
	"sync/atomic"
	"time"
)

type Proxy struct {
	ipRegexp           *regexp.Regexp
	proxyIPTypeRegexp  *regexp.Regexp
	proxyUrlTypeRegexp *regexp.Regexp
	allIps             map[string]string
	all                map[string]bool
	online             int32
	threadPool         chan bool
	AviableIp          []string
	sync.Mutex
}

func New() *Proxy {
	p := &Proxy{
		ipRegexp:           regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+`),
		proxyIPTypeRegexp:  regexp.MustCompile(`https?://([\w]*:[\w]*@)?[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+`),
		proxyUrlTypeRegexp: regexp.MustCompile(`((https?|ftp):\/\/)?(([^:\n\r]+):([^@\n\r]+)@)?((www\.)?([^/\n\r:]+)):?([0-9]{1,5})?\/?([^?\n\r]+)?\??([^#\n\r]*)?#?([^\n\r]*)`),
		allIps:             map[string]string{},
		all:                map[string]bool{},
		threadPool:         make(chan bool, 1000),
	}
	p.Update()
	return p
}

// 更新代理IP列表
func (self *Proxy) Update() *Proxy {
	f, err := os.Open("E:\\mygo\\src\\firstGo\\proxy.lib")
	if err != nil {
		// logs.Log.Error("Error: %v\n", err)
		return self
	}
	b, _ := ioutil.ReadAll(f)
	f.Close()

	proxysIPType := self.proxyIPTypeRegexp.FindAllString(string(b), -1)
	for _, proxy := range proxysIPType {
		self.allIps[proxy] = self.ipRegexp.FindString(proxy)
		self.all[proxy] = false
	}

	proxysUrlType := self.proxyUrlTypeRegexp.FindAllString(string(b), -1)
	for _, proxy := range proxysUrlType {
		gvalue := self.proxyUrlTypeRegexp.FindStringSubmatch(proxy)
		self.allIps[proxy] = gvalue[6]
		self.all[proxy] = false
	}

	log.Printf(" * 读取代理IP: %v 条\n", len(self.all))

	self.findOnline()

	return self

}

// 筛选在线的代理IP
func (self *Proxy) findOnline() *Proxy {
	log.Printf(" * 正在筛选在线的代理IP……")
	self.online = 0
	for proxy := range self.all {
		self.threadPool <- true
		go func(proxy string) {
			alive, _, _ := ping.Ping(self.allIps[proxy], 4)
			self.Lock()
			self.all[proxy] = alive
			self.Unlock()
			if alive {
				atomic.AddInt32(&self.online, 1)
			}
			<-self.threadPool
		}(proxy)
	}
	for len(self.threadPool) > 0 {
		time.Sleep(0.2e9)
	}
	self.online = atomic.LoadInt32(&self.online)
	log.Printf(" * 在线代理IP筛选完成，共计：%v 个\n", self.online)
	for key, value := range self.all {
		if value {
			self.AviableIp = append(self.AviableIp, key)
		}
	}
	return self
}
