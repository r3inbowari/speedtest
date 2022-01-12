package speedtest

import (
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"log"
	"net/url"
	"time"
)

type WebsocketClient struct {
	Url           *url.URL
	Conn          *websocket.Conn
	SendMsg       chan string
	RecvMsg       chan string
	DetectionTime time.Duration
	Ctx           context.Context
	Cancel        context.CancelFunc
}

func NewWebsocketClient(url *url.URL, dt time.Duration) *WebsocketClient {
	var sender = make(chan string, 10)
	var receiver = make(chan string, 10)
	var conn *websocket.Conn
	return &WebsocketClient{
		Url:           url,
		Conn:          conn,
		SendMsg:       sender,
		RecvMsg:       receiver,
		DetectionTime: dt,
	}
}

func (wsc *WebsocketClient) dial() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	wsc.Ctx = ctx
	wsc.Cancel = cancel
	wsc.Conn, _, err = websocket.DefaultDialer.DialContext(ctx, wsc.Url.String(), nil) // timeout?
	if err != nil {
		log.Printf("connect failed %s", err.Error())

	}
	log.Printf("connected to %s ", wsc.Url.String())
}

func (wsc *WebsocketClient) Dial() error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	wsc.Ctx = ctx
	wsc.Cancel = cancel
	wsc.Conn, _, err = websocket.DefaultDialer.DialContext(ctx, wsc.Url.String(), nil) // timeout?
	return err
}

func (wsc *WebsocketClient) sendMsgThread() {
	go func() {
		for {
			select {
			case <-wsc.Ctx.Done():
				log.Println("send close")
				return
			case msg := <-wsc.SendMsg:
				err := wsc.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					//log.Println("send:", err)
					wsc.Cancel()
					return
				}
			}
		}
	}()
}

func (wsc *WebsocketClient) readMsgThread() {
	go func() {
		for {
			if wsc.Conn != nil {
				_, message, err := wsc.Conn.ReadMessage()
				if err != nil {
					log.Println("read close")
					wsc.Cancel()
					return
				}
				wsc.RecvMsg <- string(message)
			}
		}
	}()
}

func (wsc *WebsocketClient) Start() {
	fmt.Println(" ")
	//startLoop:
	wsc.dial()
	wsc.sendMsgThread()
	wsc.readMsgThread()
	for {
		select {
		case <-wsc.Ctx.Done():
			log.Printf("disconnected -> %s\n", wsc.Conn.RemoteAddr())
			//goto startLoop
			fmt.Println(" ")
			wsc.dial()
			wsc.sendMsgThread()
			wsc.readMsgThread()
		}
	}
}
