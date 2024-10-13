// ws/ws.go
package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

const BREAK_MSG = "ws://break"

var CM = &ConnectionManager{}

// 升级器，用于将HTTP连接升级为WebSocket连接
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的连接
	},
}

// WebSocket连接结构体
type WebSocketConnection struct {
	Conn *websocket.Conn
}

type ConnectionManager struct {
	connections sync.Map // 存储多个连接(可以在多个地方接收ws消息)
}

// 添加连接
func (cm *ConnectionManager) AddConnection(id string, conn *websocket.Conn) {
	cm.connections.Store(id, &WebSocketConnection{Conn: conn})
}

// 删除连接
func (cm *ConnectionManager) RemoveConnection(id string) {
	cm.connections.Delete(id)
}

// 发送消息到指定连接
func (cm *ConnectionManager) SendMessage(id string, message string) error {
	connInterface, ok := cm.connections.Load(id)
	if !ok {
		return fmt.Errorf("connection %s not found", id)
	}

	wsConn := connInterface.(*WebSocketConnection)
	return wsConn.Conn.WriteMessage(websocket.TextMessage, []byte(message))
}

// 处理WebSocket连接
func HandleConnection(w http.ResponseWriter, r *http.Request, cm *ConnectionManager, id string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected:", id)
	cm.AddConnection(id, conn)

	// 连接建立时发送欢迎消息
	welcomeMessage := "Welcome to the WebSocket server!"
	err = cm.SendMessage(id, welcomeMessage)
	if err != nil {
		fmt.Println("Error sending welcome message:", err)
	}

	defer cm.RemoveConnection(id)

	for {
		messageType, msg, err := conn.ReadMessage() // 读取消息
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		if string(msg) == BREAK_MSG {
			fmt.Println("Client disconnected:", id)
			break
		}
		fmt.Printf("Received message from %s: %s\n", id, msg)

		// 将接收到的消息回显给客户端
		err = conn.WriteMessage(messageType, msg)
		if err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}
