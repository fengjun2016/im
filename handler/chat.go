package handler

import (
	"github.com/julienschmidt/httprouter"
	"gopkg.in/fatih/set.v0"
	"net/http"
)

const (
	CMD_SINGLE_MSG = 10
	CMD_ROOM_MSG   = 11
	CMD_HEART      = 0
)

type Message struct {
	Id      int64  `json:"id,omitempty" form:"id"`           //消息ID
	Userid  int64  `json:"userid,omitempty" form:"userid"`   //谁发的
	Cmd     int    `json:"cmd,omitempty" form:"cmd"`         //群聊还是私聊
	Dstid   int64  `json:"dstid,omitempty" form:"dstid"`     //对端用户ID/群ID
	Media   int    `json:"media,omitempty" form:"media"`     //消息按照什么样式展示
	Content string `json:"content,omitempty" form:"content"` //消息的内容
	Pic     string `json:"pic,omitempty" form:"pic"`         //预览图片
	Url     string `json:"url,omitempty" form:"url"`         //服务的URL
	Memo    string `json:"memo,omitempty" form:"memo"`       //简单描述
	Amount  int    `json:"amount,omitempty" form:"amount"`   //其他和数字相关的
}

/**
消息发送结构体
1、MEDIA_TYPE_TEXT
{id:1,userid:2,dstid:3,cmd:10,media:1,content:"hello"}
2、MEDIA_TYPE_News
{id:1,userid:2,dstid:3,cmd:10,media:2,content:"标题",pic:"http://www.baidu.com/a/log,jpg",url:"http://www.a,com/dsturl","memo":"这是描述"}
3、MEDIA_TYPE_VOICE，amount单位秒
{id:1,userid:2,dstid:3,cmd:10,media:3,url:"http://www.a,com/dsturl.mp3",anount:40}
4、MEDIA_TYPE_IMG
{id:1,userid:2,dstid:3,cmd:10,media:4,url:"http://www.baidu.com/a/log,jpg"}
5、MEDIA_TYPE_REDPACKAGR //红包amount 单位分
{id:1,userid:2,dstid:3,cmd:10,media:5,url:"http://www.baidu.com/a/b/c/redpackageaddress?id=100000","amount":300,"memo":"恭喜发财"}
6、MEDIA_TYPE_EMOJ 6
{id:1,userid:2,dstid:3,cmd:10,media:6,"content":"cry"}
7、MEDIA_TYPE_Link 6
{id:1,userid:2,dstid:3,cmd:10,media:7,"url":"http://www.a,com/dsturl.html"}

7、MEDIA_TYPE_Link 6
{id:1,userid:2,dstid:3,cmd:10,media:7,"url":"http://www.a,com/dsturl.html"}

8、MEDIA_TYPE_VIDEO 8
{id:1,userid:2,dstid:3,cmd:10,media:8,pic:"http://www.baidu.com/a/log,jpg",url:"http://www.a,com/a.mp4"}

9、MEDIA_TYPE_CONTACT 9
{id:1,userid:2,dstid:3,cmd:10,media:9,"content":"10086","pic":"http://www.baidu.com/a/avatar,jpg","memo":"胡大力"}

*/

//本核心在于形成userid和Node的映射关系
type Node struct {
	Conn *websocket.Conn
	//并行转串行
	DataQueue chan []byte
	GroupSets set.Interface
}

//映射关系表
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

//读写锁
var rwlocker sync.RWMutex

func Chat(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	//todo 检验接入是否合法
	query := req.URL.Query()
	userId := query.Get("id")
	token := query.Get("token")

	isvalid := checkToken(userId, token)
	//如果isvalid为true
	//如果isvalid为false

	// 升级http协议至websocket协议
	var upgrader = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{c.Request.Header.Get("Sec-WebSocket-Protocol")},
	}
	//完成tcp websocket握手协议
	wsConn, err = upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Println("connect ws failed.", err)
		return
	}

	//todo 获得conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	//todo 获取用户全部群id
	contact := Model.Contact{}
	if _, commIds, _, err := contact.SearchCommunityIds(userId); err != nil {
		logrus.Println("init get community userid s info failed.", err.Error())
		ResponseJson(rw, "")
		return
	}
	//刷新 group set
	for _, v := range commIds {
		node.GroupSets.Add(v)
	}

	//todo userid和node形成绑定关系
	rwlocker.Lock()
	clientMap[userId] = node
	rwlocker.Unlock()
	//todo 完成发送逻辑
	go sendproc(node)
	//todo 完成接收逻辑
	go recvproc(node)

	sendMsg(userId, []byte("hello ,world")) //发送心跳包
}

//发送协程
func sendproc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				logrus.Println("send message failed.", err.Error())
				return
			}
		}
	}
}

//接收协程
func recproc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			logrus.Println("recv message failed.", err.Error())
			return
		}
		//todo 对data进一步处理
		dispatch(data)
		//把消息广播到局域网
		broadMsg(data)
		fmt.Printf("recv<=%s", data)
	}
}

//todo 发送消息
func sendMsg(userId string, msg []byte) {
	rwlocker.RLock()
	node, ok := clientMap[userId]
	rwlocker.RULock()
	if ok {
		node.DataQueue <- msg
	}
}

//检测token是否有效
func checkToken(userId, token string) bool {
	//从数据库里面查询并进行比对
	user := model.User{}
	user.ID = userId
	if err := user.Get(); err != nil {
		logrus.Println("checkToken failed.", err)
		return false
	}
	return user.Token == token
}

//消息分发 dispatch
func dispatch(data []byte) {
	//todo 解析data为message
	msg := Message{}
	if err := json.Unmarshal(data, &msg); err != nil {
		logrus.Println("json unmarshal data to msg failed.", err.Error())
		return
	}
	//todo 根据message的cmd属性做相应的处理
	switch msg.cmd {
	//单聊
	case CMD_SINGLE_MSG:
		sendMsg(msg.Dstid, data)
	//群聊
	case CMD_ROOM_MSG:
		//todo 群聊转发逻辑
		for _, v := range clientMap {
			if v.GroupSets.Has(msg.Dstid) {
				v.DataQueue <- data
			}
		}
	//心跳
	case CMD_HEART:
		//todo 一般啥都不干 只是为了保持tcp websocket的连接长时间不断开
	}
}

//用来存放发送的要广播的数据
var udpsendchan chan []byte = make(chan []byte, 1024)

//使用基于udp协议的广播到不同服务器的节点上
func broadMsg(data []byte) error {
	udpsendchan <- data
}

//todo 完成udp数据的发送协程
func udpsendproc() {
	logrus.Println("start udpsendproc 开启udp发送数据的协程")
	//todo 使用udp协议拨号连接
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 3000,
	})
	defer conn.Close()
	if err != nil {
		logrus.Println(err.Error())
		return
	}
	//todo 通过得到的conn连接发送消息
	//for select 阻塞等待要发送的数据
	for {
		select {
		case data := <-udpsendchan:
			_, err = conn.Write(data)
			if err != nil {
				logrus.Println("stop idp send proc")
				return
			}
		}
	}
}

//项目开启运行udp发送和接收消息通知的协程运行
func init() {
	go udpsendproc()
	go udprecvproc()
}

//todo 完成udp数据的接收协程
func udprecvproc() {
	logrus.Println("start udp recv proc goroutine")
	//todo 监听udp广播端口
	conn, err := net.ListenUdP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer conn.Close()
	if err != nil {
		logrus.Println("stop udp recv proc goroutine")
		return //退出本次协程的运行
	}
	//todo 处理端口发送过来的数据
	for {
		var buf [512]byte //存储从udp接收到的广播协议的数据
		intn, err := conn.Read(buf[0:])
		if err != nil {
			logrus.Println(err.Error())
			return
		}
		//直接调用之前的消息分发函数 进行消息分发处理
		dispatch(buf[0:n])
	}
	logrus.Println("return udprecvproc")
}
