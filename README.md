# phantomgo
### a headless browser phantomjs for golang

### it is easy to use for you web download, eg: [http://www.github.com/k4s/webrowser](http://www.github.com/k4s/webrowser)
```go
import (
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/k4s/phantomgo"
)

func main() {
	p := &Param{
		Method: "POST",  //POST or GET ..
		Url:    "http://localhost/go_test/1.php",
		Header:       http.Header{"Cookie": []string{"your cookies"}},
		UsePhantomJS: true,
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
### make the phantomjs yourself to the phantomjs javaScript interface:
```go
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
### a example simulate login by cookies, so you can get web data login after:
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
var page = require('webpage').create(),
  system = require('system'),
  address;
page.settings.userAgent = 'Mozilla/5.0+(compatible;+Baiduspider/2.0;++http://www.baidu.com/search/spider.html)';
phantom.cookiesEnabled = true;

phantom.addCookie({
  'name'     : 'Apache(do yourself)',
  'value'    : '63354989(do yourself)',
  'domain'     :'.weibo.com'});
phantom.addCookie({
  'name'     : 'SINAGLOBAL(do yourself)',
  'value'    : '8156705307(do yourself)',
  'domain'     :'.weibo.com'});
phantom.addCookie({
  'name'     : 'SUB(do yourself)',
  'value'    : '_2A257NF5qDeTxGeNK6VUT8izMzjmIHXVY12Ii(do yourself)',
  'domain'     :'.weibo.com'});
phantom.addCookie({
  'name'     : 'SUBP(do yourself)',
  'value'    : '0033WrSXqPxfM725Ws9jqgMF55529P9D9WhnL77(do yourself)',
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