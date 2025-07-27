package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"miko-email/internal/result"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源，生产环境应该限制
		return true
	},
}

// WebSocket消息类型
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// 客户端连接
type Client struct {
	ID     string
	UserID int
	Conn   *websocket.Conn
	Send   chan WSMessage
}

// WebSocket管理器
type WSManager struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan WSMessage
	mutex      sync.RWMutex
}

// 全局WebSocket管理器
var wsManager = &WSManager{
	clients:    make(map[string]*Client),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	broadcast:  make(chan WSMessage),
}

// WebSocketHandler WebSocket处理器
type WebSocketHandler struct{}

// NewWebSocketHandler 创建WebSocket处理器
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{}
}

// 启动WebSocket管理器
func init() {
	go wsManager.run()
}

// 运行WebSocket管理器
func (manager *WSManager) run() {
	for {
		select {
		case client := <-manager.register:
			manager.mutex.Lock()
			manager.clients[client.ID] = client
			manager.mutex.Unlock()
			log.Printf("WebSocket client connected: %s (User: %d)", client.ID, client.UserID)

		case client := <-manager.unregister:
			manager.mutex.Lock()
			if _, ok := manager.clients[client.ID]; ok {
				delete(manager.clients, client.ID)
				close(client.Send)
				manager.mutex.Unlock()
				log.Printf("WebSocket client disconnected: %s", client.ID)
			} else {
				manager.mutex.Unlock()
			}

		case message := <-manager.broadcast:
			manager.mutex.RLock()
			for _, client := range manager.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(manager.clients, client.ID)
				}
			}
			manager.mutex.RUnlock()
		}
	}
}

// HandleWebSocket 处理WebSocket连接
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("WebSocket升级失败"))
		return
	}

	// 创建客户端
	client := &Client{
		ID:     generateClientID(),
		UserID: 0, // 初始未认证
		Conn:   conn,
		Send:   make(chan WSMessage, 256),
	}

	// 注册客户端
	wsManager.register <- client

	// 启动goroutines处理读写
	go client.writePump()
	go client.readPump()
}

// 生成客户端ID
func generateClientID() string {
	return time.Now().Format("20060102150405") + "_" + randomString(6)
}

// 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// 读取消息
func (c *Client) readPump() {
	defer func() {
		wsManager.unregister <- c
		c.Conn.Close()
	}()

	// 设置读取超时
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var message WSMessage
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// 处理消息
		c.handleMessage(message)
	}
}

// 写入消息
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 处理客户端消息
func (c *Client) handleMessage(message WSMessage) {
	switch message.Type {
	case "auth":
		// 处理认证
		if data, ok := message.Data.(map[string]interface{}); ok {
			if _, exists := data["token"].(string); exists {
				// 这里应该验证token并获取用户ID
				// 暂时简单处理
				c.UserID = 1 // 假设用户ID为1
				log.Printf("Client %s authenticated as user %d", c.ID, c.UserID)

				// 发送认证成功消息
				c.Send <- WSMessage{
					Type: "auth_success",
					Data: map[string]interface{}{
						"message": "认证成功",
						"user_id": c.UserID,
					},
				}
			}
		}

	case "ping":
		// 响应ping
		c.Send <- WSMessage{
			Type: "pong",
		}

	default:
		log.Printf("Unknown message type: %s", message.Type)
	}
}

// 广播消息给所有客户端
func BroadcastMessage(message WSMessage) {
	wsManager.broadcast <- message
}

// 发送消息给特定用户
func SendMessageToUser(userID int, message WSMessage) {
	wsManager.mutex.RLock()
	defer wsManager.mutex.RUnlock()

	for _, client := range wsManager.clients {
		if client.UserID == userID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(wsManager.clients, client.ID)
			}
		}
	}
}

// 通知新邮件
func NotifyNewEmail(userID int, from, subject string) {
	message := WSMessage{
		Type: "new_email",
		Data: map[string]interface{}{
			"from":    from,
			"subject": subject,
			"time":    time.Now(),
		},
	}
	SendMessageToUser(userID, message)
}

// 通知邮件发送成功
func NotifyEmailSent(userID int, to, subject string) {
	message := WSMessage{
		Type: "email_sent",
		Data: map[string]interface{}{
			"to":      to,
			"subject": subject,
			"time":    time.Now(),
		},
	}
	SendMessageToUser(userID, message)
}

// 通知转发规则触发
func NotifyForwardRuleTriggered(userID int, ruleName string) {
	message := WSMessage{
		Type: "forward_rule_triggered",
		Data: map[string]interface{}{
			"ruleName": ruleName,
			"time":     time.Now(),
		},
	}
	SendMessageToUser(userID, message)
}

// 通知系统消息
func NotifySystemMessage(message string) {
	wsMessage := WSMessage{
		Type: "system_notification",
		Data: map[string]interface{}{
			"message": message,
			"time":    time.Now(),
		},
	}
	BroadcastMessage(wsMessage)
}

// 获取在线用户数
func GetOnlineUserCount() int {
	wsManager.mutex.RLock()
	defer wsManager.mutex.RUnlock()
	return len(wsManager.clients)
}

// 获取在线用户列表
func GetOnlineUsers() []int {
	wsManager.mutex.RLock()
	defer wsManager.mutex.RUnlock()

	users := make(map[int]bool)
	for _, client := range wsManager.clients {
		if client.UserID > 0 {
			users[client.UserID] = true
		}
	}

	userList := make([]int, 0, len(users))
	for userID := range users {
		userList = append(userList, userID)
	}

	return userList
}
