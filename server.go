package speedtest

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/r3inbowari/common"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	PingCount                 = 20
	InitialDelta        int64 = 6564
	DefaultDownloadSize       = 25000000 * 5 // nByte
)

const (
	WebSocketSecureSchema = "wss"
	WebSocketSchema       = "ws"
	HttpSecureSchema      = "https"
	HttpSchema            = "http"
	WebSocketPath         = "/ws"
	DownloadPath          = "/download?nocache=%s&size=%d&guid=%s"
	UploadPath            = "/upload?nocache=%s&guid=%s"
)

var (
	MessageHi           = []byte{0x48, 0x49}
	MessageGetIP        = []byte{0x47, 0x45, 0x54, 0x49, 0x50}
	MessageCapabilities = []byte{0x43, 0x41, 0x50, 0x41, 0x42, 0x49, 0x4c, 0x49, 0x54, 0x49, 0x45, 0x53}
	MessagePing         = []byte{0x50, 0x49, 0x4e, 0x47, 0x20}
)

type Server struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Lat             string `json:"lat"`
	Lon             string `json:"lon"`
	Distance        int    `json:"distance"`
	Country         string `json:"country"`
	CC              string `json:"cc"`
	Url             string `json:"url"`
	Host            string `json:"host"`
	Sponsor         string `json:"sponsor"`
	Preferred       int    `json:"preferred"`
	HttpsFunctional int    `json:"https_functional"`
	ForcePingSelect int    `json:"force_ping_select"`

	Client *WebsocketClient
	Result Result
}

type Result struct {
	LocalRTT         []time.Duration // a round-trip time as observed locally.
	RemoteTimestamps []string        // response timestamp sequence. means: delta(n+1) - delta(n) = prev-trip + next-round.
	DownloadSize     int64           // actual download size
	DownloadRate     float64         // Mbps. Byte * 1000 / 1024 / 1024 * 8 (Notice that the rawRate unit is Byte/ms).
	UploadSize       int64           // actual upload size
	UploadRate       float64         // Mbps. same as above.
}

func (r *Result) String() string {
	var ret string
	if len(r.LocalRTT) > 0 {
		min := r.LocalRTT[0]
		for _, v := range r.LocalRTT {
			if v < min {
				min = v
			}
		}
		ret += fmt.Sprintf("minimum latency: %.2fms\n", float64(min.Microseconds())/1000)
	}
	ret += fmt.Sprintf("download size: %d Bytes... rate: %.2fMbps\n", r.DownloadSize, r.DownloadRate)
	ret += fmt.Sprintf("upload size: %d Bytes... rate: %.2fMbps\n", r.UploadSize, r.UploadRate)
	return ret
}

func (s *Server) Connect() {
	s.Client = NewWebsocketClient(s.ParseWebsocketUrl(), time.Minute)
	err := s.Client.Dial()
	if err != nil {
		println(err.Error())
		return
	}
	println("connected")
}

func (s *Server) Hi() string {
	s.Client.Conn.WriteMessage(websocket.TextMessage, []byte(MessageHi))
	message, _, err := s.Client.Conn.ReadMessage()
	if err != nil {
		return ""
	}
	return string(rune(message))
}

func (s *Server) Read() (string, error) {
	_, message, err := s.Client.Conn.ReadMessage()
	return string(message), err
}

func (s *Server) Write(msg []byte) error {
	return s.Client.Conn.WriteMessage(websocket.TextMessage, msg)
}

func (s *Server) Ping() (err error) {
	var pong string
	s.Result.LocalRTT = make([]time.Duration, PingCount)
	s.Result.RemoteTimestamps = make([]string, PingCount)
	var delta = InitialDelta
	for i := 0; i < PingCount; i++ {
		startTime := time.Now()
		if err = s.Write(strconv.AppendInt(MessagePing, delta, 10)); err != nil {
			break
		}
		if pong, err = s.Read(); err != nil {
			break
		}
		// parse remote prev-rtt/2 + next-rtt/2 here
		s.Result.RemoteTimestamps[i] = pong[5:18]
		s.Result.LocalRTT[i] = time.Since(startTime)
		delta += s.Result.LocalRTT[i].Milliseconds()
	}
	return err
}

// DownloadTest down-link rate test with server
// host like(https://1010b.hkspeedtest.com.prod.hosts.ooklaserver.net:8080)
// path /download
// param nocache {uuid}
// param guid {uuid}
// param size (default=25000000)
func (s *Server) DownloadTest(ctx context.Context) error {
	downloadUrl, err := url.QueryUnescape(s.ParseDownloadUrl(DefaultDownloadSize).String())
	if err != nil {
		return err
	}
	client := new(http.Client)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadUrl, nil)
	if err != nil {
		println("err " + err.Error())
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		println("err " + err.Error())
		return err
	}

	defer func() { _ = resp.Body.Close() }()
	startTime := time.Now()
	n, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		println(err.Error())
	}
	d := time.Since(startTime)
	s.Result.DownloadSize = n
	s.Result.DownloadRate = float64(n) / float64(d.Milliseconds()) * 1000 / 1024 / 1024 * 8
	return err
}

// UploadTest up-link rate test with server
// host like(https://1010b.hkspeedtest.com.prod.hosts.ooklaserver.net:8080)
// path /upload
// param nocache {uuid}
// param guid {uuid}
// param size (default=25000000)
// we can open a websocket before upload with a same guid,
// so that we can get the upload detail send by remote server
func (s *Server) UploadTest(ctx context.Context) error {
	client := new(http.Client)
	uploadUrl, err := url.QueryUnescape(s.ParseUploadUrl().String())
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Add("content", strings.Repeat("3478789494", DefaultDownloadSize/10))
	reader := strings.NewReader(v.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadUrl, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	startTime := time.Now()
	_, err = client.Do(req)
	d := time.Since(startTime)
	s.Result.UploadSize = reader.Size() - int64(reader.Len())
	s.Result.UploadRate = float64(s.Result.UploadSize) / float64(d.Milliseconds()) * 1000 / 1024 / 1024 * 8
	return err
}

func (s *Server) ParseWebsocketUrl() *url.URL {
	if s.HttpsFunctional > 0 {
		return &url.URL{Scheme: WebSocketSecureSchema, Host: s.Host, Path: WebSocketPath}
	}
	return &url.URL{Scheme: WebSocketSchema, Host: s.Host, Path: WebSocketPath}
}

func (s *Server) ParseDownloadUrl(nByte int) *url.URL {
	if nByte == 0 || nByte == -1 {
		nByte = DefaultDownloadSize
	}
	nocache := common.CreateUUID()
	guid := common.CreateUUID()
	if s.HttpsFunctional > 0 {
		return &url.URL{Scheme: HttpSecureSchema, Host: s.Host, Path: fmt.Sprintf(DownloadPath, nocache, nByte, guid)}
	}
	return &url.URL{Scheme: HttpSchema, Host: s.Host, Path: fmt.Sprintf(DownloadPath, nocache, nByte, guid)}
}

// ParseUploadUrl default 25kB
func (s *Server) ParseUploadUrl() *url.URL {
	nocache := common.CreateUUID()
	guid := common.CreateUUID()
	if s.HttpsFunctional > 0 {
		return &url.URL{Scheme: HttpSecureSchema, Host: s.Host, Path: fmt.Sprintf(UploadPath, nocache, guid)}
	}
	return &url.URL{Scheme: HttpSchema, Host: s.Host, Path: fmt.Sprintf(UploadPath, nocache, guid)}
}
