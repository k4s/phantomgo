package phantomgo

import (
	"errors"
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
	SetPhantomjsPath(string, string)
	Download(Request) (*http.Response, error)
	Exec(string, ...string) (io.ReadCloser, error)
}

type Phantom struct {
	userAgent     string
	pageEncode    string
	phantomjsPath string
	WebrowseParam
}

//浏览器参数
type WebrowseParam struct {
	method      string
	url         string
	header      http.Header
	cookie      string
	postBody    string
	dialTimeout time.Duration //拨号超时时间段
	connTimeout time.Duration //链接超时时间
	tryTimes    int           //请求失败重新请求次数
	retryPause  time.Duration //请求失败时重复试时间段
}

func NewPhantom() Phantomer {
	phantom := &Phantom{
		userAgent:     "Mozilla/5.0+(compatible;+Baiduspider/2.0;++http://www.baidu.com/search/spider.html)",
		pageEncode:    "utf-8",
		phantomjsPath: GOPATH + "/src/github.com/k4s/phantomgo/phantomjs/phantomjs",
	}
	//如果js文件不存在,则创建文件
	if !phantom.Exist(GET_JS_FILE_NAME) {
		phantom.CreatJsFile("GET")
	}
	if !phantom.Exist(POST_JS_FILE_NAME) {
		phantom.CreatJsFile("POST")
	}
	return phantom
}

func (self *Phantom) Download(req Request) (resp *http.Response, err error) {

	//请求方法
	self.method = strings.ToUpper(req.GetMethod())
	//请求地址
	self.url = req.GetUrl()
	//请求http头
	self.header = req.GetHeader()
	//postDATA
	self.postBody = req.GetPostBody()
	//请求尝试次数
	self.tryTimes = req.GetTryTimes()
	//拨号超时时间
	self.dialTimeout = req.GetDialTimeout()
	//链接超时时间
	self.connTimeout = req.GetConnTimeout()
	//请求失败重新尝试的间隔时间
	self.retryPause = req.GetRetryPause()

	//请求的cookie
	for k, v := range self.header {
		if k == "Cookie" || k == "cookie" {
			for _, vv := range v {
				self.cookie = vv
			}
		}

	}

	var pagebody io.ReadCloser
	resp = new(http.Response)

	if self.method == "GET" {
		pagebody, err = self.Open(GET_JS_FILE_NAME, self.url, self.cookie, self.pageEncode, self.userAgent)
		if err != nil {
			return nil, err
		}
		resp.Status = "200 OK"
		resp.StatusCode = 200
		resp.Body = pagebody
		return
	} else if self.method == "POST" {
		pagebody, err = self.Open(POST_JS_FILE_NAME, self.url, self.cookie, self.pageEncode, self.userAgent, self.postBody)
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

//打开远程地址
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

//动态执行js
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

//设置本地代理客户端
func (self *Phantom) SetUserAgent(userAgent string) {
	self.userAgent = userAgent
}

//设置页面编码
func (self *Phantom) SetPageEncode(pageEncode string) {
	self.pageEncode = pageEncode
}

// 动态修改执行文件的Phantomjs.exe路径
func (self *Phantom) SetPhantomjsPath(name string, filepath string) {
	self.phantomjsPath = filepath
}

//创建js临时文件
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
func (self *Phantom) Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//销毁js临时文件
func (self *Phantom) DestroyJsFile(filename string) {
	os.Remove(filename)
}
