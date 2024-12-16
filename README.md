# Unchain Proxy Server  

Unchain is a lightweight and easy-to-use proxy server designed to bypass network restrictions, censorship, and surveillance effectively.  

## Key Features  
- **Protocol Support**: Unchain proxies TCP+UDP(VLESS) packets over WebSocket (TLS/Cloudflare).  
- **Built your own VPN business**: You can build your own business on it.

## How It Works  
Unchain accepts traffic from client applications like:  
- [v2rayN](https://github.com/2dust/v2rayN)  
- [v2rayA](https://github.com/v2rayA/v2rayA)  
- [Clash](https://github.com/Dreamacro/clash)  
- [v2rayNG](https://github.com/2dust/v2rayNG)  

It processes the incoming traffic and seamlessly forwards it to the destination server, ensuring secure and efficient communication.  



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





## Build your own VPN business

Using [the cloudflare page UnchainAdmin](https://github.com/unchainese/unchainadmin) start your own VPN business. 