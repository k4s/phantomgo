package phantomgo

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

//javascript temp file name
const GET_JS_FILE_NAME = "get_jsfile_to_phantom"
const POST_JS_FILE_NAME = "post_jsfile_to_phantom"
const DIY_JS_FILE_NAME = "diy_jsfile_to_phantom"

var GOPATH = os.Getenv("GOPATH")

type Phantomer interface {
	SetUserAgent(string)
	SetProxy(string)
	SetProxyType(string)
	SetProxyAuth(string)
	SetPhantomjsPath(string, string)
	Download(Request) (*http.Response, error)
	Exec(string, ...string) (io.ReadCloser, error)
}

type Phantom struct {
	userAgent     string
	pageEncode    string
	phantomjsPath string
	proxy         string
	proxyType     string
	proxyAuth     string
	WebrowseParam
}

//web browse param
type WebrowseParam struct {
	method      string
	url         string
	header      http.Header
	cookie      string
	postBody    string
	dialTimeout time.Duration
	connTimeout time.Duration
	tryTimes    int           //if request failed,retry times
	retryPause  time.Duration //if request failed,retry time
}

func NewPhantom() Phantomer {
	phantom := &Phantom{
		userAgent:     "Mozilla/5.0+(compatible;+Baiduspider/2.0;++http://www.baidu.com/search/spider.html)",
		pageEncode:    "utf-8",
		phantomjsPath: GOPATH + "phantomjs",
	}
	//if the javascript file is no exist,creat it
	if !phantom.Exist(GET_JS_FILE_NAME) {
		phantom.CreatJsFile("GET")
	}
	if !phantom.Exist(POST_JS_FILE_NAME) {
		phantom.CreatJsFile("POST")
	}
	return phantom
}

func (self *Phantom) Download(req Request) (resp *http.Response, err error) {

	//request method
	self.method = strings.ToUpper(req.GetMethod())
	//request address
	self.url = req.GetUrl()
	//request http header
	self.header = req.GetHeader()
	//postDATA
	self.postBody = req.GetPostBody()
	//retry times
	self.tryTimes = req.GetTryTimes()
	//if request failed,retry time
	self.retryPause = req.GetRetryPause()
	self.dialTimeout = req.GetDialTimeout()
	self.connTimeout = req.GetConnTimeout()

	//set cookie
	for k, v := range self.header {
		if k == "Cookie" || k == "cookie" {
			for _, vv := range v {
				self.cookie = vv
			}
		}

	}

	var pagebody io.ReadCloser
	resp = new(http.Response)

	proxy, proxyType, proxyAuth := "", "", ""
	if self.proxy != "" {
		proxy = fmt.Sprintf("--proxy=%s ", self.proxy)
	}

	if self.proxyType != "" {
		proxyType = fmt.Sprintf("--proxy-type=%s ", self.proxyType)
	}

	if self.proxyAuth != "" {
		proxyAuth += fmt.Sprintf("--proxy-auth=%s ", self.proxyAuth)
	}

	if self.method == "GET" {
		pagebody, err = self.Open(proxy, proxyType, proxyAuth, GET_JS_FILE_NAME, self.url, self.cookie, self.pageEncode, self.userAgent)
		if err != nil {
			return nil, err
		}
		resp.Status = "200 OK"
		resp.StatusCode = 200
		resp.Body = pagebody
		return
	} else if self.method == "POST" {
		pagebody, err = self.Open(proxy, proxyType, proxyAuth, POST_JS_FILE_NAME, self.url, self.cookie, self.pageEncode, self.userAgent, self.postBody)
		if err != nil {
			return nil, err
		}
		resp.Status = "200 OK"
		resp.StatusCode = 200
		resp.Body = pagebody
		return
	}
	return nil, errors.New("Download error")
}

//open the url address
func (self *Phantom) Open(openArgs ...string) (stdout io.ReadCloser, err error) {
	cmd := exec.Command(self.phantomjsPath, openArgs...)
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return stdout, err
}

//exec javascript
func (self *Phantom) Exec(js string, args ...string) (stdout io.ReadCloser, err error) {
	file, _ := os.Create(DIY_JS_FILE_NAME)
	file.WriteString(js)
	file.Close()
	var exeCommand []string
	exeCommand = append(append(exeCommand, DIY_JS_FILE_NAME), args...)
	cmd := exec.Command(self.phantomjsPath, exeCommand...)
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return stdout, err

}

//SetUserAgent for example [chrome,firefox,IE..]
func (self *Phantom) SetUserAgent(userAgent string) {
	self.userAgent = userAgent
}

//SetProxy for example address:port
func (self *Phantom) SetProxy(proxy string) {
	self.proxy = proxy
}

//SetProxyType for example [http|socks5|none]
func (self *Phantom) SetProxyType(proxyType string) {
	self.proxyType = proxyType
}

//SetProxyAuth for example username:password
func (self *Phantom) SetProxyAuth(proxyAuth string) {
	self.proxyAuth = proxyAuth
}

//set web page decode for example [utf-8|gbk]
func (self *Phantom) SetPageEncode(pageEncode string) {
	self.pageEncode = pageEncode
}

// 动态修改执行文件的Phantomjs.exe路径
// set the phantomjs exec file
func (self *Phantom) SetPhantomjsPath(name string, filepath string) {
	self.phantomjsPath = filepath
}

//创建js临时文件
//creat temp javascript file
func (self *Phantom) CreatJsFile(jsfile string) {
	if jsfile == "GET" {
		js := getJs
		file, _ := os.Create(GET_JS_FILE_NAME)
		file.WriteString(js)
	} else if jsfile == "POST" {
		js := postJs
		file, _ := os.Create(POST_JS_FILE_NAME)
		file.WriteString(js)
	}

}

//判断js临时文件是否存在
//Is js file exist
func (self *Phantom) Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//销毁js临时文件
func (self *Phantom) DestroyJsFile(filename string) {
	os.Remove(filename)
}
