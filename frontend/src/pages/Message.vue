<template>
<div class="container">
    <h1 class="message-title"><span :class="'theme-' + parseTitle(message.title).theme"></span>{{parseTitle(message.title).title}}</h1>
    <div class="message-date">{{message.date}}</div>
    <div class="message-content"><me-content :message="message"></me-content></div>
</div>
<div class="float-tools" v-if="id > 0">
    <div class="button-cirle" @click="delMessage"><img class="icon" src="/icon/delete.svg"></div>
</div>
</template>

<script>
import { GoGetMessage, GoDelMessage } from "../../bindings/PushMe/internal/services/messageservice";
import { parseTitle, WebToast, WebConfirm, query } from "../utils/common";
import { Events } from '@wailsio/runtime'
import MeContent from "../components/MeContent.vue";

export default {
    components: { MeContent },
    data() {
        return {
            id: 0,
            message: {},
        }
    },
    created() {
        this.init();
    },
    mounted() {
        Events.On('message:new', (msg) => {
            if(msg.id == this.id) {
                console.log('message:new', msg);
                this.message = {...msg};
            }
        });
    },
    unmounted() {
        Events.Off('message:new');
    },
    methods: {
        init() {
            if (query().get('id')) {
                this.id = parseInt(query().get('id'));
                this.getMessage();
            } else {
                WebToast('消息ID不能为空');
            }
        },
        getMessage() {
            GoGetMessage(this.id).then(msg => {
                this.message = msg;
            });
        },
        delMessage() {
            if(!this.id) {
                return WebToast('消息不存在或已删除！');
            }
            WebConfirm('确定删除？', action => {
                if(action != 'ok') {
                    return;
                }
                GoDelMessage(this.id).then(res => {
                    if(res == true) {
                        Events.Emit('message:del', {...this.message});
                        WebToast('删除成功!', 500, _ => {
                            Events.Emit('message:close', this.id);
                            this.id = 0;
                            this.message = {};
                        });
                    } else {
                        WebToast('删除失败!');
                    }
                });
            });
        },
        parseTitle(title) {
            return parseTitle(title);
        },
    }
}
</script>

<style scoped>
.message-title {
    font-size: 16px;
    line-height: 1.4;
    padding-top: 10px;
}
.message-title .theme-i {
    border-left: 4px solid var(--color-theme-i);
    margin-right: 5px;
}
.message-title .theme-s {
    border-left: 4px solid var(--color-theme-s);
    margin-right: 5px;
}
.message-title .theme-w {
    border-left: 4px solid var(--color-theme-w);
    margin-right: 5px;
}
.message-title .theme-f {
    border-left: 4px solid var(--color-theme-f);
    margin-right: 5px;
}
.message-date {
    font-size: 12px;
    color: #888;
    line-height: 1.4;
    margin: 10px 0;
}
.message-content {
    font-size: 14px;
    color: #555;
    line-height: 1.5;
}
</style>