# PushMeClient
PushMeClient是一个接收自[PushMe](https://push.i-i.me/)APP转发的消息的windows(win10+)客户端，以实现上班期间，在电脑上能及时快速的收到并查看新消息。

v2.0版本基本实现了PushMeAPP上的大部分功能，包括：插件和连接自建服务，同时接口服务也更智能，支持GET请求、POST/FORM请求、POST/JSON请求，支持界面配置，支持开机自启动。

除了接收PushMe的转发消息，也可以用来接收其他应用发来的消息，例如：[短信转发器APP](https://github.com/pppscn/SmsForwarder)，实现在电脑上查看短信、其他应用的通知。

> 为了简单轻量，一开始选用纯go来开发，使用fyneUI做界面，结果发现fyneUI显示字体有点模糊，内存占用也不低，开发起来也比较困难。于是从v2.0版本开始，改用webview来开发，最终效果非常满意。另外v2.0的部分功能参考了[http-win-notice](https://github.com/shanghaobo/http-win-notice)的代码，感谢作者。

## 仓库地址

Github：https://github.com/yafoo/pushme-client

Gitee：https://gitee.com/yafu/pushme-client

## 使用教程

请在[Releases](https://github.com/yafoo/pushme-client/releases)中下载最新版客户端，打开即可免安装运行。

软件默认运行在`3200`端口，在你的PushMe安卓客户端上，设置消息转发url为：`http://您局域网IP:3200`即可。

如果您的电脑有公网IP，也可以设置公网IP，这样，手机和电脑不在一个局域网也可以接收消息。

PushMeClient目前不能独立接收官网接口推送的消息，只能接收PushMeAPP转发的消息。如需独立使用，可以自行搭建一个接口服务，参考[PushMeServer](https://github.com/yafoo/pushme-server)。PushMeClient支持连接自建服务独立使用。

## 特色功能

除了PushMeAPP上的文本消息、markdown消息、数据消息、数据小屏外，PushMeClient还支持图表消息，图表支持实时动态更新。

PushMeAPP后面也会支持图表消息。

## 其他

1. 本人是第一次学习开发go，代码写的不好，欢迎提出意见建议。

## 部分界面