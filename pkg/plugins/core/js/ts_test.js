var Task = /** @class */ (function () {
    function Task(id) {
        this.id = id;
    }
    // 错误信息打印
    Task.prototype.error = function (message) {
        console.log("[error] ".concat(this.id), message);
    };
    // 普通日志打印
    Task.prototype.info = function (message) {
        console.log("[info] ".concat(this.id), message);
    };
    // 警告日志打印
    Task.prototype.waring = function (message) {
        console.log("[waring] ".concat(this.id), message);
    };
    return Task;
}());
var task = new Task('15');
task.info("你好");