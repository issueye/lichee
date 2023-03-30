var http = require('http');
var utils = require("utils");

let rootUrl = "http://127.0.0.1:10066/api"

function testGet() {
    let url = `${rootUrl}/test/get?name=${encodeURIComponent("get张三")}`;
    let headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36 Edg/111.0.1661.54"
    };
    let body = "";
    let rs = http.request("GET", url, headers, body);
    rs = utils.toString(rs.body);
    if (typeof rs == "string") {
        rs = JSON.parse(rs);
    }
    console.log("test get rs:", rs.data);
}

function testPost() {
    let url = `${rootUrl}/test/post`;
    let headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36 Edg/111.0.1661.54"
    };
    let body = {
        "name": "post李四"
    };
    let rs = http.request("POST", url, headers, JSON.stringify(body));
    rs = utils.toString(rs.body);
    if (typeof rs == "string") {
        rs = JSON.parse(rs);
    }
    console.log("test post rs:", rs.data);
}

function start() {
    testGet();
    testPost();
}

start();