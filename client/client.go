package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"sync"
	"time"
)

type ClientLocalSocks5Server struct {
	AddrSocks5 string
	Timeout    time.Duration
	proxy      *Proxy
}

func NewClientLocalSocks5Server(addr string) (*ClientLocalSocks5Server, error) {
	return &ClientLocalSocks5Server{
		AddrSocks5: addr,
		Timeout:    5 * time.Minute,
	}, nil
}

func (ss *ClientLocalSocks5Server) fetchActiveProxy() {

	ss.proxy = &Proxy{}
}

func (ss *ClientLocalSocks5Server) Run(ctx context.Context) {
	ss.fetchActiveProxy()

	listener, err := net.Listen("tcp", ss.AddrSocks5)
	if err != nil {
		listener, err = net.Listen("tcp4", "127.0.0.1:0")
	}
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", ss.AddrSocks5, err)
	}
	ss.AddrSocks5 = listener.Addr().String()
	slog.Info("socks5 server listening on", "addr", ss.AddrSocks5)

	defer listener.Close()
	log.Println("SOCKS5 server listening on: " + ss.AddrSocks5)
	//proxySettingOn(ss.AddrSocks5)
	//defer proxySettingOff()

	// Channel to receive new connections
	connCh := make(chan net.Conn, 1)
	// Channel to signal accept goroutine to stop
	done := make(chan struct{})
	defer close(done)

	// Start accept goroutine
	go func() {
		defer close(connCh)
		for {
			select {
			case <-done:
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					// Check if the error is due to listener being closed
					select {
					case <-done:
						return // Expected closure
					default:
						log.Printf("Failed to accept connection: %v", err)
						continue
					}
				}
				select {
				case connCh <- conn:
				case <-done:
					conn.Close() // Close connection if we can't send it
					return
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("socks5 server exit")
			return
		case conn, ok := <-connCh:
			if !ok {
				// Connection channel closed, exit
				return
			}
			go ss.handleConnection(ctx, conn)
		}
	}
}

func (ss *ClientLocalSocks5Server) socks5HandShake(conn net.Conn) error {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return fmt.Errorf("failed to read version and nmethods: %w", err)
	}
	if buf[0] != socks5Version {
		return fmt.Errorf("socks5 only. unsupported SOCKS version: %d", buf[0])
	}

	// Read the supported authentication methods
	nMethods := int(buf[1])
	nMethodsData := make([]byte, nMethods)
	if _, err := io.ReadFull(conn, nMethodsData); err != nil {
		return fmt.Errorf("failed to read methods: %w", err)
	}

	// Select no authentication (0x00)
	if _, err := conn.Write([]byte{socks5Version, 0x00}); err != nil {
		return fmt.Errorf("failed to write method selection: %w", err)
	}
	return nil
}

func (ss *ClientLocalSocks5Server) socks5Request(conn net.Conn) (*Socks5Request, error) {
	buf := make([]byte, 8<<10)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read request: %w", err)
	}
	data := buf[:n]
	if len(data) < 4 {
		return nil, fmt.Errorf("request too short")
	}
	return parseSocks5Request(data)
}

func (ss *ClientLocalSocks5Server) handleConnection(outerCtx context.Context, conn net.Conn) {
	defer conn.Close() // the outer for loop is not suitable for defer, so defer close here
	ctx, cf := context.WithTimeout(outerCtx, ss.Timeout)
	defer cf()

	err := ss.socks5HandShake(conn)
	if err != nil {
		slog.Error("failed to handshake", "err", err.Error())
		socks5Response(conn, net.IPv4zero, 0, socks5ReplyFail)
		return
	}
	req, err := ss.socks5Request(conn)
	if err != nil {
		slog.Error("failed to parse socks5 request", "err", err.Error())
		socks5Response(conn, net.IPv4zero, 0, socks5ReplyFail)
		return
	}
	req.Logger().Info("remote target")
	if req.socks5Cmd == socks5CmdConnect { //tcp
		relayTcpSvr, err := ss.dispatchRelayTcpServer(ctx, req)
		if err != nil {
			slog.Error("failed to dispatch relay tcp server", "err", err.Error())
			socks5Response(conn, net.IPv4zero, 0, socks5ReplyFail)
			return
		}
		socks5Response(conn, net.IPv4zero, 0, socks5ReplyOkay)
		defer relayTcpSvr.Close()
		ss.pipeTcp(ctx, conn, relayTcpSvr)
		return
	} else if req.socks5Cmd == socks5CmdUdpAssoc {
		udpH, err := NewRelayUdpDirect(conn)
		if err != nil {
			slog.Error("failed to create udp handler", "err", err.Error())
			socks5Response(conn, net.IPv4zero, 0, socks5ReplyFail)
			return
		}

		defer udpH.Close()
		udpH.PipeUdp()
		return
	} else if req.socks5Cmd == socks5CmdBind {
		relayBind(conn, req)
		return
	} else {
		err = fmt.Errorf("unknown command: %d", req.socks5Cmd)
		slog.Error("unknown command", "err", err.Error())
		socks5Response(conn, net.IPv4zero, 0, socks5ReplyFail)
	}
}

func (ss *ClientLocalSocks5Server) shouldGoDirect(req *Socks5Request) (goDirect bool) {

	if req.CountryCode == "CN" || req.CountryCode == "" {
		//empty means geo ip failed or local address
		return true
	}

	return false
}

func (ss *ClientLocalSocks5Server) dispatchRelayTcpServer(ctx context.Context, req *Socks5Request) (io.ReadWriteCloser, error) {
	if ss.shouldGoDirect(req) {
		req.Logger().Info("go direct")
		return NewRelayTcpDirect(req)
	}
	return NewRelayTcpSocks5e(ctx, ss.proxy, req)
}

func (ss *ClientLocalSocks5Server) pipeTcp(ctx context.Context, s5 net.Conn, relayRw io.ReadWriter) {
	// Create cancellable context for proper goroutine cleanup
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // Ensure both goroutines exit when function returns

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		span := slog.With("fn", "ws -> s5")
		defer func() {
			span.Debug("wg1 done")
			cancel() // Cancel context if this goroutine exits
			wg.Done()
		}()
		for {
			select {
			case <-ctx.Done():
				span.Info("ctx.Done exit")
				return
			default:
				//ws.SetReadDeadline(time.Now().Add(1 * time.Second))
				buf := make([]byte, 8<<10)
				n, err := relayRw.Read(buf)
				if err != nil {
					span.Error("relay read", "err", err.Error())
					return
				}
				_, err = s5.Write(buf[:n])
				if err != nil {
					span.Error("s5 write", "err", err.Error())
					return
				}
			}
		}
	}()
	go func() { // s5 -> ws
		span := slog.With("fn", "s5 -> ws")
		defer func() {
			span.Debug("wg2 done")
			cancel() // Cancel context if this goroutine exits
			wg.Done()
		}()
		for {
			select {
			case <-ctx.Done():
				span.Debug("ctx.Done exit")
				return
			default:
				buf := make([]byte, 8<<10)
				//s5.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
				n, err := s5.Read(buf)
				if errors.Is(err, io.EOF) {
					slog.Info("s5 read EOF")
					return
				}
				if err != nil {
					et := fmt.Sprintf("%T", err)
					span.With("errType", et).Error("s5 read", "err", err.Error())
					return
				}
				//ws.SetWriteDeadline(time.Now().Add(1 * time.Second))
				n, err = relayRw.Write(buf[:n])
				if err != nil {
					span.Error("relay write", "err", err.Error())
					return
				}
			}
		}
	}()
	wg.Wait()
	slog.Debug("2 goroutines is Done")
}
