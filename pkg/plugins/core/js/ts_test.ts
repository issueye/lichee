class Task {
    id: String;

    constructor(id: String) {
        this.id = id;
    }

    // 错误信息打印
    error(message) {
        console.log(`[error] ${this.id}`, message)
    }

    // 普通日志打印
    info(message) {
        console.log(`[info] ${this.id}`, message)
    }

    // 警告日志打印
    waring(message) {
        console.log(`[waring] ${this.id}`, message)
    }
}


var task = new Task('15')
task.info("你好")