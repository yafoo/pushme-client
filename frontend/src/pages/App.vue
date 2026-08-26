<template>
<div class="container">
    <div class="list-empty" v-if="!list.length">暂无消息</div>
    <me-item class="list-item" v-for="message in list" :key="message.id" @click="openMessage(message)" :message="message"></me-item>
    <div class="button" v-if="hasMore" style="margin: 15px 0;" @click="getList(page+1)">加载更多</div>
</div>

<div class="float-tools">
    <div class="button-cirle" @click="getList()"><img class="icon" src="/icon/refresh.svg"></div>
    <div class="button-cirle" @click="openSetting"><img class="icon" src="/icon/setting.svg"></div>
    <div class="button-cirle" @click="openDashboard"><img class="icon" src="/icon/dashboard.svg"></div>
</div>

<div class="status-bar">
    <div class="status-item">
        <span class="status-dot" :class="{'status-on': apiStatus}"></span>
        <span class="status-text">接口</span>
    </div>
    <div class="status-item">
        <span class="status-dot" :class="{'status-on': hostStatus.enable}"></span>
        <span class="status-text">自建</span>
    </div>
    <div class="status-item" v-if="hostStatus.enable">
        <span class="status-dot" :class="{'status-connected': hostStatus.connected}"></span>
        <span class="status-text">{{ hostStatus.connected ? '已连接' : '连接..' }}</span>
    </div>
</div>
</template>

<script>
import { GoGetMessageListGrouped, GoAddMessage } from "../../bindings/PushMe/internal/services/messageservice";
import { GoOpenMessage, GoOpenDashboard, GoOpenSetting, GoOpenUser } from "../../bindings/PushMe/internal/services/appservice";
import { GoOpenBrowser } from "../../bindings/PushMe/internal/services/utilsservice";
import { GoGetSettingNotice, GoGetSetting } from "../../bindings/PushMe/internal/services/settingservice";

import { Events } from '@wailsio/runtime';
import {isTextMsg, isDataMsg, isUrlMsg, WebCheckVersion, WebToast, WebNotification, WebSpeakMsg } from "../utils/common";
import { initPlugin } from "../utils/plugin";
import { initHost } from "../utils/host";
import { markRaw } from 'vue';
import MeItem from "../components/MeItem.vue";

export default {
    components: { MeItem },
    data() {
        return {
            list: [],
            page: 1,
            pageSize: 10,
            hasMore: false,
            plugin: null,
            settingNotice: {},
            apiStatus: false,
            hostStatus: { enable: false, connected: false },
        }
    },
    created() {
        this.getList();
        this.getSettingNotice();
        this.getApiStatus();
    },
    mounted() {
        this.init();

        Events.On('message:api', ({data}) => {
            console.log('message:api', data);
            this.calcMessage(data);
        });
        Events.On('message:del', ({data}) => {
            console.log('message:del', data);
            isTextMsg(data) && this.getList();
        });
        Events.On('plugin:change', () => {
            this.initPlugin();
        });
        Events.On('setting:change', () => {
            this.getSettingNotice();
            this.getApiStatus();
            this.initHost();
        });

        Events.On('notification:action', ({data}) => {
            console.log('notification:action', data);
            const userInfo = data.userInfo;
            if(isUrlMsg({type: userInfo.type})) {
                GoOpenBrowser(data.body).then(res => {
                    if(!res) {
                        WebToast('打开链接失败', 2000);
                    }
                });
            } else {
                GoOpenMessage(userInfo.id).catch(err => {
                    WebToast('消息打开失败:' + err.message);
                });
            }
        });
        Events.On('event:toast', ({data}) => {
            console.log('event:toast', data);
            WebToast(data);
        });
    },
    unmounted() {
        Events.Off('message:api');
        Events.Off('message:del');
        Events.Off('plugin:change');

        Events.Off('notification:action');
        Events.Off('event:toast');
        Events.Off('setting:change');
    },
    methods: {
        init() {
            this.initPlugin();
            this.initHost();

            WebCheckVersion();
        },
        getList(page = 1) {
            this.page = page;
            GoGetMessageListGrouped(this.page, this.pageSize).then(list => {
                this.hasMore = list.length >= this.pageSize;

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
                WebToast(err.message);
            });
        },
        addMessage(msg) {
            GoAddMessage(msg).then(newMsg => {
                if(newMsg.id > 0) {
                    if(isDataMsg(newMsg)) {
                        Events.Emit('message:new', newMsg);
                    } else {
                        this.settingNotice.enable && WebNotification(newMsg);
                        this.settingNotice.speak && WebSpeakMsg(newMsg, this.settingNotice.speak)
                        this.getList();
                    }
                } else {
                    WebToast('消息添加失败');
                }
            }).catch(err => {
                WebToast('消息添加失败:' + err.message);
            });
        },
        openMessage(msg) {
            if(window.getSelection().toString() === '') {
                if(msg.user) {
                    GoOpenUser(msg.user);
                } else if(isUrlMsg(msg)) {
                    GoOpenBrowser(msg.content).then(res => {
                        if(!res) {
                            WebToast('打开链接失败', 2000);
                        }
                    });
                } else {
                    GoOpenMessage(msg.id);
                }
            }
        },
        openDashboard() {
            GoOpenDashboard();
        },
        openSetting() {
            GoOpenSetting();
        },
        calcMessage(msg) {
            msg.title = msg.title.replace(/\[~.+\]$/g, ''); //过滤通道信息
            this.plugin && this.plugin(msg);
        },
        async initPlugin() {
            this.plugin = markRaw(await initPlugin((msg => {
                this.addMessage(msg);
            })));
        },
        async initHost() {
            initHost(msg => {
                console.log('host:message', msg);
                this.calcMessage(msg);
            }, status => {
                this.hostStatus = { ...status };
            });
        },
        async getSettingNotice() {
            this.settingNotice = await GoGetSettingNotice();
        },
        async getApiStatus() {
            try {
                const setting = await GoGetSetting();
                this.apiStatus = setting.api.enable;
            } catch (err) {
                console.error('获取接口状态失败:', err);
            }
        },
    }
}
</script>

<style>
::-webkit-scrollbar {
    background-color: #f5f5f5;
}
</style>

<style scoped>
.container {
    min-height: 100vh;
    background-color: #f5f5f5;
    overflow: hidden;
}

.list-empty {
    margin: 8px 0;
    padding: 8px 8px;
    background-color: #fff;
    border-radius: 3px;
    box-shadow: 0 3px 5px rgba(0,0,0,.01);
    line-height: 1.4;
    cursor: pointer;
    user-select: text;
    text-align: center;
    font-size: 12px;
}
.list-item {
    margin: 8px 0;
}

.status-bar {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: 14px;
    background-color: rgba(255, 255, 255, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 16px;
    font-size: 10px;
    color: #666;
    z-index: 1;
}

.status-item {
    display: flex;
    align-items: center;
    gap: 4px;
}

.status-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background-color: #ccc;
}

.status-dot.status-on {
    background-color: #52c41a;
}

.status-dot.status-connected {
    background-color: #1890ff;
}

.status-text {
    line-height: 1;
}
</style>
