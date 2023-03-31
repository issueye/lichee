var query = require('go/query')

let header = {
    "Cookie": "_ga=GA1.1.944084037.1651746013; _clck=1vnk3ca|1|f1u|0; yp_riddler_id=5a7ce9bd-f608-4054-b9fe-21831537a77e; oscid=OrZwnxvdOPW5ir7ppqH/m1fasSzPg8OqJnfoUk9bNjOMCh4zRY9c25KF/5eXuKY5mNvIajpXaWsRopMRqQI8j1uu/pagSmsohEnWcddfjNZMG/WcSJzZeI7PoOASoRYNIla3dLYJkT2wB8Fp7kLuPLAHwWnuQu486wjL1FTgO0c=; _user_behavior_=cfe82fc3-623a-44f7-9dc9-ffc01bbeedd0; Hm_lvt_a411c4d1664dd70048ee98afe7b28f0b=1680221445; _ga_TK89C9ZD80=GS1.1.1680221444.17.1.1680221740.0.0.0; Hm_lpvt_a411c4d1664dd70048ee98afe7b28f0b=1680221752",
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