package main

import ("fmt"
	"net"
	"os"
	"log"
	"strconv"
	"strings"
	"bufio"
)

/*
This Project is a chat room based on muti threads with golang
introduction:

run server:  go run chitter.go IP port
run client:  go run chitter.go -c IP port

client operation list:
"<content>" or "all:<content>" : broadcast
"<ID> : <content>" : personal message
"whoami:" : To achieve your ID
"exit:" : exit the server

author : Jijun Jiang
e-mail : jiangjijun1994@gmail.com
 */

type Client struct {
	id int
	inwords chan string
	outwords chan string
	online bool
	con net.Conn
}


type message struct{
	senderid int
	receiverid int
	content string
}

var userchannel chan *Client
var idchannel chan int
var personalchannel chan *message
var broadcastchannel chan *message
var Severchannel chan string

func Clientrecieve(conn *net.TCPConn){
	for {
		var buff [1024]byte
		_,err :=conn.Read(buff[0:])
		checkErr(err, "client read")
		outstr := string(buff[:])
		fmt.Println(outstr)
		if len(outstr) >= 9 && string(buff[:])[0:9] == "you exit!" {
			conn.Close()
			os.Exit(3)
		}

	}
}


func main() {
	if len(os.Args) == 3 && os.Args[1] == "-c" {
		var tcpadder *net.TCPAddr
		tcpadder, _ = net.ResolveTCPAddr("tcp", "localhost:"+os.Args[2])
		conn, _ := net.DialTCP("tcp", nil, tcpadder)
		fmt.Println("connect to NuChitter")
		defer conn.Close()
		go Clientrecieve(conn)
		for {
			var msg string
			readtext := bufio.NewReader(os.Stdin)
			msg, _ = readtext.ReadString('\n')
			msg = strings.Replace(msg, "\n", "", -1)
			b := []byte(msg)
			conn.Write(b)
		}

	} else if len(os.Args) == 2 {

		service := os.Args[1]
		listerner, err := net.Listen("tcp", "localhost:" + service)
		checkErr(err, "listerning")
		userchannel = make(chan *Client)
		idchannel = make(chan int, 2)
		personalchannel = make(chan *message)
		broadcastchannel = make(chan *message)
		Severchannel = make(chan string)
		idchannel <- 1
		var userlist = make([]*Client,0)
		go handlemessage(userlist)
		go serverPlay()
		defer listerner.Close()
		for {
			conn, err := listerner.Accept()
			id := <-idchannel
			checkErr(err, "listen accept")
			fmt.Println("new client id = " + strconv.Itoa(id) + " join")
			newuser := new(Client)
			newuser.id = id
			newuser.inwords = make(chan string)
			newuser.outwords = make(chan string)
			newuser.online = true
			newuser.con = conn
			userchannel <- newuser
			go handleClient(*newuser)
		}

	} else {
		println("pls enter in the formate:\nserver:./chitter <port>\nclient:./chitter <port> -c")
		os.Exit(1)
	}
}

//server read clients' typed wordes
func Read(client Client) {
	for {
		var buff [1024]byte
		_,err :=client.con.Read(buff[0:])
		if (err != nil){
			str := "User " + strconv.Itoa(client.id) + " leave"
			Severchannel <- str
			client.con.Close()
			break
		}
		client.inwords <- string(buff[:])
	}
}

//server write to clients' terminal
func Write(client Client) {
	for {
		stream := <-client.outwords
		client.con.Write([]byte(stream))
	}
}


func serverPlay() {
	for {
		tips :=<- Severchannel
		println(tips)
	}
}

func handlemessage(userlist []*Client){

	for {
		select {
		case newuser := <- userchannel:
			userlist = append(userlist, newuser)
			id := len(userlist) + 1
			idchannel <- id

		case messages := <- personalchannel:

			if messages.receiverid > len(userlist) {
				userlist[messages.senderid - 1].outwords <- "no this id exists, pls try again!"
			} else if !userlist[messages.receiverid - 1].online {
				userlist[messages.senderid - 1].outwords <- "this user is off line, pls try again!"
			} else {
				userlist[messages.receiverid - 1].outwords <- messages.content
			}
			if len(messages.content) >= 9 && messages.content[0:9] == "you exit!" {
				userlist[messages.senderid - 1].online = false
			}
		case broadcast := <- broadcastchannel:
			for i := 0; i < len(userlist); i++ {
				if userlist[i].id != broadcast.senderid{
					userlist[i].outwords <- broadcast.content
				}
			}
		}
	}
}


func handleClient(user Client){
	go Read(user)
	go Write(user)

	for {content := <- user.inwords
		messages := new(message)
		pair := strings.Split(content,":")
		if strings.Index(content,":") == -1{
			content = strconv.Itoa(user.id) + ":" + content[0:len(content) - 2]
			messages.content = content
			messages.senderid = user.id
			broadcastchannel <- messages
		}else if isToPersonal(content){
			messages.senderid = user.id
			recid, err := strconv.Atoi(strings.Split(pair[0], " ")[0])
			str := strconv.Itoa(user.id) + ":" + content[strings.Index(content, ":")+1:]
			messages.content = str
			checkErr(err, "transfer personal message")
			messages.receiverid = recid
			personalchannel <- messages
		}else if len(pair[0]) >=3 && pair[0][0:3] == "all" && strings.Index(content,":") != -1{
			content = strconv.Itoa(user.id) + ":" + content[strings.Index(content, ":")+1:]
			messages.content = content
			messages.senderid = user.id
			broadcastchannel <- messages
		}else if len(pair[0]) >=4 && pair[0][0:4] == "exit" && strings.Index(content,":") != -1{
			messages.senderid = user.id
			messages.receiverid = user.id
			messages.content = "you exit!"
			personalchannel <- messages
		}else if len(pair[0]) >=6 && pair[0][0:6] == "whoami" && strings.Index(content,":") != -1{
			str := "chitter: your id is: " + strconv.Itoa(user.id)
			messages.content = str
			messages.senderid = user.id
			messages.receiverid = user.id
			personalchannel <- messages
		} else {
			words := "unkonwn request! pls try again!"
			user.outwords <- words
		}
	}
}


func isToPersonal(str string) bool{
	pair := strings.Split(str,":")

	_, err := strconv.Atoi(strings.Split(pair[0]," ")[0])
	if err == nil{
		return true
	}else{
		return false
	}
}

func checkErr(err error, str string){
	if err!= nil {
		fmt.Println(str+" error")
		log.Fatal(err)
	}
}