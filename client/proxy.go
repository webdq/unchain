package client

import (
	"fmt"
	"strings"
)

type Proxy struct {
	Name     string
	Protocol string //ws,wss,http2,tls,http3
	Host     string
	Uri      string
	Sni      string
	Version  string // one socks5

	UserID    string
	Password  string
	TrafficKb int64
	SpeedMs   int64
	Status    string //active, inactive
}

func (p *Proxy) RelayURL() string {
	switch p.Protocol {
	case "ws":
		return fmt.Sprintf("ws://%s/%s", p.Host, strings.TrimPrefix(p.Uri, "/"))
	case "wss":
		return fmt.Sprintf("wss://%s/%s", p.Host, strings.TrimPrefix(p.Uri, "/"))
	case "tcp+tls":
		return fmt.Sprintf("tcp-tls://%s/%s", p.Host, strings.TrimPrefix(p.Uri, "/"))
	default:
		return ""
	}
}

func (p *Proxy) EnableEarlyData() bool {
	//https://xtls.github.io/config/transports/websocket.html#websocketobject
	if strings.Contains(p.Uri, "?ed=") {
		return true
	}
	return false
}
