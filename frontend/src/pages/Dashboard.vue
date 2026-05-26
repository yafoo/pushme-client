<template>
<div class="container">
    <div class="dashboard">
        <div class="dashboard-col" v-for="num in this.column" :key="num" :ref="'col' + num">
            <div class="card dashboard-item" v-for="message in this.colList[num-1]" :key="message.id" @click="openMessage(message)">
                <div class="dashboard-title">{{message.title}}</div>
                <div class="dashboard-date">{{message.date}}</div>
                <div class="dashboard-content"><me-content :message="message"></me-content></div>
            </div>
        </div>
    </div>
</div>
</template>

<script>
import { GoGetDataList } from "../../bindings/PushMe/internal/services/messageservice";
import { GoOpenMessage } from "../../bindings/PushMe/internal/services/appservice";
import { parseTitle, WebToast, isDataMsg, isMarkMsg, isChartMsg, isEChartsMsg } from "../utils/common";
import { Events } from '@wailsio/runtime'
import MeContent from "../components/MeContent.vue";

export default {
    components: { MeContent },
    data() {
        return {
            list: [],
            column: 2,
            colList: [],
            isMonting: false,
        }
    },
    created() {
        this.init();
    },
    mounted() {
        setTimeout(() => {
            this.resize();
        }, 100);
        window.addEventListener('resize', _ => {
            this.$nextTick(() => {
                this.resize();
            });
        });
        Events.On('message:new', (msg) => {
            if(isDataMsg(msg)) {
                console.log('message:new', msg);
                this.updateMessage(msg);
            }
        });
        Events.On('message:del', (detail) => {
            if(isDataMsg(detail)) {
                this.getList();
            }
        });
    },
    unmounted() {
        Events.Off('message:new');
        Events.Off('message:del');
    },
    methods: {
        init() {
            this.getList();
        },
        getList() {
            GoGetDataList().then(list => {
                this.list = list;
                this.initColumn();
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
        isMarkMsg(msg) {
            return isMarkMsg(msg);
        },
        isChartMsg(msg) {
            return isChartMsg(msg);
        },
        isEChartsMsg(msg) {
            return isEChartsMsg(msg);
        },
        renderMark(content) {
            if(!this.md) {
                this.md = markdownit({html: true, linkify: false});
            }
            return this.md.render(content);
        },
        resize() {
            const column = Math.round(window.innerWidth / 150);
            if(column != this.column) {
                this.column = column;
                this.initColumn();
                this.$nextTick(() => {
                    this.mountBoard();
                });
            }
        },
        initColumn() {
            this.colList = [];
            for(let i=0; i<this.column; i++) {
                this.colList[i] = [];
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
                const col_i = this.getColumn();
                this.colList[col_i].push(this.list[index]);
                this.$nextTick(() => {
                    this.mountBoard(index + 1);
                });
            } else {
                this.isMonting = false;
            }
        },
        getColumn() {
            const heights = [];
            for(let i=1; i<=this.column; i++) {
                heights.push(this.$refs['col'+i][0].offsetHeight);
            }
            const min_height = Math.min(...heights);

            let col_i = 0;
            for(let i=0; i<this.column; i++) {
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
.dashboard .markdown>*:first-child {
    margin-top: 8px !important;
}
</style>

<style scoped>
.container {
    min-height: 100vh;
    background-color: #f5f5f5;
}
.dashboard {
    display: flex;
    gap: 5px;
    align-items: flex-start;
}
.dashboard-col {
    flex: 1;
    overflow: hidden;
}
.dashboard-item {
    margin: 8px 0;
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