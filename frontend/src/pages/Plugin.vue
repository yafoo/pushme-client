<template>
<div class="container">
    <div class="card plugin-item" v-for="(item, index) in list" :key="item.id">
        <div class="plugin-title" @click="addPlugin(item.id)">{{item.title}}</div>
        <div class="plugin-action">
            <div class="plugin-up" @click="SwapSort(index, 'up')">
                <svg xmlns="http://www.w3.org/2000/svg" :fill="index == 0 ? grayColor : iconColor" height="20px" viewBox="0 -960 960 960" width="20px"><path d="M480-525 291-336l-51-51 240-240 240 240-51 51-189-189Z"/></svg>
            </div>
            <div class="plugin-down" @click="SwapSort(index, 'down')">
                <svg xmlns="http://www.w3.org/2000/svg" :fill="index == listLen - 1 ? grayColor : iconColor" height="20px" viewBox="0 -960 960 960" width="20px"><path d="M480-333 240-573l51-51 189 189 189-189 51 51-240 240Z"/></svg>
            </div>
            <div class="plugin-check" @click="toggleState(item)">
                <svg xmlns="http://www.w3.org/2000/svg" v-show="item.state == 0" :fill="grayColor" height="20px" viewBox="0 -960 960 960" width="20px"><path d="M480.28-96Q401-96 331-126t-122.5-82.5Q156-261 126-330.96t-30-149.5Q96-560 126-629.5q30-69.5 82.5-122T330.96-834q69.96-30 149.5-30t149.04 30q69.5 30 122 82.5T834-629.28q30 69.73 30 149Q864-401 834-331t-82.5 122.5Q699-156 629.28-126q-69.73 30-149 30Zm-.28-72q130 0 221-91t91-221q0-130-91-221t-221-91q-130 0-221 91t-91 221q0 130 91 221t221 91Zm0-312Z"/></svg>
                <svg xmlns="http://www.w3.org/2000/svg" v-show="item.state == 1" :fill="iconColor" height="20px" viewBox="0 -960 960 960" width="20px"><path d="M480.23-288Q560-288 616-344.23q56-56.22 56-136Q672-560 615.77-616q-56.22-56-136-56Q400-672 344-615.77q-56 56.22-56 136Q288-400 344.23-344q56.22 56 136 56Zm.05 192Q401-96 331-126t-122.5-82.5Q156-261 126-330.96t-30-149.5Q96-560 126-629.5q30-69.5 82.5-122T330.96-834q69.96-30 149.5-30t149.04 30q69.5 30 122 82.5T834-629.28q30 69.73 30 149Q864-401 834-331t-82.5 122.5Q699-156 629.28-126q-69.73 30-149 30Zm-.28-72q130 0 221-91t91-221q0-130-91-221t-221-91q-130 0-221 91t-91 221q0 130 91 221t221 91Zm0-312Z"/></svg>
            </div>
        </div>
    </div>
</div>

<div class="float-tools">
    <div class="button-cirle" @click="addPlugin()"><img class="icon" src="/icon/add.svg"></div>
</div>
</template>

<script>
import { GoGetPluginList, GoSwapPlugin, GoEditPluginState } from "../../bindings/PushMe/internal/services/pluginservice";
import { GoOpenPluginEdit } from "../../bindings/PushMe/internal/services/appservice";
import { WebToast } from "../utils/common";
import { Events } from '@wailsio/runtime'

export default {
    data() {
        return {
            iconColor: '#46bc99',
            grayColor: '#999999',
            list: [],
            listLen: 0,
        }
    },
    created() {
        this.init();
    },
    mounted() {
        Events.On('plugin:change', () => {
            this.getList();
        });
    },
    unmounted() {
        Events.Off('plugin:change');
    },
    methods: {
        init() {
            this.getList();
        },
        getList() {
            GoGetPluginList().then(list => {
                this.list = list;
                this.listLen = list.length;
            });
        },
        addPlugin(id = 0) {
            GoOpenPluginEdit(id);
        },
        SwapSort(index, updown) {
            if(index == 0 && updown == 'up') {
                return WebToast('已经是第一个了');
            } else if(index == this.listLen - 1 && updown == 'down') {
                return WebToast('已经是最后一个了');
            }
            GoSwapPlugin(this.list[index], this.list[updown == 'up' ? index - 1 : index + 1]).then(res => {
                Events.Emit('plugin:change');
            }).catch(err => {
                WebToast('操作失败:' + err.message);
            });
        },
        toggleState(plugin) {
            plugin.state = plugin.state == 0 ? 1 : 0;
            GoEditPluginState(plugin).then(res => {
                Events.Emit('plugin:change');
            }).catch(err => {
                WebToast('操作失败:' + err.message);
            });
        }
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

.plugin-item {
    margin: 8px 0;
    padding: 1px;
    line-height: 1.4;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 10px;
}
.plugin-title {
    padding-left: 8px;
    flex: 1;
    display: -webkit-box;
    line-clamp: 1;
    -webkit-line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    cursor: pointer;
    line-height: 36px;
}
.plugin-action {
    display: flex;
    line-height: 1;
}
.plugin-action>div {
    padding: 7px;
    border-radius: 100%;
    transform: all 0.3s;
    cursor: pointer;
}
.plugin-action>div:hover {
    background-color: #f5f5f5;
}
</style>