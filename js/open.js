//system 用于
var system = require('system');
var page = require('webpage').create();
console.log(system.args[0],system.args[1],system.args[2])
page.settings.userAgent = 'Mozilla/5.0+(compatible;+Baiduspider/2.0;++http://www.baidu.com/search/spider.html)';
if(system.args.length ==1){
	phantom.exit();
}else{
	var url = system.args[1];
	var encode = system.args[2];

	if(encode != undefined){
		//设置编码
		phantom.outputEncoding=encode;
	}
	if(system.args[3] != undefined){
		
		//设置客户端代理设备
		page.settings.userAgent = system.args[3]
	}
	
	page.open(url, function (status) {
	    if (status !== 'success') {
	        console.log('Unable to access network');
	    } else {
	        // var ua = page.evaluate(function () {
	        //     return page.content;
	        // });
	        console.log(page.content);
	    }
	    phantom.exit();
	});
}

