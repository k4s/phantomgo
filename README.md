# phantomgo
a headless browser phantomjs for golang

全新的框架,更容易使用,更容易嵌套到自己的下载器
eg:[http://www.github.com/k4s/webrowser](http://www.github.com/k4s/webrowser)
```
import (
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/k4s/phantomgo"
)

func main() {
	p := &Param{
		Method: "POST",  //有POST，GET方法
		Url:    "http://localhost/go_test/1.php",
		Header:       http.Header{"Cookie": []string{"your cookies"}},
		UsePhontomJS: true,
		PostBody:     "aaa=111",
	}
	brower := NewPhantom()
	resp, err := brower.Download(p)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
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

模拟登录，新浪微博轻松抓数据
```
package main

import (
	//	"time"
	"fmt"
	"github.com/k4s/phantomgo"
	"io/ioutil"
)

func main() {
	p := phantomgo.NewPhantom()

	js := `
var page = require('webpage').create(),
  system = require('system'),
  address;
page.settings.userAgent = 'Mozilla/5.0+(compatible;+Baiduspider/2.0;++http://www.baidu.com/search/spider.html)';
phantom.cookiesEnabled = true;

phantom.addCookie({
  'name'     : 'Apache(换成自己的)',
  'value'    : '63354989(换成自己的)',
  'domain'     :'.weibo.com'});
phantom.addCookie({
  'name'     : 'SINAGLOBAL(换成自己的)',
  'value'    : '8156705307(换成自己的)',
  'domain'     :'.weibo.com'});
phantom.addCookie({
  'name'     : 'SUB(换成自己的)',
  'value'    : '_2A257NF5qDeTxGeNK6VUT8izMzjmIHXVY12Ii(换成自己的)',
  'domain'     :'.weibo.com'});
phantom.addCookie({
  'name'     : 'SUBP(换成自己的)',
  'value'    : '0033WrSXqPxfM725Ws9jqgMF55529P9D9WhnL77(换成自己的)',
  'domain'     :'.weibo.com'});
if (system.args.length === 1) {
  phantom.exit(1);
} else {
  address = system.args[1];
  page.open(address, function (status) {
    console.log(page.content);
    phantom.exit();
  });
}
`
	res, _ := p.Exec(js, "http://weibo.com/55555555/fans?rightmod=1&wvr=6")
	output, _ := ioutil.ReadAll(res)
	fmt.Println(string(output))

}

```