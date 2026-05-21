<template>
<div class="container">
    <div class="list-empty" style="text-align: center; font-size: 12px;" v-if="!list.length">暂无消息</div>
    <template v-for="message in list" :key="message.id">
        <me-item class="list-item" @click="openMessage(message)" v-if="isTextMsg(message)" :message="message"></me-item>
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
import { GoCheckVersion } from "../../bindings/PushMe/internal/services/utilsservice";
import { GoOpenMessage, GoOpenDashboard, GoOpenSetting } from "../../bindings/PushMe/internal/services/appservice";
import { Events } from '@wailsio/runtime'
import {isTextMsg, isDataMsg, WebCheckVersion, WebToast, Notification } from "../utils/common";
import { initPlugin } from "../utils/plugin";
import { initHost } from "../utils/host";
import { markRaw } from 'vue';
import MeItem from "../components/MeItem.vue";

export default {
    components: { MeItem },
    data() {
        return {
            id: 0,
            list: [],
            page: 1,
            pageSize: 10,
            hasMore: false,
            plugin: null,
        }
    },
    created() {
        this.getList(1);
    },
    mounted() {
        this.init();
        Events.On('message:api', ({data}) => {
            console.log('message:api', data[0]);
            this.newMessage(data[0]);
        });
        Events.On('message:del', (detail) => {
            isTextMsg(detail) && this.getList(1);
        });
        Events.On('plugin:change', () => {
            this.initPlugin();
        });

        Events.On('notification:action', ({data}) => {
            console.log('notification:action', data[0]);
            GoOpenMessage(data[0].userInfo.id).catch(err => {
                WebToast('消息打开失败:' + err.message);
            });
        });
        Events.On('event:toast', ({data}) => {
            console.log('Events Toast', data[0]);
            WebToast(data[0]);
        });
    },
    unmounted() {
        Events.Off('message:api');
        Events.Off('message:del');
        Events.Off('plugin:change');

        Events.Off('notification:action');
        Events.Off('event:toast');
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
        newMessage(msg) {
            this.plugin && this.plugin(msg);
        },
        addMessage(msg) {
            msg.title = msg.title.replace(/\[~.+\]$/g, ''); //过滤通道信息

            GoAddMessage(msg).then(newMsg => {
                if(newMsg.id > 0) {
                    if(isDataMsg(newMsg)) {
                        Events.Emit('message:new', newMsg);
                    } else if(isTextMsg(newMsg)) {
                        Notification(newMsg);
                        this.getList(1);
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
                GoOpenMessage(msg.id);
            }
        },
        openDashboard() {
            GoOpenDashboard();
        },
        openSetting() {
            GoOpenSetting();
        },
        isTextMsg(msg) {
            return isTextMsg(msg);
        },
        async initPlugin() {
            this.plugin = markRaw(await initPlugin((msg => {
                this.addMessage(msg);
            })));
        },
        async initHost() {
            initHost(msg => {
                this.newMessage(msg);
            })
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
}
.list-item {
    margin: 8px 0;
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
