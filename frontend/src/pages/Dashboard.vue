<template>
<div class="container">
    <div class="dashboard">
        <div class="dashboard-col" v-for="i in columns" :key="i" :ref="'col' + i">
            <div class="card dashboard-item" :class="{'note-card': message.type === 'note'}" v-for="message in this.columnList[i-1]" :key="message.id" @click="openMessage(message)">
                <div class="dashboard-title">{{message.title}}</div>
                <div class="dashboard-date">{{message.date}}</div>
                <div class="dashboard-content"><me-content :message="message"></me-content></div>
            </div>
        </div>
    </div>
</div>
<div class="float-tools">
    <div class="button-cirle" @click="addNote"><img class="icon" src="/icon/add.svg"></div>
</div>
</template>

<script>
import { GoGetDataList } from "../../bindings/PushMe/internal/services/messageservice";
import { GoOpenMessage, GoOpenNoteEdit } from "../../bindings/PushMe/internal/services/appservice";
import { WebToast, isDataMsg } from "../utils/common";
import { Events } from '@wailsio/runtime'
import MeContent from "../components/MeContent.vue";

export default {
    components: { MeContent },
    data() {
        return {
            list: [],
            columns: 2,
            columnList: [],
            isMonting: false,
        }
    },
    mounted() {
        this.getList();

        Events.On('message:new', ({data}) => {
            if(isDataMsg(data)) {
                console.log('message:new', data);
                this.updateMessage(data);
            }
        });
        Events.On('message:del', ({data}) => {
            if(isDataMsg(data)) {
                console.log('message:del', data);
                this.getList();
            }
        });

        window.addEventListener('resize', _ => {
            this.$nextTick(() => {
                this.resize();
            });
        });
        setTimeout(() => {
            this.resize();
        }, 100);
    },
    unmounted() {
        Events.Off('message:new');
        Events.Off('message:del');
    },
    methods: {
        getList() {
            GoGetDataList().then(list => {
                this.list = list;
                this.resetColumnList();
                this.$nextTick(() => {
                    this.mountBoard();
                });
            }).catch(err => {
                console.log(err);
            });
        },
        openMessage(msg) {
            if(window.getSelection().toString() === '') {
                GoOpenMessage(msg.id);
            }
        },
        addNote() {
            GoOpenNoteEdit(0);
        },
        updateMessage(msg) {
            let is_update = false;
            for(let i = 0; i < this.list.length; i++) {
                if(this.list[i].id == msg.id) {
                    // this.list.splice(i, 1, msg);
                    Object.keys(msg).forEach(key => {
                        this.list[i][key] = msg[key];
                    });
                    is_update = true;
                    break;
                }
            }
            if(!is_update) {
                this.getList();
            }
        },
        resize() {
            const newColumns = Math.round(window.innerWidth / 150);
            if(newColumns != this.columns) {
                this.columns = newColumns;
                this.resetColumnList();
                this.$nextTick(() => {
                    this.mountBoard();
                });
            }
        },
        resetColumnList() {
            this.columnList = [];
            for(let i=0; i<this.columns; i++) {
                this.columnList[i] = [];
            }
        },
        mountBoard(index=0) {
            if(index == 0) {
                if(this.isMonting) {
                    return;
                }
                this.isMonting = true;
            }
            if(this.list.length > index) {
                const col_i = this.getNextColumnIndex();
                this.columnList[col_i].push(this.list[index]);
                this.$nextTick(() => {
                    this.mountBoard(index + 1);
                });
            } else {
                this.isMonting = false;
            }
        },
        getNextColumnIndex() {
            const heights = [];
            for(let i=1; i<=this.columns; i++) {
                heights.push(this.$refs['col'+i][0].offsetHeight);
            }
            const min_height = Math.min(...heights);

            let col_i = 0;
            for(let i=0; i<this.columns; i++) {
                if(heights[i] == min_height) {
                    col_i = i;
                    break;
                }
            }

            return col_i;
        },
    }
}
</script>

<style>
::-webkit-scrollbar {
    background-color: #f5f5f5;
}
.dashboard .data {
    font-size: 18px;
    color: var(--color-primary);
    text-align: center;
    padding: 25px 0;
    line-height: 1.2;
    white-space: pre-wrap;
}
.dashboard .markdown {
    font-size: 14px;
    color: #555;
    line-height: 1.2;
}
.dashboard .markdown>* {
    margin: 0.5em 0 !important;
}
.dashboard .markdown>*:last-child {
    margin-bottom: 5px !important;
}
.dashboard .markdown>table:first-child {
    margin-top: 0.5em !important;
}
.dashboard .markdown h1, .dashboard .markdown h2 {
    font-size: 1.1em;
    padding-bottom: 0.2em;
}
.dashboard .markdown h3, .dashboard .markdown h4 {
    font-size: 1em;
}
.dashboard .markdown blockquote {
    padding: 0.4em 0.6em;
    border-left-width: 3px;
}
.dashboard .markdown pre {
    padding: 0.5em 0.6em;
    font-size: 11px;
    max-height: 15em;
}
.dashboard .markdown code {
    font-size: 0.9em;
}
.dashboard .markdown table th, .dashboard .markdown table td {
    padding: 0.3em 0.4em;
    font-size: 13px;
}
.dashboard .markdown li {
    margin: 0.3em 0 0.3em 16px;
}
.dashboard .markdown img {
    max-height: 150px;
    object-fit: contain;
}
.dashboard-content canvas {
    margin-top: 5px !important;
}
.dashboard-content .svg {
    margin-top: 5px !important;
}
</style>

<style scoped>
.container {
    min-height: 100vh;
    background-color: #f5f5f5;
    padding-top: 8px;
}
.dashboard {
    display: flex;
    gap: 5px;
    align-items: flex-start;
}
.dashboard-col {
    flex: 1;
}
.dashboard-item {
    margin-bottom: 8px;
    padding: 3px;
    cursor: pointer;
}
.dashboard-title {
    font-size: 12px;
    padding-top: 3px;
    text-align: center;
}
.dashboard-date {
    font-size: 8px;
    color: #999;
    line-height: 1.2;
    text-align: center;
}
</style>