<template>
<div class="row">
    <div class="face" v-if="showFace && user !== ''" @click.stop="openUser">
        <me-avatar :user="user" :face="face"></me-avatar>
        <div class="face-user">{{user}}</div>
    </div>
    <div class="card item" :class="themeClass">
        <div class="item-head"><div class="item-title">{{title}}</div><div class="item-date">{{date}}</div></div>
        <div class="item-content" v-if="content !== ''">{{content}}</div>
    </div>
</div>
</template>

<script>
import { markRaw } from 'vue'
import { isMarkMsg, isHtmlMsg, removeStyleScript, getShortDate } from "../utils/common";
import { GoOpenUser } from "../../bindings/PushMe/internal/services/appservice";
import MeAvatar from "./MeAvatar.vue";

export default {
    components: { MeAvatar },
    props: {
        message: {
            type: Object,
            default: {}
        },
        showFace: {
            type: Boolean,
            default: true
        }
    },
    data() {
        return {
            themeClass: '',
            user: '',
            face: '',
            title: '',
            content: '',
            date: '',
            md: null,
        }
    },
    watch: {
        message: {
            deep: true,
            handler() {
                this.parseMsg();
            },
            immediate: true,
        }
    },
    methods: {
        openUser() {
            if(this.user) {
                GoOpenUser(this.user);
            }
        },
        parseMsg() {
            // 直接从 message 对象读取已解析的字段
            this.themeClass = this.message.theme ? `theme ${this.message.theme}` : '';
            this.user = this.message.user || '';
            this.face = this.message.face || '';
            this.title = this.message.title;
            this.date = getShortDate(this.message.date);
            this.renderContent();
        },
        async renderContent() {
            let content = this.message.content;
            let isHTml = false;
            if(isHtmlMsg(this.message)) {
                content = removeStyleScript(content);
                isHTml = true;
            }
            if(!isHTml && !isMarkMsg(this.message)) {
                this.content = content;
                return;
            }

            if(!this.md) {
                const { getMd } = await import('../utils/md');
                this.md = getMd();
            }
            this.content = this.md.parseInline(content).replace(/<[^>]+>|&[^>]+;/g, '').replace(/[#*]+\s/g, '');
        },
    }
}
</script>

<style scoped>
.row {
    display: flex;
    align-items: center;
    gap: 5px;
}
.face {
    width: 42px;
    display: flex;
    flex-direction: column;
    gap: 3px;
    cursor: pointer;
}
.face-user {
    height: 10px;
    font-size: 10px;
    line-height: 1;
    letter-spacing: 0.5px;
    text-align: center;
    overflow: hidden;
}
.item {
    flex: 1;
    line-height: 1.4;
    cursor: pointer;
    user-select: text;
}
.item.theme {
    border-left: 3px solid #fff;
    padding-left: 5px;
}
.item.theme.i {
    border-left-color: var(--color-theme-i);
}
.item.theme.s {
    border-left-color: var(--color-theme-s);
}
.item.theme.w {
    border-left-color: var(--color-theme-w);
}
.item.theme.f {
    border-left-color: var(--color-theme-f);
}

.item-head {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.item-title {
    flex: 1;
    font-size: 14px;
    font-weight: 400;
    display: -webkit-box;
    line-clamp: 1;
    -webkit-line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-all;
    color: #000;
}
.item-date {
    font-size: 10px;
    color: #555;
}
.item-content {
    margin-top: 5px;
    display: -webkit-box;
    line-clamp: 2;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-all;
    font-size: 12px;
    color: #666;
}
</style>