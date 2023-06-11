package main

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell"
	ws "github.com/gorilla/websocket"
	"go-chess/client/client"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	Token  string
	screen tcell.Screen
)

const (
	ModeChat = 0
	ModeCmd  = 1
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	file, err := os.Create("log.txt")
	if err != nil {
		log.Fatal("无法创建日志文件:", err)
	}
	defer file.Close()
	// 设置日志输出位置为文件
	log.SetOutput(file)
	for {
		fmt.Println("[1]创建房间 [2]加入房间 [3]注册 [4]登录 [5]退出")
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			fmt.Println("创建房间")
			newRoom()
		case 2:
			var roomId string
			fmt.Println("请输入房间号:")
			fmt.Scanln(&roomId)
			// post to server
			// 设置房间ID
			fmt.Println("加入房间")
			joinRoom(roomId)
		case 3:
			fmt.Println("注册")
			registry()
		case 4:
			fmt.Println("登录")
			login()
		case 5:
			fmt.Println("退出")
			return
		default:
			fmt.Println("无效的选项")
		}
	}
}

func joinRoom(roomId string) {

	url := fmt.Sprintf("ws://127.0.0.1:8080/join?&roomID=%s", roomId)
	log.Println("连接到服务器:", url)
	//with token
	headers := http.Header{}
	if Token == "" {
		fmt.Println("请先登录")
		return
	}
	headers.Add("Authorization", Token)
	dialer := ws.Dialer{}
	conn, _, err := dialer.Dial(url, headers)
	if err != nil {
		log.Fatalln(err)
	}
	screen, _ = tcell.NewScreen()

	defer screen.Fini()
	if err := screen.Init(); err != nil {
		panic(err)
	}
	screen.Clear()

	c := client.NewClient(roomId, conn, screen)
	mode := ModeCmd
	go func() {
		for {
			screen.Clear()
			switch mode {
			case ModeChat:
				writeStringToScreen(1, 12, "聊天: ")
			case ModeCmd:
				writeStringToScreen(1, 12, "命令:")
			}
			writeStringToScreen(8, 12, c.Input)

			switch c.Turn {
			case client.TurnWhite:
				writeStringToScreen(1, 11, "白方回合")
			case client.TurnBlack:
				writeStringToScreen(1, 11, "黑方回合")
			case client.TurnNone:
				writeStringToScreen(1, 11, "等待玩家加入")
			}

			c.Show()
			screen.Show()
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {
		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyESC {
					switch mode {
					case ModeChat:
						mode = ModeCmd
					case ModeCmd:
						mode = ModeChat
					}
				} else if ev.Key() == tcell.KeyRune {
					c.Input += string(ev.Rune())
				} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
					if len(c.Input) > 0 {
						c.Input = c.Input[:len(c.Input)-1]
					}
				} else if ev.Key() == tcell.KeyEnter {
					switch mode {
					case ModeChat:
						c.ChatChan <- c.Input
						c.Input = ""

					case ModeCmd:
						if client.RegexTypeCastle.MatchString(c.Input) {
							c.CmdChan <- c.Input
							c.Input = ""
						} else if client.RegexTypeMove.MatchString(c.Input) {
							c.CmdChan <- c.Input
							c.Input = ""
						} else if client.RegexTypePromotion.MatchString(c.Input) {
							c.CmdChan <- c.Input
							c.Input = ""
						} else {
							c.Input = "invalid command"
						}
					}
				} else if ev.Key() == tcell.KeyCtrlC {
					c.Exit <- struct{}{}
					return
				}
			}
		}
	}()
	<-c.Exit
}

type RoomCreateResponse struct {
	RoomID string `json:"roomID"`
}

func newRoom() {

	URL := "http://localhost:8080/room" // 设置创建房间的URL

	// 发送POST请求
	resp, err := http.Post(URL, "application/x-www-form-urlencoded", nil)
	if err != nil {
		log.Fatal("POST request failed:", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println("创建房间失败", resp.Status)
	}
	var roomCreateResponse RoomCreateResponse
	err = json.NewDecoder(resp.Body).Decode(&roomCreateResponse)
	if err != nil {
		log.Fatal("JSON decoding failed:", err)
	}
	fmt.Println("创建房间成功 房间号为:", roomCreateResponse.RoomID)

	// 进入房间
	joinRoom(roomCreateResponse.RoomID)
}

func writeStringToScreen(x, y int, s string) {
	for i, v := range s {
		screen.SetContent(x+i, y, v, nil, tcell.StyleDefault)
	}
}

func registry() {
	var username, password string
	fmt.Println("请输入用户名:")
	fmt.Scanln(&username)
	fmt.Println("请输入密码:")
	fmt.Scanln(&password)

	URL := "http://localhost:8080/registry" // 设置创建房间的URL
	// 构造请求的参数
	payload := url.Values{}
	payload.Add("username", username)
	payload.Add("password", password)

	// 发送POST请求
	resp, err := http.PostForm(URL, payload)
	if err != nil {
		log.Fatal("POST request failed:", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("注册失败 ", resp.Status)
	}
	fmt.Println("注册成功")
}

func login() {
	var username, password string
	fmt.Println("请输入用户名:")
	fmt.Scanln(&username)
	fmt.Println("请输入密码:")
	fmt.Scanln(&password)
	URL := "http://localhost:8080/login" // 设置创建房间的URL
	// 构造请求的参数
	payload := url.Values{}
	payload.Add("username", username)
	payload.Add("password", password)

	// 发送POST请求
	resp, err := http.PostForm(URL, payload)
	if err != nil {
		log.Fatal("POST request failed:", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("登录失败 ", resp.Status)
	}
	//set token
	Token = resp.Header.Get("Authorization")
}
