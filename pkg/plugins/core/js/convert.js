var utils = require('utils')

utils.arrayToMap = function (key, arr) {
    var data = {};
    for (let index = 0; index < arr.length; index++) {
        const element = arr[index];
        if (element[key]) {
            data[element[key]] = element
        }
    }
    return data
}