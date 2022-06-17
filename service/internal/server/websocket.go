package server

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"service/internal/server/util"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//channel通道
var ToClientChan chan clientInfo

// channel通道结构体
type clientInfo struct {
	ClientId   string
	SendUserId string
	MessageId  string
	Code       int
	Msg        string
	Data       *string
}

type RetData struct {
	MessageId  string      `json:"messageId"`
	SendUserId string      `json:"sendUserId"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

// 心跳间隔
var heartbeatInterval = 25 * time.Second

func init() {
	ToClientChan = make(chan clientInfo, 1000)
}

const (
	// 最大的消息大小
	maxMessageSize = 8192
)

var Manager = NewClientManager() // 管理者

func WsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("module", err)
		return
	}
	// 解析参数
	sessionId := r.FormValue("_sessionId")
	c.SetReadLimit(maxMessageSize) // 设置读取消息大小上线
	clientId := util.GenClientId()
	clientSocket := NewClient(clientId, sessionId, c)
	Manager.AddClient2SystemClient(sessionId, clientSocket)

	//读取客户端消息
	clientSocket.Read()

	// 用户连接事件
	Manager.Connect <- clientSocket
	go Manager.Start()
}

// 发送信息到指定客户端
func SendMessage2Client(clientId string, sendUserId string, code int, msg string, data *string) (messageId string) {
	messageId = util.GenUUID()
	// 如果是单机服务，则只发送到本机
	SendMessage2LocalClient(messageId, clientId, sendUserId, code, msg, data)
	return
}

//关闭客户端
func CloseClient(clientId, systemId string) {
	//如果是单机服务，则只发送到本机
	CloseLocalClient(clientId, systemId)
	return
}

//添加客户端到分组
func AddClient2Group(systemId string, groupName string, clientId string, userId string, extend string) {
	if client, err := Manager.GetByClientId(clientId); err == nil {
		//如果是单机，就直接添加到本地group了
		Manager.AddClient2LocalGroup(groupName, client, userId, extend)
	}
}

// 发送信息到指定分组
func SendMessage2Group(systemId, sendUserId, groupName string, code int, msg string, data *string) (messageId string) {
	messageId = util.GenUUID()
	// 如果是单机服务，则只发送到本机
	Manager.SendMessage2LocalGroup(systemId, messageId, sendUserId, groupName, code, msg, data)
	return
}

//发送信息到指定系统
func SendMessage2System(systemId, sendUserId string, code int, msg string, data string) {
	messageId := util.GenUUID()
	//如果是单机服务，则只发送到本机
	Manager.SendMessage2LocalSystem(systemId, messageId, sendUserId, code, msg, &data)
}

//获取分组列表
func GetOnlineList(systemId *string, groupName *string) map[string]interface{} {
	var clientList []string
	//如果是单机服务，则只发送到本机
	retList := Manager.GetGroupClientList(util.GenGroupKey(*systemId, *groupName))
	clientList = append(clientList, retList...)

	return map[string]interface{}{
		"count": len(clientList),
		"list":  clientList,
	}
}

// 通过本服务器发送信息
func SendMessage2LocalClient(messageId, clientId string, sendUserId string, code int, msg string, data *string) {
	log.WithFields(log.Fields{
		"clientId": clientId,
	}).Info("发送到通道")
	ToClientChan <- clientInfo{ClientId: clientId, MessageId: messageId, SendUserId: sendUserId, Code: code, Msg: msg, Data: data}
	return
}

// 发送关闭信号
func CloseLocalClient(clientId, systemId string) {
	if conn, err := Manager.GetByClientId(clientId); err == nil && conn != nil {
		if conn.SystemId != systemId {
			return
		}
		Manager.DisConnect <- conn
		log.WithFields(log.Fields{
			//"host":     setting.GlobalSetting.LocalHost,
			//"port":     setting.CommonSetting.HttpPort,
			"clientId": clientId,
		}).Info("主动踢掉客户端")
	}
	return
}

// 监听并发送给客户端信息
func WriteMessage() {
	for {
		clientInfo := <-ToClientChan
		log.WithFields(log.Fields{
			//"host":       setting.GlobalSetting.LocalHost,
			//"port":       setting.CommonSetting.HttpPort,
			"clientId":   clientInfo.ClientId,
			"messageId":  clientInfo.MessageId,
			"sendUserId": clientInfo.SendUserId,
			"code":       clientInfo.Code,
			"msg":        clientInfo.Msg,
			"data":       clientInfo.Data,
		}).Info("发送到本机")
		if conn, err := Manager.GetByClientId(clientInfo.ClientId); err == nil && conn != nil {
			if err := Render(conn.Socket, clientInfo.MessageId, clientInfo.SendUserId, clientInfo.Code, clientInfo.Msg, clientInfo.Data); err != nil {
				Manager.DisConnect <- conn
				log.WithFields(log.Fields{
					//"host":     setting.GlobalSetting.LocalHost,
					//"port":     setting.CommonSetting.HttpPort,
					"clientId": clientInfo.ClientId,
					"msg":      clientInfo.Msg,
				}).Error("客户端异常离线：" + err.Error())
			}
		}
	}
}

func Render(conn *websocket.Conn, messageId string, sendUserId string, code int, message string, data interface{}) error {
	return conn.WriteJSON(RetData{
		Code:       code,
		MessageId:  messageId,
		SendUserId: sendUserId,
		Msg:        message,
		Data:       data,
	})
}

//启动定时器进行心跳检测
func PingTimer() {
	go func() {
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()
		for {
			<-ticker.C
			//发送心跳
			for clientId, conn := range Manager.AllClient() {
				if err := conn.Socket.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
					Manager.DisConnect <- conn
					log.Errorf("发送心跳失败: %s 总连接数：%d", clientId, Manager.Count())
				}
			}
		}

	}()
}
