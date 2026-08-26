<template>
<div class="container">
    <div class="list-empty" v-if="!list.length && !loading">暂无消息</div>
    <me-item class="list-item" v-for="message in list" :key="message.id" @click="openMessage(message)" :message="message" :showFace="false"></me-item>
    <div class="button" v-if="hasMore && !loading" style="margin: 15px 0;" @click="getList(page+1)">加载更多</div>
</div>

<div class="float-tools">
    <div class="button-cirle" @click="deleteAll"><img class="icon" src="/icon/delete.svg"></div>
    <div class="button-cirle float-info">{{totalCount}}条<br>消息</div>
    <me-avatar class="float-avatar" :user="user" :face="face" :size="42"></me-avatar>
</div>
</template>

<script>
import { GoGetMessageListByUser, GoGetCountByUser, GoDelMessageByUser } from "../../bindings/PushMe/internal/services/messageservice";
import { GoOpenMessage } from "../../bindings/PushMe/internal/services/appservice";
import { GoOpenBrowser } from "../../bindings/PushMe/internal/services/utilsservice";
import { WebToast, WebConfirm, query, isUrlMsg } from "../utils/common";
import { Events } from '@wailsio/runtime';
import MeItem from "../components/MeItem.vue";
import MeAvatar from "../components/MeAvatar.vue";

export default {
    components: { MeItem, MeAvatar },
    data() {
        return {
            user: '',
            face: '',
            list: [],
            page: 1,
            pageSize: 10,
            hasMore: false,
            loading: false,
            totalCount: 0,
        }
    },
    created() {
        this.init();
    },
    mounted() {
        Events.On('message:new', ({data}) => {
            if(data.user === this.user) {
                this.list.unshift(data);
                this.totalCount++;
            }
        });
        Events.On('message:del', ({data}) => {
            if(data.user === this.user) {
                this.list = this.list.filter(m => m.id !== data.id);
                this.totalCount--;
            }
        });
    },
    unmounted() {
        Events.Off('message:new');
        Events.Off('message:del');
    },
    methods: {
        init() {
            const user = query().get('user');
            if(!user) {
                return WebToast('用户名为空');
            }
            this.user = user;
            this.getCount();
            this.getList();
        },
        getList(page = 1) {
            this.loading = true;
            GoGetMessageListByUser(this.user, page, this.pageSize).then(list => {
                this.loading = false;
                this.hasMore = list.length >= this.pageSize;
                if(page === 1) {
                    this.list = list;
                    this.$nextTick(() => {
                        const dom = document.getElementById('app');
                        dom && dom.scrollTo({
                            top: 0,
                            behavior: 'smooth'
                        });
                    });
                } else {
                    list.forEach(item => this.list.push(item));
                }
                this.page = page;
                // 从第一条消息中提取 face（如果尚未获取）
                if(!this.face && list.length > 0) {
                    this.face = list[0].face || '';
                }
            }).catch(err => {
                this.loading = false;
                WebToast('获取消息失败:' + err.message);
            });
        },
        getCount() {
            GoGetCountByUser(this.user).then(count => {
                this.totalCount = count;
            });
        },
        openMessage(msg) {
            if(window.getSelection().toString() === '') {
                if(isUrlMsg(msg)) {
                    GoOpenBrowser(msg.content);
                } else {
                    GoOpenMessage(msg.id);
                }
            }
        },
        deleteAll() {
            WebConfirm(`确定删除 "${this.user}" 的全部 ${this.totalCount} 条消息？`, '删除确认', action => {
                if(action != 'ok') {
                    return;
                }
                GoDelMessageByUser(this.user).then(res => {
                    if(res) {
                        WebToast('删除成功');
                        this.list = [];
                        this.totalCount = 0;
                        this.getList();
                        Events.Emit('message:del', {type: 'text'});
                    } else {
                        WebToast('删除失败');
                    }
                });
            });
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

.float-avatar {
    cursor: pointer;
}
.float-info {
    color: #fff;
    cursor: pointer;
    user-select: text;
    font-size: 10px;
    line-height: 1.2;
    text-align: center;
}
</style>
