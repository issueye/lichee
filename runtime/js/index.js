var utils = require('utils');
var url = require("url");

var lichee = {
    // 获取路由  获取最后一个路由节点
    getUrl: function () {
        console.log("请求路由:",request.getUrl().getPath());
        var path = request.getUrl().getPath();
        if (path.indexOf("/api") == 0) {
            path = path.substring(5);
        }
        return path.split("/")
    },
    getRequetInfo: function() {
        let info = {
            bodyString: request.getBodyString().value,
            query: JSON.parse(url.parseQuery(request.getUrl().getRawQuery()).value),
            cookies: request.cookies(),
            headers: request.getHeaders(),
        };
        // console.log("request info:",JSON.stringify(info));
        return info;
    },
    getRequestBodyString: function() {
        return request.getBodyString().value;
    },
    getRequestParams: function(){
        url.parse(request.getUrl().getPath())
    },
    main: function () {
        let route = this.getUrl();
        let str = route.map(v => v).join(".");
        let e = `controller.${str}(this.getRequetInfo())`;
        console.log("调用函数:", e);
        try {
            eval(e);   
        } catch (error) {
            response.write(JSON.stringify({
                code: 500,
                errMsg: "没找到路由"
            }));   
        }
    }
}

// 路由实现，规则按/拆分，如：访问:/api/hello/world，函数实现如下
var controller = {
    "hello": {
        "world": function () {
            response.write(JSON.stringify({
                code: 200,
                data: "hello world!"
            })); 
        }
    },
    "test": {
        "get": function (requestInfo) {
            // console.log("request info:",JSON.stringify(requestInfo));
            let name = requestInfo.query.name;
            if (name) {
                response.write(JSON.stringify({
                    code: 200,
                    data: `你好${name}`
                }));
            } else {
                response.write(JSON.stringify({
                    code: 500,
                    errMsg: "缺少name参数"
                }));
            }
        },
        "post": function(requestInfo) {
            // console.log("request info:",JSON.stringify(requestInfo));
            let name = "";
            if (requestInfo.bodyString.length > 2) {
                name = JSON.parse(requestInfo.bodyString).name || "";
            }
            if (name) {
                response.write(JSON.stringify({
                    code: 200,
                    data: name
                }));
            } else {
                response.write(JSON.stringify({
                    code: 500,
                    errMsg: name
                }));
            }
        }
    }
}

lichee.main();