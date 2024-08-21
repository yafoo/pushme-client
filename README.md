# PushMeClient
PushMeClient是一个接收来自[PushMe](https://push.i-i.me/) APP转发的消息的windows客户端

## 仓库地址

Github：https://github.com/yafoo/pushme-client

Gitee：https://gitee.com/yafu/pushme-client

## 开发目的

此软件主要解决一个需求：对于上班族，在上班期间可能大部分时间都在操作电脑，没有时间查看手机，可能会错过重要消息。有了此软件，在电脑上就可以随时看到新消息，方便很多。

## 使用教程

请在[Releases](https://github.com/yafoo/pushme-client/releases)中下载最新版客户端，打开即可免安装运行。

软件默认运行在`3200`端口，在你的PushMe安卓客户端上，设置消息转发url为：`http://您局域网IP:3200`即可。

如果您的电脑有公网IP，也可以设置公网IP，这样，手机和电脑不在一个局域网也可以接收消息。

## 更换服务端口

如果您想使用其他端口号，直接在客户端文件名加上`-端口号`即可，例如：`PushMeClient-1234.exe`，再次打开运行，服务会运行在`1234`端口。

## 其他

1. 本人是第一次学习开发go，代码写的不好，欢迎提出意见建议。

2. 如果您有开发实力，也可以编写自己的http服务，实现在电脑上接收消息。

3. 软件暂时没有直接接收PushMe接口推送消息的计划，只支持接收客户端转发消息。
