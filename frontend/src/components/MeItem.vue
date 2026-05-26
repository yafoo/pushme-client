<template>
<div class="card item" :class="theme">
    <div class="item-head"><div class="item-title">{{title}}</div><div class="item-date">{{date}}</div></div>
    <div class="item-content" v-if="content !== ''">{{content}}</div>
</div>
</template>

<script>
import { markRaw } from 'vue'
import { parseTitle, isMarkMsg, isHtmlMsg, removeStyleScript, getShortDate } from "../utils/common";

export default {
    props: {
        message: {
            type: Object,
            default: {}
        },
    },
    data() {
        return {
            title: '',
            theme: '',
            content: '',
            date: '',
            md: null,
        }
    },
    watch: {
        message: {
            deep: true,
            handler() {
                this.date = getShortDate(this.message.date);
                this.parseTitle();
                this.renderContent();
            },
            immediate: true,
        }
    },
    methods: {
        parseTitle() {
            const res = parseTitle(this.message.title);
            this.title = (res.title + '' || '').replace(/^\[#/, '[');
            this.theme = res.theme ? 'theme ' + res.theme : '';
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
            this.content = this.md.renderInline(content).replace(/<[^>]+>|&[^>]+;/g, '').replace(/[#*]+\s/g, '');
        },
    }
}
</script>

<style scoped>
.item {
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