# Testing UDP over SOCKS5 with VLESS Using V2Ray Client

This guide outlines the steps to test UDP functionality over SOCKS5 with VLESS using the V2Ray client.

## Prerequisites
- Ensure you have the necessary files: `udp_echo_svr.go`, `main.go`, `config.json`, and `udpcheck.py` in the `testkit` directory.
- Access to the server `s3.mojotv.cn` for deploying the UDP echo server.

## Steps

1. **Deploy UDP Echo Server**:
   - Copy `testkit/udp_echo_svr.go` to the server `s3.mojotv.cn` using SCP.
   - On the server, run the following command to start the UDP echo server in the background:
     ```
     go run udp_echo_svr.go &
     ```

2. **Start VLESS Proxy Server**:


config.toml
```toml
#toml 格式的文件,建议字符串使用单引号,大小敏感
#

SubAddresses = 'n-us1.libragen.cn:80'# 可以被广域网访问的域名端口,可以是域名也可以是ip,多个地址用逗号分隔
AppPort = '8013' # import same as config.json
RegisterUrl = '' #主控服务器地址,主要作用是控制用户授权和流量计费,可以为空则为个人模式
RegisterToken = 'unchain.people.from.censorship.and.surveillance'# 主控服务器的token
AllowUsers = '6fe57e3f-e618-4873-ba96-a76adec22ccd' # important! same as config.json uuid
LogFile = '' # 
DebugLevel = 'debug' # 日志基本debug, info, warn, error
EnableDataUsageMetering = 'false'
```

   - In your local environment, navigate to the project root and run:
     ```
     go run main.go
     ```
     This starts the VLESS over WebSocket proxy server.

3. **Run V2Ray Client**:
   - Follow the installation guide for V2Ray: https://www.v2fly.org/guide/install.html
   - Change to the `testkit` directory and start the V2Ray client using the provided configuration:
     ```
     cd testkit && v2ray run
     ```
     The client will use `testkit/config.json` for configuration.

4. **Perform UDP Check**:
   - From the `testkit` directory, run the UDP check script:
     ```
     cd testkit && python3 udpcheck.py
     ```
     This will test UDP connectivity through the SOCKS5 proxy.

## Notes
- Ensure all services are running and accessible before proceeding to the next step.
- Monitor the terminal output for any errors or confirmations during the process.