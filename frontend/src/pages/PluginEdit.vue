<template>
<div class="container plugin-form">
    <div class="plugin-form-item">
        <input class="input" v-model="plugin.title" placeholder="插件标题" />
    </div>
    <div class="plugin-form-item plugin-form-textarea">
        <textarea class="textarea" v-model="plugin.content" placeholder="插件代码"></textarea>
    </div>
    <div class="plugin-form-item plugin-form-button">
        <div class="button" @click="edit">保存</div>
        <div class="button second" v-if="id > 0" @click="del">删除</div>
    </div>
</div>
</template>

<script>
import { GoGetPlugin, GoEditPlugin, GoDelPlugin } from "../../bindings/PushMe/internal/services/pluginservice";
import { query, WebToast, WebConfirm } from "../utils/common";
import { Events } from '@wailsio/runtime';

export default {
    data() {
        return {
            id: 0,
            plugin: {},
            defaultCode: `/**
* @name 插件名
* @param {Object} msg - 消息对象
* @param {String} msg.title - 标题
* @param {String} msg.content - 内容
* @param {String} msg.date - 时间
* @param {String} msg.type - 类型
* @param {Function} next - 下一个插件
*/
function(msg, next) {
// your code
next();
}`
        }
    },
    created() {
        this.init();
    },
    methods: {
        init() {
            if(query().get('id') > 0) {
                this.id = parseInt(query().get('id'));
                this.get();
            } else {
                this.plugin.content = this.defaultCode;
            }
        },
        get() {
            GoGetPlugin(this.id).then(plugin => {
                this.plugin = plugin;
            }).catch(err => {
                WebToast('消息插件失败:' + err.message);
            });
        },
        edit() {
            if(!this.plugin.title) {
                return WebToast('请填写插件标题！');
            }
            GoEditPlugin(this.plugin).then(id => {
                if(id > 0) {
                    Events.Emit('plugin:change');
                    WebToast('保存成功！', 500, _ => {
                        Events.Emit('plugin:close', this.id);
                    });
                } else {
                    WebToast('保存失败！');
                }
            }).catch(err => {
                WebToast('保存失败:' + err.message);
            });
        },
        del() {
            if(!this.id) {
                return WebToast('插件不存在或已删除！');
            }

            WebConfirm('确定删除？', action => {
                if(action != 'ok') {
                    return;
                }
                GoDelPlugin(this.id).then(res => {
                    if(res === true) {
                        Events.Emit('plugin:change');
                        WebToast('删除成功！', 500, _ => {
                            Events.Emit('plugin:close', this.id);
                            this.id = 0;
                            this.plugin = {};
                        });
                    } else {
                        WebToast('删除失败！');
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
.plugin-form {
    display: flex;
    height: 100%;
    flex-direction: column;
    padding-top: 8px;
}
.plugin-form-item {
    margin-bottom: 8px;
}
.plugin-form-textarea {
    flex: 1;
}
.plugin-form-textarea .textarea {
    height: 100%;
}
.plugin-form-button {
    display: flex;
    gap: 8px;
    justify-content: center;
}
</style>