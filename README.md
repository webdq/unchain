# Unchain Proxy Server  

Unchain is a lightweight and easy-to-use proxy server designed to bypass network restrictions, censorship, and surveillance effectively.  


## Key Features  
- **Protocol Support**: Seamlessly handles TCP and UDP (VLESS) packets over WebSocket with TLS/Cloudflare support.  
- **Build Your Own VPN Business**: Provides a robust platform for starting your own VPN services.  
- **Compatibility**: Fully compatible with popular proxy clients like v2rayN or any application supporting the VLESS + WebSocket protocol.  


## How It Works  

Unchain operates as a proxy/VPN server, compatible with popular proxy clients such as v2rayN or any application that supports the VLESS+WebSocket protocol. It accepts traffic from various client applications, including:  

- [v2rayN](https://github.com/2dust/v2rayN)  
- [v2rayA](https://github.com/v2rayA/v2rayA)  
- [Clash](https://github.com/Dreamacro/clash)  
- [v2rayNG](https://github.com/2dust/v2rayNG)  
- [iOS app Shadowrocket](https://apps.apple.com/us/app/shadowrocket/id932747118)

Unchain processes incoming traffic and securely forwards it to the destination server, ensuring both security and efficiency in communication.  

## Unchain Architecture



Unchain is a dead simple VLESS over websocket proxy server.
The core biz logic is only 200 lines of code.  [app_ws_vless.go](/internal/node/app_ws_vless.go).

Unchain server uses a simple architecture that is VLESS over WebSocket (WS) + TLS.


```
             V2rayN,V2rayA,Clash or ShadowRocket                          
                 +------------------+
                 |   VLESS Client   |
                 |   +-----------+  |
                 |   | TLS Layer  | |
                 |   +-----------+  |
                 |   | WebSocket  | |
                 |   +-----------+  |
                 +--------|---------+
                          |
                          | Encrypted VLESS Traffic (wss://)
                          |
           +--------------------------------------+
           |         Internet (TLS Secured)       |
           +--------------------------------------+
                          |
                          |
        +-----------------------------------+
        |        Reverse Proxy Server       |
        | (e.g., Nginx or Cloudflare)       |
        |                                   |
        |   +---------------------------+   |
        |   | HTTPS/TLS Termination     |   |
        |   +---------------------------+   |
        |   | WebSocket Proxy (wss://)  |   |
        |   +---------------------------+   |
        |     Forward to VLESS Server       |
        +------------------|----------------+
                           |
           +--------------------------------+
           |     Unchain       Server       |
           |                                |
           |   +------------------------+   |
           |   | WebSocket Handler      |   |
           |   +------------------------+   |
           |   | VLESS Core Processing  |   |
           |   +------------------------+   |
           |                                |
           |   Forward Traffic to Target    |
           +------------------|-------------+
                              |
                     +-----------------+
                     | Target Server   |
                     | or Destination  |
                     +-----------------+

```



## Usage

### 1. Build from Source

To build from source, follow these steps:

1. Clone the repository and navigate to the `cmd/node` directory:
   ```sh
   cd cmd/node
   ```
2. Copy the example configuration file and customize it:
   ```sh
   cp config.example.standalone.toml config.toml
   ```
3. Run the application:
   ```sh
   go run main.go
   ```

### 2. Deploying on Your Own Ubuntu Server Using GitHub Actions

You can deploy the application on an Ubuntu server using GitHub Actions. Here's how:

1. **Fork the repository** to your GitHub account.
2. **Create an Environment** named `production` in your repository settings.
3. **Add the following SSH connection details** to the Environment Secrets:
   - `EC2_HOST`: The SSH host with port (e.g., `1.1.1.1:20`).
   - `EC2_USER`: The SSH user (e.g., `ubuntu`).
   - `EC2_KEY`: Your SSH private key.

4. **Add your TOML configuration file content** to the Environment Variables:
   - `CONFIG_TOML`: Copy the content of your `config.toml` file, replace all `"` with `'`, and paste it here.

learn more in [.github/workflows/deploy.sh](/.github/workflows/deploy.sh)




### 3. Running the Application

Once the application is running, you will see a VLESS connection schema URL in the standard output. Copy and paste this URL into your V2rayN client.

Congratulations! You now have your self-hosted proxy server up and running.




### 4. (Optional) create your own admin app for Auth and Data-traffic

create an RESTful API for [chain proxy server push](https://github.com/unchainese/unchain/blob/5ece8c39814684a8a54e8e009d7c888e5988a017/internal/node/app.go#L161) :
[Register API example code](https://github.com/unchainese/unchainadmin/blob/035b2232d4262c24ef70b8ad7abb9faebaaecc96/functions/api/nodes.ts#L34)


## Build your own VPN business

Using [the cloudflare page UnchainAdmin](https://github.com/unchainese/unchainadmin) start your own VPN business. 