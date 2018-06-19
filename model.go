package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type Fimp struct {
	Adapter string `json:"adapter"`
	Address string `json:"address"`
	Group   string `json:"group"`
}

type Device struct {
	Fimp   Fimp                   `json:"_fimp"`
	Client VClient                `json:"client"`
	Param  map[string]interface{} `json:"param"`
	Room   int                    `json:"room"`
}

type House struct {
	Mode string    `json:"mode"`
	Time time.Time `json:"time"`
}

type Room struct {
	ID     int         `json:"id"`
	Param  interface{} `json:"param"`
	Client VClient     `json:"client"`
	Type   string      `json:"type"`
}

type Data struct {
	Cmd       string `json:"cmd"`
	Param     Param  `json:"param"`
	RequestID int    `json:"requestId"`
}

type Param struct {
	Components []string `json:"components"`
	Device     []Device `json:"device,omitempty"`
	Room       []Room   `json:"room,omitempty"`
	House      House    `json:"house,omitempty"`
}

type Msg struct {
	Type string `json:"type"`
	Src  string `json:"src"`
	Dst  string `json:"dst"`
	Data Data   `json:"data"`
}

type VinculumMsg struct {
	Ver string `json:"ver"`
	Msg Msg    `json:"msg"`
}

type VClient struct {
	Host            string
	Client          *websocket.Conn
	Msg             chan []byte
	IsRunning       bool
	RunningRequests map[int]chan VinculumMsg
	Subscribers     []chan VinculumMsg
}
