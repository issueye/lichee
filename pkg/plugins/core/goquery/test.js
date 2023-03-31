var query = require('go/query')

let header = {
    "User-Agent": "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Mobile Safari/537.36 Edg/111.0.1661.54"
}

let doc = query.do('GET', 'https://www.oschina.net/news', header)

let selection = doc.find('div[id=newsList]>div>div')
if (selection.err) {
    throw selection.err;
}

selection.each(function (i, s) {
    let content = s.find("div>div>p").text()
    console.log(`num ${i}`, content.value)
})