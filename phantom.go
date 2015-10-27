package phantomgo

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

//javascript temp file name
const JS_FILE_NAME = "tem_js_to_phantom"

var GOPATH = os.Getenv("GOPATH")

type Phantomer interface {
	SetUserAgent(string)
	SetPhantomjsPath(string, string)
	Start(args []string) (io.ReadCloser, error)
	Exec(string, ...string) (io.ReadCloser, error)
}

type Phantom struct {
	userAgent     string
	jsFileName    string
	phantomjsPath string
}

func NewPhantom() Phantomer {
	phantom := &Phantom{
		jsFileName:    JS_FILE_NAME,
		phantomjsPath: GOPATH + "/src/github.com/k4s/phantomgo/phantomjs/phantomjs",
	}
	phantom.CreatJsFile()
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
			result, err = self.Open(self.jsFileName, args[1])
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
				result, err := self.Open(self.jsFileName, args[1], args[2], self.userAgent)
				if err != nil {
					return nil, err
				}
				return result, err
			} else {
				result, err := self.Open(self.jsFileName, args[1], args[2])
				if err != nil {
					return nil, err
				}
				return result, err
			}
		}
	}
	return nil, errors.New("args error")
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
//js为执行代码，args为命令行参数
func (self *Phantom) Exec(js string, args ...string) (stdout io.ReadCloser, err error) {
	file, _ := os.Create(self.jsFileName)
	file.WriteString(js)
	file.Close()
	var exeCommand []string
	exeCommand = append(append(exeCommand, self.jsFileName), args...)
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

// 动态修改执行文件路径
func (self *Phantom) SetPhantomjsPath(name string, filepath string) {
	self.phantomjsPath = filepath
}

//创建js临时文件
func (self *Phantom) CreatJsFile() {
	js := loadjs()
	file, _ := os.Create(self.jsFileName)
	file.WriteString(js)

}

func (self *Phantom) DestroyJsFile() {
	os.Remove(self.jsFileName)
}
