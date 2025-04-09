# Unchain 翻墙服务器搭建和使用教程

## 服务器要求
建议使用linux ubuntu 服务器. 服务器配配置没有内存和CPU限制,任意配置都可以. 最便宜的服务器都可以.

服务求区域最好选择日本或者美国(访问OpenAI,Claude,Google Gemini友好).


## 安装部署

从 [https://github.com/unchainese/unchain/releases/tag/v0.0.3](https://github.com/unchainese/unchain/releases/tag/v0.0.3)
现在对应服务器架构的二进制文件,解压到任意目录,然后运行即可.

```bash
wget https://github.com/unchainese/unchain/releases/download/v0.0.3/unchain-linux-amd64.unchain.tar.gz
tar -zxvf unchain-linux-amd64.unchain.tar.gz
```

在上一步解压之后的可执行文件相同的目录创建 `config.toml`配置文件.
创建文件命令 `vim config.toml`.
文件内容详见 [https://github.com/unchainese/unchain/blob/v0.0.3/config.example.standalone.toml](https://github.com/unchainese/unchain/blob/v0.0.3/config.example.standalone.toml)

使用下面命令,来测试配置文件是否正确.
```bash
# cd 到 unchain 和 config.toml 所在目录
chmod +x unchain
./unchain
```
如果没有报错,说明配置文件正确.


## 使用systemctl 管理服务

在 `/etc/systemd/system/` 目录下创建 `unchain.service` 文件,
使用命令 `vim /etc/systemd/system/unchain.service` 创建文件.
文件内容如下:
[](https://github.com/unchainese/unchain/blob/v0.0.3/unchain.service)

 - `systemctl daemon-reload` 重新加载服务
 - `systemctl start unchain` 启动服务
 - `systemctl stop unchain` 停止服务
 - `systemctl restart unchain` 重启服务
 - `systemctl status unchain` 查看服务状态


## 使用V2ray/Clash/ShadowRocket 客户端连接

复制日志中的 VLESS链接,在V2ray/Clash/ShadowRocket 客户端中添加VLESS链接.