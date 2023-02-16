package service

import (
	telephone "Hwgen/app/controller"
	"Hwgen/global"
	"fmt"
	"github.com/gofrs/uuid"
	"net"
	"strings"
	"time"
)

// ClientManager 客户端管理
type ClientManager struct {
	clients    map[*Client]bool //客户端 map 储存并管理所有的长连接client，在线的为true，不在的为false
	broadcast  chan []byte      //web端发送来的的message我们用broadcast来接收，并最后分发给所有的client
	register   chan *Client     //新创建的长连接client
	unregister chan *Client     //新注销的长连接client
}

type Client struct {
	uuid uuid.UUID
	conn net.Conn
	send chan []byte
}

var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}
var (
	method telephone.Origin
)

func Run() {
	listener, err := net.Listen("tcp", "0.0.0.0:"+fmt.Sprintf("%d", global.H_CONFIG.System.Addr))
	if err != nil {
		fmt.Println("Listen tcp server failed,err:", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listen.Accept failed,err:", err)
			continue
		} else {
			client := &Client{
				uuid: uuid.Must(uuid.NewV4()),
				conn: conn,
				send: make(chan []byte),
			}
			go manager.Start()
			// manager.Quantity()

			manager.register <- client
			go client.Process()
			go client.Ping()
		}
	}
}

func (manager *ClientManager) Start() {
LOOP:
	for {
		select {
		case conn := <-manager.register: //如果有新的连接接入,就通过channel把连接传递给conn
			fmt.Println("Register uuid", conn.uuid)
			manager.clients[conn] = true
			manager.Quantity(len(manager.clients))
		case conn := <-manager.unregister: //断开连接时
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
			}
			fmt.Println("Client", conn.uuid)
			manager.Quantity(len(manager.clients))
			break LOOP
		}
	}
	fmt.Println("has been Destroy!")
}

func (manager *ClientManager) Quantity(ll int) {
	fmt.Println("Client quantity->", ll)
}

type Message struct {
	Type     int
	Describe string
	Content  string
	Pid      int
	Sid      int
}

func (c *Client) Process() bool {
	// conn := c.conn
	defer c.conn.Close()
	for {
		var buf [512]byte //接受数据
		n, err := c.conn.Read(buf[:])
		if err != nil {
			fmt.Printf("read from connect failed, ERR: %v\n", err)
			manager.unregister <- c
			break
		}
		originstr := string(buf[:n])
		originstr = strings.Replace(originstr, " ", "", -1)
		fmt.Println("origin", originstr)
		// piece1 := originstr[0:4]
		piece2 := originstr[4:6]
		piece3 := originstr[6:10]

		var instruction string
		switch piece2 {
		case "05":
			fmt.Println("网络连接状态查询")
			piece3 = piece3 + "1"
			instruction = "0011" + piece2 + piece3
		case "10":
			fmt.Println("公话认证")
			instruction, _ = method.Operation_10(originstr)
		case "04":
			fmt.Println("学生签到记录")
			piece3 = piece3 + "1"
			instruction = "0011" + piece2 + piece3
		case "01":
			fmt.Println("获取亲情号码")
			instruction, _ = method.Operation_01(originstr)
		case "03":
			fmt.Println("亲情电话订单处理")
			instruction, _ = method.Operation_03(originstr)
		case "13":
			fmt.Println("计费通话订单处理")
			instruction, _ = method.Operation_13(originstr)
		case "75":
			fmt.Println("获取通话费率")
			instruction, _ = method.Operation_75(originstr)
		case "81":
			fmt.Println("公话状态告警")
			instruction, _ = method.Operation_81(originstr)
		case "82":
			fmt.Println("获取公话状态")
			instruction, _ = method.Operation_82(originstr)
			continue
		default:
			fmt.Println("未识别指令", string(buf[:n]))
			continue
		}
		fmt.Println("last instruction", instruction)
		// if _, err = conn.Write([]byte(instruction)); err != nil {
		// 	fmt.Printf("write to client failed, err: %v\n", err)
		// 	break
		// }

		err = W(c.conn, instruction)
		if err != nil {
			fmt.Println("c.conn:", c.conn)
			fmt.Println("Write failed,ERR:", err)
			break
		}
	}
	return true
}

func (c *Client) Ping() {
	for {
		// time.Sleep(1 * time.Hour)
		time.Sleep(1 * time.Minute)
		fmt.Println("Ping...", c.uuid)
		instruction, _ := method.TelephoneState()
		err := W(c.conn, instruction)
		if err != nil {
			fmt.Println("Write failed,ERR:", err)
			break
		}
	}
}

func W(conn net.Conn, msg string) error {
	_, err := conn.Write([]byte(msg))
	return err
}
