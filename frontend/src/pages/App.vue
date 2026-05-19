<template>
<div class="container">
    <div class="list-item" style="text-align: center; font-size: 12px;" v-if="!list.length">暂无消息</div>
    <template v-for="message in list" :key="message.id">
        <div class="list-item" @click="openMessage(message)" v-if="isTextMsg(message)" :class="'theme-' + parseTitle(message.title).theme">
            <div class="list-title">{{parseTitle(message.title).title}}</div>
            <div class="list-content" v-if="isMarkMsg(message)" v-html="renderMark(message.content)"></div>
            <div class="list-content" v-else-if="isHtmlMsg(message)" v-html="renderMark(removeStyleScript(message.content))"></div>
            <div class="list-content" v-else>{{message.content}}</div>
            <div class="list-date">{{message.date}}</div>
        </div>
    </template>
    <div class="button" v-if="hasMore" style="margin: 15px 0;" @click="getList(page+1)">加载更多</div>
</div>
<div class="float-tools">
    <div class="button-cirle button-refresh" @click="getList(1)"><img class="icon" src="/icon/refresh.svg"></div>
    <div class="button-cirle button-setting" @click="openSetting"><img class="icon" src="/icon/setting.svg"></div>
    <div class="button-cirle button-dashboard" @click="openDashboard"><img class="icon" src="/icon/dashboard.svg"></div>
</div>
</template>

<script>
import { GoGetMessageList, GoAddMessage } from "../../bindings/PushMe/internal/services/messageservice";
import { GoGetPluginListWithState } from "../../bindings/PushMe/internal/services/pluginservice";
import { GoGetHost, GoCheckVersion } from "../../bindings/PushMe/internal/services/settingservice";
import { GoOpenMessage, GoOpenDashboard, GoOpenSetting, GoNotification } from "../../bindings/PushMe/internal/services/appservice";
import { Events } from '@wailsio/runtime'
import {isTextMsg, isMarkMsg, isHtmlMsg, parseTitle, isDataMsg, WebCheckVersion, WebToast} from "../utils/common";

export default {
    data() {
        return {
            id: 0,
            list: [],
            page: 1,
            pageSize: 10,
            hasMore: false,
            md: null,
            host: null,
            client: null,
            plugin: null,
        }
    },
    created() {
        this.getList(1);
    },
    mounted() {
        this.init();
        Events.On('msg', ({data}) => {
            console.log('Events Api NewMessage', data[0]);
            this.newMessage(data[0]);
        });
        Events.On('toast', ({data}) => {
            console.log('Events Toast', data[0]);
            WebToast(data[0]);
        });
        Events.On('notification:action', ({data}) => {
            console.log('Events Notification action', data[0]);
            GoOpenMessage(data[0].userInfo.id).then(() => {
                console.log('OpenMessage Success');
            }).catch(err => {
                WebToast('消息打开失败:' + err.message);
            });
        });
        Events.On('plugin_change', () => {
            this.initPlugin();
        });
        Events.On('message_del', (detail) => {
            isTextMsg(detail) && this.getList(1);
        });
    },
    unmounted() {
        Events.Off('msg');
        Events.Off('toast');
        Events.Off('notification:action');
        Events.Off('plugin_change');
        Events.Off('message_del');
    },
    methods: {
        init() {
            this.initPlugin();

            this.initHost();

            this.checkVersion();
        },
        getList(page = 1) {
            this.page = page;
            GoGetMessageList(this.page, this.pageSize).then(list => {
                if(list.length < this.pageSize) {
                    this.hasMore = false;
                } else {
                    this.hasMore = true;
                }

                if(this.page == 1) {
                    this.list = list;
                    this.$nextTick(() => {
                        const dom = document.getElementById('app');
                        dom && dom.scrollTo({
                            top: 0,
                            behavior: 'smooth'
                        });
                    });
                } else {
                    list.forEach(item => {
                        this.list.push(item);
                    });
                }
            }).catch(err => {
                console.log(err);
            });
        },
        openMessage(msg) {
            if(window.getSelection().toString() === '') {
                GoOpenMessage(msg.id);
            }
        },
        newMessage(msg) {
            if(this.plugin) {
                this.pluginCalc(msg);
            } else {
                this.addMessage(msg);
            }
        },
        addMessage(msg) {
            msg.title = msg.title.replace(/\[~.+\]$/g, ''); //过滤通道信息

            GoAddMessage(msg).then(newMsg => {
                if(newMsg.id > 0) {
                    if(isDataMsg(newMsg)) {
                        Events.Emit('new_message', newMsg);
                    } else if(isTextMsg(newMsg)) {
                        this.getList(1);
                        const titles = parseTitle(newMsg.title);
                        const title = ({'': '', i: '[info]', s: '[success]', f: '[failure]', w: '[warning]'})[titles.theme] + ' ' + titles.title;
                        let content = newMsg.content;
                        if(newMsg.type == 'html') {
                            content = this.removeStyleScript(newMsg.content);
                        }
                        content = content.replace(/\s+/g, ' ').replace(/[#*]+\s/g, '').replace(/<[^>]+>|&[^>]+;/g, '');
                        GoNotification({...newMsg, title, content}).then(() => {
                            console.log('Notification Success');
                        }).catch(err => {
                            WebToast('消息发送失败:' + err.message);
                        });
                    }
                } else {
                    WebToast('消息发送失败');
                }
            });
        },
        openDashboard() {
            GoOpenDashboard();
        },
        getDashboard() {
            const list = [];
            this.list.forEach(item => {
                if(isDataMsg(item)) {
                    list.push({...item});
                }
            });
            return list;
        },
        openSetting() {
            GoOpenSetting();
        },
        isTextMsg(msg) {
            return isTextMsg(msg);
        },
        isMarkMsg(msg) {
            return isMarkMsg(msg);
        },
        isHtmlMsg(msg) {
            return isHtmlMsg(msg);
        },
        async renderMark(content) {
            if(!this.md) {
                const { getMd } = await import('../utils/md');
                this.md = getMd();
            }
            return this.md.renderInline(content).replace(/<[^>]+>|&[^>]+;/g, '').replace(/[#*]+\s/g, '');
        },
        removeStyleScript(html) {
            return html.replace(/<(style|script)[^>]*>[\s\S]*?<\/\1>|<\s*(style|script)[^>]*\/>/gi, '');
        },
        parseTitle(title) {
            return parseTitle(title);
        },
        async initHost() {
            let host =null;
            try {
                const res = await GoGetHost();
                console.log('GetSettingHost', res);
                host = res;
            } catch (err) {
                return console.log(err);
            }

            if(!host) {
                return;
            }

            if(!host.enable) {
                return console.log('自建服务未开启');
            }
            if(!host.ip || !host.port || !host.push_key) {
                return WebToast('自建服务参数配置不全', {...host});
            }

            let ip = host.ip + '';
            if(~ip.indexOf(':')) {
                !~ip.indexOf('[') && (ip = `[${ip}]`);
                !~ip.indexOf('[[') && (ip = `[${ip}]`);
            }
            const protocol = !host.tls || host.tls == "无证书" ? "ws" : "wss";
            const url = `${protocol}://${ip}:${host.port}`;
            const sub_topic = host.push_key;
            const client_id = 'pc_' + sub_topic;
            const options = {
                clientId: client_id,
                keepalive: 300,
            };
            if(host.offline_msg) {
                options.clean = false;
            }

            let mqtt = null;
            try {
                const res = await import('mqtt');
                mqtt = res.default;
            } catch (err) {
                return console.log(err);
            }
            if(!mqtt) {
                return WebToast('MQTT 模块加载失败');
            }

            this.client = mqtt.connect(url, options);

            this.client.on('connect', res => {
                console.log('Host Connected');
                this.client.subscribe(sub_topic, {qos: host.offline_msg ? 1 : 0});
                GoLog(`connect to ${url} success`);
            });

            this.client.on('message', (topic, payload) => {
                if(topic == sub_topic) {
                    let data = payload.toString();
                    console.log('Host NewMessage', data);
                    try {
                        data = JSON.parse(data);
                        this.newMessage(data);
                    } catch(e) {
                        WebToast(e);
                    }
                }
            });
        },
        async initPlugin() {
            try {
                const list = await GoGetPluginListWithState(1);
                console.log('PluginEnabledCount', list.length);
                this.plugin && this.plugin.terminate && this.plugin.terminate();
                if(list.length == 0) {
                    this.plugin = null;
                    return;
                }

                let {plugin: pluginJs} = await import('../utils/plugin.js');
                let pluginStr = '';
                list.forEach(p => {
                    pluginStr += `plugins.push(\n    ${p.content.replaceAll("\n", "\n    ")}\n);\n`;
                });
                pluginJs = pluginJs.replace('// plugin_list', pluginStr);

                const blob = new Blob([pluginJs], { type: 'application/javascript' });
                const url = URL.createObjectURL(blob);
                this.plugin = new Worker(url);
                this.plugin.onmessage = e => {
                    this.addMessage(e.data);
                };

                URL.revokeObjectURL(url);
            } catch (err) {
                WebToast('插件初始化失败:' + err.message);
            }
        },
        pluginCalc(msg) {
            this.plugin && this.plugin.postMessage && this.plugin.postMessage(msg);
        },
        checkVersion() {
            Events.Once("version", ({data}) => {
                console.log('Events version', data[0]);
                data[0] && WebCheckVersion(JSON.parse(data[0]))
            })
            GoCheckVersion();
        },
    }
}
</script>

<style scoped>
.container {
    min-height: 100vh;
    background-color: #f5f5f5;
}

.list-item {
    margin: 8px 0;
    padding: 8px 10px;
    background-color: #fff;
    border-radius: 3px;
    box-shadow: 0 3px 5px rgba(0,0,0,.01);
    line-height: 1.4;
    cursor: pointer;
    user-select: text;
}
.list-item.theme-i {
    border-left: 4px solid #AAAAAA;
    padding-left: 6px;
}
.list-item.theme-s {
    border-left: 4px solid #4CAF50;
    padding-left: 6px;
}
.list-item.theme-w {
    border-left: 4px solid #FDD835;
    padding-left: 6px;
}
.list-item.theme-f {
    border-left: 4px solid #E91E63;
    padding-left: 6px;
}
.list-title {
    font-size: 16px;
}
.list-content {
    margin: 5px 0;
    display: -webkit-box;
    line-clamp: 2;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-all;
}
.list-content,
.list-date {
    font-size: 14px;
    color: #555;
}
.button-dashboard {
    position: relative;
    bottom: 0;
    right: 0;
    z-index: 2;
}
.button-refresh,
.button-setting {
    margin-bottom: -42px;
    opacity: 0;
    transition: all 0.3s;
}
.float-tools:hover .button-refresh,
.float-tools:hover .button-setting {
    margin-bottom: 8px;
    opacity: 1;
}
.float-tools:hover .button-refresh:hover,
.float-tools:hover .button-setting:hover {
    opacity: 0.8;
}
</style>
