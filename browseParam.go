package phantomgo

import (
	"net/http"
	"time"
)

//供内部调用
type Request interface {
	GetMethod() string
	GetUrl() string
	GetHeader() http.Header
	GetPostBody() string
	GetRedirectTimes() int
	GetDialTimeout() time.Duration
	GetConnTimeout() time.Duration
	GetRetryPause() time.Duration
	GetTryTimes() int
	GetusePhantomJS() bool
}

//供外部调用
type Param struct {
	Method        string
	Url           string
	Header        http.Header
	PostBody      string
	RedirectTimes int           //重定向次数
	DialTimeout   time.Duration //拨号超时时间段
	ConnTimeout   time.Duration //链接超时时间
	RetryPause    time.Duration //请求失败时重复试时间段
	TryTimes      int           //请求失败重新请求次数
	UsePhantomJS  bool
}

func (self *Param) GetMethod() string {
	return self.Method
}

func (self *Param) GetUrl() string {
	return self.Url
}
func (self *Param) GetHeader() http.Header {
	return self.Header
}
func (self *Param) GetPostBody() string {
	return self.PostBody
}
func (self *Param) GetRedirectTimes() int {
	return self.RedirectTimes
}
func (self *Param) GetDialTimeout() time.Duration {
	return self.DialTimeout
}
func (self *Param) GetConnTimeout() time.Duration {
	return self.ConnTimeout
}
func (self *Param) GetRetryPause() time.Duration {
	return self.RetryPause
}
func (self *Param) GetTryTimes() int {
	return self.TryTimes
}
func (self *Param) GetusePhantomJS() bool {
	return self.UsePhantomJS
}
