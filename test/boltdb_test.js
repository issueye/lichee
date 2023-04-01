var bolt = require('db/bolt')
var time = require('time')


let BUCKET_NAME = 'test_bucket'

// 创建一个 bucket
bolt.createBucketIfNotExists(BUCKET_NAME)

// 写入数据
bolt.update(BUCKET_NAME, function (db) {

    for (i = 0; i < 20; i++) {
        // 写入

        let data = {
            code: 200,
            index: i,
            message: time.nowString(),
        }
        db.put(`test_00${i + 1}`, JSON.stringify(data))
    }
})

bolt.view(BUCKET_NAME, function (db) {
    // 读取
    let val = db.get('test_002')
    if (val.err) {
        throw err;
    }

    console.log('val', val.value);

    db.foreach(function (k, v) {
        console.log(`${k}`, v);
    })
})