package wsbridge

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Bridge is websocket connection bridge
type Bridge interface {
	Start() error
}

type bridge struct {
	Bridge
	// host name used to listen a websocket connection request
	Host string
	// port number used to listen a websocket connection request
	Port int
	// target websocket server URI
	Server url.URL
	// upgrader instance used to upgrade HTTP request to Websocket
	upgrader *websocket.Upgrader
	// dialer used to connect the target websocet server
	dialer *websocket.Dialer
}

// NewBridge creates a new bridge instance
func NewBridge(config *Config) Bridge {
	var proxy func(*http.Request) (*url.URL, error)
	if config.Proxy != nil {
		proxy = http.ProxyURL(config.Proxy)
	} else {
		proxy = http.ProxyFromEnvironment
	}
	return &bridge{
		Host:   config.Host,
		Port:   config.Port,
		Server: config.Server,
		upgrader: &websocket.Upgrader{
			HandshakeTimeout: config.HandshakeTimeout,
		},
		dialer: &websocket.Dialer{
			Proxy:            proxy,
			HandshakeTimeout: config.HandshakeTimeout,
		},
	}
}

// Start starts listening
func (b *bridge) Start() error {
	http.HandleFunc(b.Server.Path, b.handleWebsocketRequest)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", b.Host, b.Port), nil)
}

func (b *bridge) handleWebsocketRequest(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Connecting to client...")
	cc, err := b.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"r":   r,
			"err": err,
		}).Error("Failed to upgrade")
		return
	}
	defer cc.Close()
	logrus.Info("Client has connected")

	logrus.Info("Connecting to server...")
	cs, _, err := b.dialer.Dial(b.Server.String(), nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"uri": b.Server,
			"err": err,
		}).Error("Failed to connect")
		return
	}
	defer cs.Close()
	logrus.Info("Server has connected")

	var once sync.Once
	done := make(chan struct{})

	// client to server
	go func() {
		defer once.Do(func() { close(done) })
		for {
			if err := bypass(cc, cs); err != nil {
				logrus.Debugln(err)
				return
			}
		}
	}()
	// server to client
	go func() {
		defer once.Do(func() { close(done) })
		for {
			if err := bypass(cs, cc); err != nil {
				logrus.Debugln(err)
				return
			}
		}
	}()
	// Wait
	<-done
	logrus.Info("Terminate bridge connection")
}
