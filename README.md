# Unchain

A lightweight, high-performance proxy server for bypassing network restrictions using VLESS over WebSocket with TLS.

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Supported-blue.svg)](Dockerfile)

## Features

- **VLESS Protocol Support**: Full VLESS over WebSocket with TLS encryption
- **Client Compatible**: Works with v2rayN, v2rayA, Clash, ShadowRocket, and more
- **Lightweight**: Minimal resource footprint with core logic in ~200 lines
- **Production Ready**: Includes traffic metering, health checks, and graceful shutdown
- **Flexible Deployment**: Standalone or integrated with admin servers
- **Memory Efficient**: Optimized goroutine management and resource cleanup
- **High Performance**: Concurrent connections optimized for Go 1.23

## Quick Start

### Prerequisites
- Go 1.23+ or Docker
- Basic proxy/VPN knowledge

### Install from Source

```bash
git clone https://github.com/unchainese/unchain.git
cd unchain
go mod download
cp config.example.standalone.toml config.toml
# Edit config.toml
go run main.go
```

### Docker

```bash
docker build -t unchain .
docker run -p 80:80 \
  -e SUB_ADDRESSES="your-domain.com:443" \
  -e ALLOW_USERS="your-uuid" \
  unchain
```

## Configuration

Unchain uses TOML config or environment variables.

### config.toml

```toml
SubAddresses = 'domain.com:443'
AppPort = '80'
AllowUsers = 'uuid1,uuid2'
LogFile = ''
DebugLevel = 'info'
EnableDataUsageMetering = 'true'
```

### Environment Variables

```bash
APP_PORT=80
SUB_ADDRESSES=domain.com:443
ALLOW_USERS=uuid1,uuid2
```

## Usage

### Endpoints
- `/wsv/{uid}` - VLESS WebSocket endpoint
- `/sub/{uid}` - Subscription URL generator
- `/` - Health check

### Get VLESS URLs
```bash
curl http://localhost:80/sub/your-uuid
```

### Client Setup
Import the generated VLESS URL into your client (v2rayN, Clash, etc.).

## Architecture

```
Client --VLESS/WS/TLS--> Reverse Proxy --WS--> Unchain --TCP/UDP--> Target
```

Unchain runs behind a reverse proxy (Nginx/Cloudflare) handling TLS and WebSocket upgrades.

## Project Structure

```
├── main.go                 # Entry point
├── server/                 # Core server
│   ├── app.go             # HTTP server
│   ├── app_ws_vless.go    # VLESS handler
│   ├── app_ping.go        # Health check
│   └── app_sub.go         # Subscription
├── global/                 # Utilities
│   ├── config.go          # Config management
│   └── logger.go          # Logging
├── schema/                 # Protocols
│   ├── vless.go           # VLESS parser
│   └── trojan.go          # Trojan support
├── client/                 # Client utilities
│   ├── client.go          # SOCKS5 proxy
│   ├── websocket.go       # WS client
│   ├── proxy.go           # Coordination
│   ├── relay_*.go         # Relays
│   ├── socks5_*.go        # SOCKS5
│   └── geo.go             # GeoIP
├── config.example.standalone.toml
└── Dockerfile
```

## Technology Stack

- **Go 1.23**
- **gorilla/websocket** v1.5.3
- **BurntSushi/toml** v1.4.0
- **google/uuid** v1.6.0
- **oschwald/geoip2-golang** v1.11.0
- **sirupsen/logrus** v1.9.3

## Troubleshooting

### Common Issues

1. **Connection Failed**
   - Check server status: `curl http://localhost:80/`
   - Verify firewall and DNS
   - Ensure reverse proxy is configured

2. **WebSocket Errors**
   - Confirm proxy supports WS upgrades
   - Check TLS certificates

3. **Auth Failed**
   - Verify UUID in `AllowUsers`
   - Check admin server if used

4. **High Memory**
   - Monitor goroutines via `/`
   - Restart if >1000 goroutines

### Debug
Set `DebugLevel = 'debug'` and check logs.

## Client Compatibility

- v2rayN (Windows)
- v2rayA (Cross-platform)
- Clash (Cross-platform)
- v2rayNG (Android)
- ShadowRocket (iOS)

## Business Use

Integrate with admin server for user management, traffic metering, and billing.

See [UnchainAdmin](https://github.com/unchainese/unchainadmin) for example.

## Performance

- **RAM**: 512MB min, 1GB+ recommended
- **CPU**: 1 core min, 2+ for production
- **Connections**: Thousands concurrent
- **Memory**: ~20MB base + ~1-2MB/100 connections

## Contributing

1. Fork the repo
2. Create feature branch
3. Commit changes
4. Push and open PR

### Development

```bash
git clone https://github.com/yourusername/unchain.git
cd unchain
go mod download
go test ./...
go build
```

## License

Apache License 2.0 - see [LICENSE](LICENSE)

---

⭐ Star if useful! [Issues](https://github.com/unchainese/unchain/issues) 
