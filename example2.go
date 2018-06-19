package main

import (
	"errors"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

func NewClient(host string) *VClient {
	vc := VClient{Host: host, IsRunning: false}
	log.Info(vc)
	return &vc
}

func (vc *VClient) Connect() error {
	defer func() {
		if r := recover(); r != nil {
			log.Error("<Client> Process CRASHED with error : ", r)
		}
	}()
	vc.RunningRequests = make(map[int]chan VinculumMsg)
	u := url.URL{Scheme: "ws", Host: vc.Host, Path: "/ws"}
	log.Infof("<Client> Connecting to %s", u.String())
	var err error
	vc.Client, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error("<Client> Dial error", err)
		vc.IsRunning = false
		return err
	}
	log.Info("client connected")

	go func() {
		log.Info("in recv func")
		defer func() {
			if r := recover(); r != nil {
				log.Error("<VincClient> Process CRASHED with error : ", r)
				vc.IsRunning = false
			}
		}()

		for {
			vincMsg := VinculumMsg{}

			err := vc.Client.ReadJSON(&vincMsg)
			log.Info("vincMsg : ")
			log.Info(vincMsg)
			if err != nil {
				continue
			}
			if vincMsg.Msg.Type == "response" {
				log.Info("msg type is response")
				for k, vchan := range vc.RunningRequests {
					if k == vincMsg.Msg.Data.RequestID {
						vchan <- vincMsg
					}
				}
			} else {
				log.Info("in else block")
				str := "len of subscribers : " + strconv.Itoa(len(vc.Subscribers))
				log.Info(str)
				for i := range vc.Subscribers {
					log.Info(i)
					select {
					case vc.Subscribers[i] <- vincMsg:
						log.Info("subscribers <- vincMsg")
					default:
						log.Infof("<VincClient> No listeners on the channel")
					}
				}
			}

			if !vc.IsRunning {
				break
			}
		}
		vc.Client.Close()
	}()

	return nil
}

func (vc *VClient) GetTransmissionNr(components []string) (VinculumMsg, error) {
	log.Info("GetTransmissionNr")

	if !vc.IsRunning {
		log.Info("not running client")
		err := vc.Connect()
		if err != nil {
			return VinculumMsg{}, errors.New("Vinculum is Not connected ")
		}
	}

	reqId := rand.Intn(1000)
	msg := VinculumMsg{Ver: "sevenOfNine", Msg: Msg{Type: "request", Src: "chat", Dst: "vinculum", Data: Data{Cmd: "get", RequestID: reqId, Param: Param{Components: components}}}}

	log.Info(msg)
	vc.Client.WriteJSON(msg)
	log.Info(reqId)
	vc.RunningRequests[reqId] = make(chan VinculumMsg)
	select {
	case msg := <-vc.RunningRequests[reqId]:
		delete(vc.RunningRequests, reqId)
		log.Info(msg)
		return msg, nil
	case <-time.After(time.Second * 5):
		log.Info("<VincClient> timeout 5 sec")
	}
	delete(vc.RunningRequests, reqId)

	return VinculumMsg{}, errors.New("Timeout")
}

func main() {
	vc := NewClient("localhost:1989")
	components := []string{"area", "room", "devices"}
	log.Info(vc.GetTransmissionNr(components))
}
