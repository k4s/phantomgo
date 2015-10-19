# phantomgo
a headless browser phantomjs for golang
```
package main

import (
	"fmt"
	"github.com/k4s/phantomgo"

	//	"io"
	"io/ioutil"
	"os"
	// "log"
	// "os/exec"
)

func main() {
	p := phantomgo.NewPhantom()
	res, _ := p.Start([]string{"open", "http://weibo.com"})
	output, _ = ioutil.ReadAll(res)
	fmt.Println(string(output))
}
```
