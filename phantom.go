package phantomgo

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

var GOPATH = os.Getenv("GOPATH")

type Phantomer interface {
	SetUserAgent(string)
	SetFilePath(string, string)
	Start(args []string) (io.ReadCloser, error)
}

type Phantom struct {
	userAgent string
	filePath  map[string]string
}

func NewPhantom() Phantomer {
	phantom := &Phantom{filePath: map[string]string{
		"phantomjs": GOPATH + "/src/github.com/k4s/phantomgo/phantomjs/phantomjs",
		"openjs":    GOPATH + "/src/github.com/k4s/phantomgo/js/open.js",
	}}
	return phantom
}

func (self *Phantom) Start(args []string) (result io.ReadCloser, err error) {
	if len(args) == 2 {
		//exec 是执行js命令
		if args[0] == "exec" {
			result, err = self.Exec(args[1])
			if err != nil {
				return nil, err
			}
			return result, err
		} else if args[0] == "open" {
			result, err = self.Open(self.filePath["openjs"], args[1])
			if err != nil {
				return nil, err
			}

			return result, err

		}
	} else if len(args) == 3 {
		//open 是打开url链接
		if args[0] == "open" {

			if self.userAgent != "" {
				//args[1] 为请求地址
				//args[2] 为页面编码
				//userAgent 为客户端代理请求设备，默认为百度爬虫
				result, err := self.Open(self.filePath["openjs"], args[1], args[2], self.userAgent)
				if err != nil {
					return nil, err
				}
				return result, err
			} else {
				result, err := self.Open(self.filePath["openjs"], args[1], args[2])
				if err != nil {
					return nil, err
				}
				return result, err
			}
		}
	}
	return nil, errors.New("args error")
}

func (self *Phantom) Open(openArgs ...string) (stdout io.ReadCloser, err error) {
	cmd := exec.Command(self.filePath["phantomjs"], openArgs...)
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

func (self *Phantom) Exec(js string) (stdout io.ReadCloser, err error) {
	return stdout, err
}

func (self *Phantom) SetUserAgent(userAgent string) {
	self.userAgent = userAgent
}

// 动态修改执行文件路径
func (self *Phantom) SetFilePath(name string, filepath string) {
	self.filePath[name] = filepath
}
