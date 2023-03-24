var local = require('db/local')
var other = require('db/other')

var dbTrans = function (db) {
    return function (callback) {
        if (!callback) {
            return;
        }
        var tx = db.begin()
        if (tx.err) {
            throw tx.err;
        }
        var tx = tx.value
        try {
            var err = callback(tx)
            if (err) {
                throw (err)
            }
            tx.commit()
        } catch (e) {
            tx.rollback()
        }
    }
}

// 将方法注入到对象
local.transaction = dbTrans(local)
other.transaction = dbTrans(other)