package service

import (
	telephone "Hwgen/app/controller"
	"fmt"
	"github.com/gofrs/uuid"
	"net"
	"strconv"
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
	id   int
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

func Run() {
	listener, err := net.Listen("tcp", "0.0.0.0:10087")
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
			go manager.start()
			manager.register <- client
			// go client.Read()
			go client.Process()
			go client.Write()
			go client.Ping()
		}
	}
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register: //如果有新的连接接入,就通过channel把连接传递给conn
			fmt.Println("Register uuid", conn.uuid)
			manager.clients[conn] = true
			fmt.Println("Client quantity", len(manager.clients))
		case conn := <-manager.unregister: //断开连接时
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
			}
			fmt.Println("Disconnected", conn.uuid)
			fmt.Println("Client quantity ", len(manager.clients))
			return
		}
	}
}

type Message struct {
	Type     int
	Describe string
	Content  string
	Pid      int
	Sid      int
}

func (c *Client) Process() bool {
	conn := c.conn
	defer conn.Close()
	for {
		var buf [128]byte
		//接受数据
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("read from connect failed, err: %v\n", err)
			break
		}
		originstr := string(buf[:n])
		originstr = strings.Replace(originstr, " ", "", -1)
		fmt.Println("origin", originstr)

		piece1 := originstr[0:4]
		piece2 := originstr[4:6]
		piece3 := originstr[6:10]
		fmt.Println("piece1", piece1)
		fmt.Println("piece2", piece2)
		fmt.Println("piece3", piece3)
		var instruction string
		switch piece2 {
		case "05": //网络连接状态查询
			fmt.Println("网络连接状态查询")
			instruction = piece1 + piece2 + piece3
		case "10": //验证
			fmt.Println("验证")
			piece3 = piece3 + "1"
			instruction = piece1 + piece2 + piece3
		case "04": //登录
			fmt.Println("登录")
			piece3 = piece3 + "1"
			instruction = piece1 + piece2 + piece3
		case "01": //亲情号码
			instruction, err = telephone.Operation_01(originstr)
			if err != nil {
				fmt.Println("instructionErr", err)
			}
		case "03": //话单
			piece3 = piece3 + "1"
			instruction = piece1 + piece2 + piece3
		default:
			return false
		}
		fmt.Println("last instruction", instruction)
		if _, err = conn.Write([]byte(instruction)); err != nil {
			fmt.Printf("write to client failed, err: %v\n", err)
			break
		}
	}
	return true
}

func Hex2Dec(val string) int {
	n, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		fmt.Println(err)
	}
	return int(n)
}

// func (c *Client) Read() {
// 	defer c.conn.Close()
// 	reader := bufio.NewReader(c.conn)
// 	for {
// 		msg, err := proto.Decode(reader)
// 		if err == io.EOF {
// 			return
// 		}
// 		if err != nil {
// 			fmt.Println("decode msg failed, err:", err)
// 			return
// 		}
// 		// var Msg Message
// 		// err = json.Unmarshal([]byte(msg), &Msg)
// 		// if err != nil {
// 		// 	fmt.Println("json.Unmarshal error:", err)
// 		// }
// 		fmt.Println("Recived from client,data", msg)
// 		// if Msg.Type == 2 {
// 		// 	c.id = Msg.Pid
// 		// 	go c.Operations(c.id)
// 		// }
// 	}
// }

func (c *Client) Operations(id int) {
	for conn := range manager.clients {
		fmt.Println("conn.id", conn.id)
		if conn.id == id {
			conn.send <- []byte("action!")
		}
	}
	time.Sleep(3 * time.Second)
	c.Operations(id)
}

func (c *Client) Write() {
	defer c.conn.Close()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.Write([]byte{})
				return
			}
			_, err := c.conn.Write([]byte(message))
			if err != nil {
				fmt.Println("Server Write failed,err:", err)
				return
			}
		}
	}
}

func W(conn net.Conn, msg string) bool {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("Server Write failed,err:", err)
		return false
	}
	return true
}

func (c *Client) Ping() {
	for {
		time.Sleep(5 * time.Second)
		msg := "ping ==>" + c.uuid.String()
		d := W(c.conn, msg)
		if d != true {
			manager.unregister <- c
			break
		}
	}
}
