# phantomgo
a headless browser phantomjs for golang

可打开需要加载js，或者防御爬虫的网站
```
package main

import (
	"fmt"
	"github.com/k4s/phantomgo"
	"io/ioutil"
)

func main() {
	p := phantomgo.NewPhantom()
	res, _ := p.Start([]string{"open", "http://weibo.com"})
	output, _ := ioutil.ReadAll(res)
	fmt.Println(string(output))
}
```
可动态执行phantomjs提供的JavaScript接口
```
package main

import (
	"fmt"
	"github.com/k4s/phantomgo"
	"io/ioutil"
)

func main() {
	p := phantomgo.NewPhantom()
	js := `
		var system = require('system');
		console.log(system.args[0],system.args[1],system.args[2]);
		phantom.exit();
		`
	res, _ := p.Exec(js,"11","22")
	output, _ := ioutil.ReadAll(res)
	fmt.Println(string(output))
	
}
```