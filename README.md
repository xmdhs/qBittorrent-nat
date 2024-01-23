# qBittorrent-nat
qBittorrent 打洞工具，思路来自 https://github.com/Mythologyli/qBittorrent-NAT-TCP-Hole-Punching

区别在于打包了打洞功能，使用上可能会更加方便一点。

## 要求
* nat1
* iptables

需要运行程序的设备是 nat 1，可以设置 dmz，或者路由器开启 upnp。具体参照 https://github.com/xmdhs/natupnp?tab=readme-ov-file#%E5%85%89%E7%8C%AB%E6%8B%A8%E5%8F%B7%E7%BD%91%E7%BA%BF%E6%8F%92%E5%9C%A8%E8%B7%AF%E7%94%B1%E5%99%A8%E7%9A%84-lan-%E5%8F%A3

## 使用
配置文件见 config.json，修改账号密码成 qBittorrent web ui 上设置的。

`qBittorrent-nat -c config.json`