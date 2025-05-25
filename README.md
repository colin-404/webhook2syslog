# syslog-webhook
监听http端口，将http请求写入到syslog端口


# 使用
先修改config.yaml配置文件
```
go build
webhook2syslog -c ./config.yaml
```