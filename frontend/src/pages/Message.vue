<template>
<div class="container">
    <h1 class="message-title" :class="theme">{{title}}</h1>
    <div class="message-date">{{message.date}}</div>
    <div class="message-content"><me-content :message="message"></me-content></div>
</div>
<div class="float-tools" v-if="id > 0">
    <div class="button-cirle button-edit" v-if="isNote" @click="editNote"><img class="icon" src="/icon/edit.svg"></div>
    <div class="button-cirle" @click="delMessage"><img class="icon" src="/icon/delete.svg"></div>
</div>
</template>

<script>
import { GoGetMessage, GoDelMessage } from "../../bindings/PushMe/internal/services/messageservice";
import { GoOpenNoteEdit } from "../../bindings/PushMe/internal/services/appservice";
import { WebToast, WebConfirm, query, isNoteMsg } from "../utils/common";
import { calcTitleInfo } from "../utils/message";
import { Events } from '@wailsio/runtime'
import MeContent from "../components/MeContent.vue";

export default {
    components: { MeContent },
    data() {
        return {
            id: 0,
            message: {},
            title: '',
            theme: '',
            isNote: false,
        }
    },
    watch: {
        message: {
            deep: true,
            handler() {
                const res = calcTitleInfo(this.message.title);
                this.title = (res.title + '' || '');
                this.theme = res.theme ? 'theme ' + res.theme : '';
                this.isNote = isNoteMsg(this.message);
            },
            immediate: true,
        }
    },
    created() {
        this.init();
    },
    mounted() {
        Events.On('message:new', ({data}) => {
            if(data.id == this.id) {
                console.log('message:new', data);
                this.message = {...data};
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
            }).catch(err => {
                WebToast('消息获取失败:' + err.message);
            });
        },
        editNote() {
            if(!this.id) {
                return WebToast('便签不存在！');
            }
            GoOpenNoteEdit(this.id);
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
                }).catch(err => {
                    WebToast('删除失败:' + err.message);
                });
            });
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
.theme::before {
    content: '';
    border-left: 4px solid #fff;
    margin-right: 5px;
}
.theme.i::before {
    content: '';
    display: inline-block;
    border-left: 4px solid #fff;
    margin-right: 5px;
}
.theme.i::before {
    border-left-color: var(--color-theme-i);
}
.theme.s::before {
    border-left-color: var(--color-theme-s);
}
.theme.w::before {
    border-left-color: var(--color-theme-w);
}
.theme.f::before {
    border-left-color: var(--color-theme-f);
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

.button-edit {
    margin-bottom: -42px;
    opacity: 0;
    transition: all 0.3s;
}
#app:hover .float-tools .button-edit {
    margin-bottom: 8px;
    opacity: 1;
}
.float-tools:hover .button-edit:hover {
    opacity: 0.8;
}
</style>