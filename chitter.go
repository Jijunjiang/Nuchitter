package main

import ("fmt"
	"net"
	"os"
	"log"
	"strconv"
	"strings"
	"bufio"
)


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


func Clientrecieve(conn *net.TCPConn){
	for {
		var buff [1024]byte
		_,err :=conn.Read(buff[0:])
		checkErr(err, "client read")
		fmt.Println(string(buff[:]))

	}
}


func main() {
	if len(os.Args) == 3 && os.Args[2] == "-c" {
		var tcpadder *net.TCPAddr
		tcpadder, _ = net.ResolveTCPAddr("tcp", "localhost:"+os.Args[1])
		conn, _ := net.DialTCP("tcp", nil, tcpadder)
		fmt.Println("connect to NuChitter")
		go Clientrecieve(conn)

		for {
			var msg string
			readtext := bufio.NewReader(os.Stdin)
			msg, _ = readtext.ReadString('\n')
			msg = strings.Replace(msg, "\n", "", -1)

			if msg == "quit" {
				break
			}
			b := []byte(msg)
			conn.Write(b)
		}

	} else if len(os.Args) == 2 {

		service := os.Args[1]
		listerner, err := net.Listen("tcp", "localhost:"+service)
		checkErr(err, "listerning")
		userchannel = make(chan *Client)
		idchannel = make(chan int, 2)
		personalchannel = make(chan *message)
		broadcastchannel = make(chan *message)
		idchannel <- 1
		var userlist = make([]*Client,0)
		go handlemessage(userlist)
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
			println("User " + strconv.Itoa(client.id) + " leave")
			client.online = false
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

func handlemessage(userlist []*Client){

	for {
		select {
		case newuser := <- userchannel:
			userlist = append(userlist, newuser)
			id := len(userlist) + 1
			idchannel <- id

		case messages := <- personalchannel:
			if messages.receiverid <= len(userlist) {
				userlist[messages.receiverid-1].outwords <- messages.content
			}else{userlist[messages.senderid - 1].outwords <- "no this id exists, pls try again!"}
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
		} else if len(pair[0]) < 3 {
			words := "unkonwn request! pls try again!"
			user.outwords <- words
		} else if pair[0][0:3] == "all"{
			content = strconv.Itoa(user.id) + ":" + content[strings.Index(content,":") + 1:]
			messages.content = content
			messages.senderid = user.id
			broadcastchannel <- messages
		} else if len(pair[0]) <6 && (len(pair[0]) >= 3 && pair[0][0:3] != "all") {
			words := "unkonwn request! pls try again!"
			user.outwords <- words
		} else if pair[0][0:6] == "whoami" && strings.Index(content,":") != -1{
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