package ws

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/issueye/lichee/utils"
	"github.com/tidwall/gjson"
)

const (
	// PingPeriod 心跳检测 时间间隔
	PingPeriod = 5 * time.Second
)

const (
	WS_CALL_MACHINE = "CALL_MACHINE"
	WS_BIG_SCREEN   = "YxPDJHAZXP_BY"
)

// WsResponse 响应结构
type WsResponse struct {
	Msg          string `json:"msg"`
	SendDatetime string `json:"sendDatetime"`
}

type UnReg struct {
	Group string
	Id    string
}

var UnRegChan = make(chan *UnReg, 10)

// WsResponseData 返回数据
type WsResponseData struct {
	WsResponse
	Data interface{} `json:"data"`
}

func (res *WsResponse) ToJson() string {
	return utils.Ljson{}.Struct2Json(res)
}

// WebsocketManager 初始化 wsManager 管理器
var WebsocketManager = Manager{
	Group:            make(map[string]map[string]*Client),
	Register:         make(chan *Client, 128),
	UnRegister:       make(chan *Client, 128),
	GroupMessage:     make(chan *GroupMessageData, 128),
	Message:          make(chan *MessageData, 128),
	BroadCastMessage: make(chan *BroadCastMessageData, 128),
	groupCount:       0,
	clientCount:      0,
}

func RunWs(log *zap.SugaredLogger) {
	WebsocketManager.log = log

	go WebsocketManager.Start()
	go WebsocketManager.SendService()
	go WebsocketManager.SendGroupService()
	go WebsocketManager.SendGroupService()
	go WebsocketManager.SendAllService()
	go WebsocketManager.SendAllService()

	// 心跳检测
	go WebsocketManager.Ping()
}

type Manager struct {
	Group                   map[string]map[string]*Client
	groupCount, clientCount uint
	Lock                    sync.Mutex
	Register, UnRegister    chan *Client
	Message                 chan *MessageData
	GroupMessage            chan *GroupMessageData
	BroadCastMessage        chan *BroadCastMessageData
	RegisterEvent           ClientEventFunc // 注册事件
	UnRegisterEvent         ClientEventFunc // 注销事件

	log *zap.SugaredLogger // 日志对象
}

// ClientEventFunc
// t 0 注册 1 注销
type ClientEventFunc func(t int64, client *Client)

// Client 单个 websocket 信息
type Client struct {
	Id          string    // 唯一标识
	Group       string    // 组名
	IP          string    // ip
	LastPing    time.Time // 最后一次心跳
	PingLoseNum int       // 心跳检测失败次数
	Socket      *websocket.Conn
	Message     chan []byte
}

// MessageData
// messageData 单个发送数据信息
type MessageData struct {
	Id      string
	Group   string
	Message []byte
}

// GroupMessageData
// groupMessageData 组广播数据信息
type GroupMessageData struct {
	Group   string
	Message []byte
}

// BroadCastMessageData
// 广播发送数据信息
type BroadCastMessageData struct {
	Message []byte
}

func writeLog(msg string) {
	if WebsocketManager.log != nil {
		WebsocketManager.log.Debug(msg)
	} else {
		log.Print(msg)
	}
}

// 服务端向客户端发送消息
func S2CConsole(id, msg string) {
	writeLog(fmt.Sprintf("[%s]  [服务端]    >>    [客户端 : %s]    ::  发送消息：%s", time.Now().Format(utils.FormatDateTimeMs), id, msg))
}

// 客户端向服务端发送消息
func C2SConsole(id, msg string) {
	writeLog(fmt.Sprintf("[%s]  [客户端 : %s]    >>    [服务端]    ::  发送消息：%s", time.Now().Format(utils.FormatDateTimeMs), id, msg))
}

func Console(msg string) {
	writeLog(fmt.Sprintf("[%s]  %s", time.Now().Format(utils.FormatDateTimeMs), msg))
}

// Read
// 读信息，从 websocket 连接直接读取数据
func (c *Client) Read() {
	defer func() {
		WebsocketManager.UnRegister <- c
		Console(fmt.Sprintf("客户端[%s]断开连接", c.Id))
		if err := c.Socket.Close(); err != nil {
			Console(fmt.Sprintf("客户端[%s]断开连接失败，失败原因：%s", c.Id, err))
		}
	}()

	for {
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}
		C2SConsole(c.IP, string(message))

		// {"type":"heartbeat"}
		value := gjson.Get(string(message), "type")
		if value.Exists() {
			if value.String() == "heartbeat" {
				Console(fmt.Sprintf("客户端[%s]心跳检测，心跳内容【%s】", c.IP, &value))
			}
		}

		c.Message <- message
	}
}

// Write
// 写信息，从 channel 变量 Send 中读取数据写入 websocket 连接
func (c *Client) Write() {
	defer func() {
		Console(fmt.Sprintf("客户端【%s】断开连接", c.IP))
		if err := c.Socket.Close(); err != nil {
			Console(fmt.Sprintf("客户端【%s】断开连接失败，失败原因: %s", c.IP, err))
		}
	}()

	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if strings.ToUpper(string(message)) == "PING" {
				S2CConsole(c.IP, "心跳检测")
				err := c.Socket.WriteMessage(websocket.TextMessage, []byte(`{"code": 200, "type": "heartbeat", "message":"pong"}`))
				if err != nil {
					S2CConsole(c.IP, fmt.Sprintf("发送消息失败，失败原因: %s", err.Error()))
				}
			} else {
				err := c.Socket.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					S2CConsole(c.IP, fmt.Sprintf("发送消息失败，失败原因:  %s", err.Error()))
				}
			}

		}
	}
}

// Ping
// 心跳检测
func (manager *Manager) Ping() {

	openPing := os.Getenv("OPEN_PING")
	if openPing == "false" || openPing == "" {
		return
	}

	for {
		time.Sleep(PingPeriod)
		// 遍历组，检测心跳
		for _, v := range manager.Group {
			// 判断如果当前组的心跳时间超过了10秒，则认为心跳失败
			for _, c := range v {
				err := c.Socket.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"code": 200, "msg": "ping", "sendDatetime":"%s"}`, time.Now().Format("2006-01-02 15:04:05"))))
				if err != nil {
					S2CConsole(c.IP, fmt.Sprintf("写入数据失败,失败原因: %s", err.Error()))
					c.PingLoseNum++
					c.LastPing = time.Now()
					if c.PingLoseNum > 3 {
						WebsocketManager.UnRegister <- c
						Console(fmt.Sprintf("客户端【%s】断开连接", c.IP))
						if err := c.Socket.Close(); err != nil {
							Console(fmt.Sprintf("客户端【%s】断开连接失败，失败原因: %s", c.IP, err))
						}
					}
				} else {
					c.PingLoseNum = 0
					c.LastPing = time.Now()
				}
			}
		}
	}
}

// Start
// 启动 websocket 管理器
func (manager *Manager) Start() {
	fmt.Println("WEBSOCKET 管理器启动...")
	for {
		select {
		// 注册
		case client := <-manager.Register:
			{
				manager.Lock.Lock()
				if manager.Group[client.Group] == nil {
					manager.Group[client.Group] = make(map[string]*Client)
					manager.groupCount += 1
				}
				manager.Group[client.Group][client.Id] = client
				manager.clientCount += 1
				manager.Lock.Unlock()
				if manager.RegisterEvent != nil {
					manager.RegisterEvent(0, client)
				}
				Console(fmt.Sprintf("客户端【%s】注册到【%s】组", client.IP, client.Group))
			}

		// 注销
		case client := <-manager.UnRegister:
			Console(fmt.Sprintf("客户端【%s】从【%s】组注销", client.IP, client.Group))
			manager.Lock.Lock()
			if manager.UnRegisterEvent != nil {
				manager.UnRegisterEvent(1, client)
			}
			var ok bool
			_, ok = manager.Group[client.Group]
			if ok {
				_, ok = manager.Group[client.Group][client.Id]
				if ok {
					close(client.Message)
					delete(manager.Group[client.Group], client.Id)
					manager.clientCount -= 1
					if len(manager.Group[client.Group]) == 0 {
						delete(manager.Group, client.Group)
						manager.groupCount -= 1
					}
				}
			}

			// 通知外部
			UnRegChan <- &UnReg{
				Group: client.Group,
				Id:    client.Id,
			}

			manager.Lock.Unlock()
		}
	}
}

// Close
// 关闭 websocket 管理器
func (manager *Manager) Close() {
	Console("WEBSOCKET 服务关闭")
	// 关闭注册管道
	close(manager.Register)
	close(manager.UnRegister)
	// 关闭消息管道
	close(manager.Message)
	close(manager.GroupMessage)
	close(manager.BroadCastMessage)
}

// SendService
// 处理单个 client 发送数据
func (manager *Manager) SendService() {
	for {
		select {
		case data := <-manager.Message:
			groupMap, ok := manager.Group[data.Group]
			if ok {
				var conn *Client
				conn, ok = groupMap[data.Id]
				if ok {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// SendGroupService
// 处理 group 广播数据
func (manager *Manager) SendGroupService() {
	for {
		select {
		// 发送广播数据到某个组的 channel 变量 Send 中
		case data := <-manager.GroupMessage:
			if groupMap, ok := manager.Group[data.Group]; ok {
				for _, conn := range groupMap {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// SendAllService
// 处理广播数据
func (manager *Manager) SendAllService() {
	for {
		select {
		case data := <-manager.BroadCastMessage:
			for _, v := range manager.Group {
				for _, conn := range v {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// Send
// 向指定的 client 发送数据
func (manager *Manager) Send(id string, group string, message []byte) {
	data := &MessageData{
		Id:      id,
		Group:   group,
		Message: message,
	}
	manager.Message <- data
}

// SendGroup
// 向指定的 Group 广播
func (manager *Manager) SendGroup(group string, message []byte) {
	data := &GroupMessageData{
		Group:   group,
		Message: message,
	}
	manager.GroupMessage <- data
}

// SendAll
// 广播
func (manager *Manager) SendAll(message []byte) {
	data := &BroadCastMessageData{
		Message: message,
	}
	manager.BroadCastMessage <- data
}

// RegisterClient
// 注册
func (manager *Manager) RegisterClient(client *Client) {
	manager.Register <- client
}

// UnRegisterClient
// 注销
func (manager *Manager) UnRegisterClient(client *Client) {
	manager.UnRegister <- client
}

// LenGroup
// 当前组个数
func (manager *Manager) LenGroup() uint {
	return manager.groupCount
}

// LenClient
// 当前连接个数
func (manager *Manager) LenClient() uint {
	return manager.clientCount
}

// Info
// 获取 wsManager 管理器信息
func (manager *Manager) Info() map[string]interface{} {
	managerInfo := make(map[string]interface{})
	managerInfo["groupLen"] = manager.LenGroup()
	managerInfo["clientLen"] = manager.LenClient()
	managerInfo["chanRegisterLen"] = len(manager.Register)
	managerInfo["chanUnregisterLen"] = len(manager.UnRegister)
	managerInfo["chanMessageLen"] = len(manager.Message)
	managerInfo["chanGroupMessageLen"] = len(manager.GroupMessage)
	managerInfo["chanBroadCastMessageLen"] = len(manager.BroadCastMessage)
	return managerInfo
}

// WsClient
// gin 处理 websocket handler
func (manager *Manager) WsClient(ctx *gin.Context, group string, args ...interface{}) (id string) {
	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("WEBSOCKET 连接失败")
		return ""
	}

	var ip string

	if len(args) > 0 {
		ip = args[0].(string)
	} else {
		ip = ctx.ClientIP()
	}

	client := &Client{
		Id:      utils.Lid{}.GetUUID(),
		Group:   group,
		IP:      ip, // 以 ip 作为 对应窗口的标识
		Socket:  conn,
		Message: make(chan []byte, 1024),
	}

	fmt.Printf("客户端【%s】注册到服务器 \n", ip)
	manager.RegisterClient(client)
	go client.Read()
	go client.Write()

	return client.Id
}
