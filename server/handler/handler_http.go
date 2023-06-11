package handler

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/websocket"
	"go-chess/server/config"
	"go-chess/server/dao"
	"go-chess/server/utils/jwt"
	myjwt "go-chess/server/utils/jwt"
	"log"
	"net/http"
	"time"
)

func HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	// 创建房间
	room := NewRoom()
	go room.Run()

	// 将房间添加到房间列表
	rooms[room.ID] = room

	log.Printf("Room created. ID: %s\n", room.ID)

	// 返回房间ID
	response := map[string]string{"roomID": room.ID}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("JSON marshal error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func HandleJoinRoom(w http.ResponseWriter, r *http.Request) {

	// 加入房间
	//get query params
	query := r.URL.Query()
	roomID := query.Get("roomID")

	//header获取token
	token := r.Header.Get("Authorization")
	if token == "" {
		log.Println("Token not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//log.Println("get the Token:", token)
	playerId, err := jwtAuth(token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//log.Println("get the playerId:", playerId)

	//get username from mysql

	username, e := dao.GetUsername(playerId)
	if e != nil {
		log.Println(e)
		return
	}
	//log.Println("get the username:", username)

	room, ok := rooms[roomID]
	if !ok {
		log.Printf("Room %s not found\n", roomID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if room.WhitePlayer != nil && room.BlackPlayer != nil {
		log.Printf("Room %s is full\n", roomID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn, playerId, username, room)
	//log.Println("New client joined username: ", username)
	room.registry <- client

	//TODO: 这里有并发问题，需要加锁
	time.Sleep(1 * time.Second)
	if room.WhitePlayer != nil && room.BlackPlayer != nil {
		runGame(room)
	}
}

func runGame(r *Room) {
	log.Printf("Game started. Room ID: %s\n", r.ID)

	msg := &Message{
		Type: MessageTypeGameStart,
		GameStart: gameStart{
			Board:   r.Board,
			Content: "你是白方",
			Turn:    r.turn,
		},
	}
	r.sendToWhite <- msg

	msg = &Message{
		Type: MessageTypeGameStart,
		GameStart: gameStart{
			Board:   r.Board,
			Content: "你是黑方",
			Turn:    r.turn,
		},
	}
	r.sendToBlack <- msg
}

// HandleRegistry 用户注册
func HandleRegistry(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	log.Println("get username and password", username, password)
	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var u dao.User
	u.Username = username
	u.Password = password

	u.UserID = generateUid()
	//encode password
	h := md5.New()
	h.Write([]byte(u.Password))
	u.Password = base64.StdEncoding.EncodeToString(h.Sum(nil))

	if err := dao.Db.Create(&u).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var u dao.User
	u.Username = username
	u.Password = password
	//encode password
	h := md5.New()
	h.Write([]byte(u.Password))
	u.Password = base64.StdEncoding.EncodeToString(h.Sum(nil))

	//get password from db
	var user dao.User
	if err := dao.Db.Where("username = ?", u.Username).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//log.Println("get user from db:", user)

	if user.Password != u.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwtConfig := config.GlobalConfig.Jwt
	j := myjwt.NewJWT(&myjwt.Config{
		SecretKey:  jwtConfig.SecretKey,
		ExpireTime: jwtConfig.ExpiresTime,
		Issuer:     jwtConfig.Issuer,
	})

	claims := j.CreateClaims(&jwt.BaseClaims{
		Id:         user.UserID,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})
	tokenString, err := j.GenerateToken(&claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//set token to response header
	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusOK)
}

func jwtAuth(token string) (int64, error) {
	//parse token
	jwtConfig := config.GlobalConfig.Jwt
	j := myjwt.NewJWT(&myjwt.Config{
		SecretKey:  jwtConfig.SecretKey,
		ExpireTime: jwtConfig.ExpiresTime,
		Issuer:     jwtConfig.Issuer,
	})
	mc, err := j.ParseToken(token)
	if err != nil {
		return 0, err
	}

	return mc.BaseClaims.Id, nil
}

func generateUid() int64 {
	sf, err := snowflake.NewNode(config.GlobalConfig.Snowflake.MachineId)
	if err != nil {
		panic(err)
	}
	return sf.Generate().Int64()
}
