<template>
<div class="container note-edit">
    <div class="note-form">
        <div class="note-form-item">
            <input class="input" v-model="title" placeholder="标题" />
        </div>
        <div class="note-form-item note-form-textarea">
            <textarea class="textarea" v-model="content" placeholder="内容（支持 Markdown 格式）"></textarea>
        </div>
        <div class="note-form-item note-form-button">
            <div class="button" @click="save">保存</div>
            <div class="button second" v-if="id > 0" @click="del">删除</div>
        </div>
    </div>
</div>
</template>

<script>
import { GoGetMessage, GoUpdateMessage, GoAddMessage, GoDelMessage } from "../../bindings/PushMe/internal/services/messageservice";
import { query, WebToast, WebConfirm } from "../utils/common";
import { Events } from '@wailsio/runtime';

export default {
    data() {
        return {
            id: 0,
            title: '',
            content: '',
            message: null,
        }
    },
    created() {
        this.init();
    },
    methods: {
        init() {
            const idParam = query().get('id');
            if(idParam && parseInt(idParam) > 0) {
                this.id = parseInt(idParam);
                this.load();
            }
        },
        load() {
            GoGetMessage(this.id).then(msg => {
                if(msg && msg.id > 0) {
                    this.message = msg;
                    this.title = msg.title || '';
                    this.content = msg.content || '';
                } else {
                    WebToast('便签不存在或已删除！');
                }
            }).catch(err => {
                WebToast('加载失败:' + err.message);
            });
        },
        save() {
            if(!this.title && !this.content) {
                return WebToast('标题和内容不能同时为空！');
            }
            
            const now = new Date();
            const dateStr = now.getFullYear() + '-' + 
                String(now.getMonth() + 1).padStart(2, '0') + '-' + 
                String(now.getDate()).padStart(2, '0') + ' ' + 
                String(now.getHours()).padStart(2, '0') + ':' + 
                String(now.getMinutes()).padStart(2, '0') + ':' + 
                String(now.getSeconds()).padStart(2, '0');

            if(this.id > 0) {
                // 更新
                const msg = {
                    id: this.id,
                    title: this.title,
                    content: this.content,
                    date: dateStr,
                    type: 'note'
                };
                GoUpdateMessage(msg).then(res => {
                    if(res) {
                        Events.Emit('message:new', msg);
                        WebToast('保存成功！', 500, _ => {
                            Events.Emit('note:close', this.id);
                        });
                    } else {
                        WebToast('保存失败！请检查标题是否重复！');
                    }
                }).catch(err => {
                    WebToast('保存失败:' + err.message);
                });
            } else {
                // 新增
                const msg = {
                    title: this.title,
                    content: this.content,
                    date: dateStr,
                    type: 'note'
                };
                GoAddMessage(msg).then(newMsg => {
                    if(newMsg && newMsg.id > 0) {
                        Events.Emit('message:new', newMsg);
                        WebToast('添加成功！', 500, _ => {
                            Events.Emit('note:close', this.id);
                        });
                    } else {
                        WebToast('添加失败！请检查标题是否重复！');
                    }
                }).catch(err => {
                    WebToast('添加失败:' + err.message);
                });
            }
        },
        del() {
            if(!this.id) {
                return WebToast('便签不存在！');
            }
            WebConfirm('确定删除此便签？', action => {
                if(action != 'ok') {
                    return;
                }
                GoDelMessage(this.id).then(res => {
                    if(res) {
                        Events.Emit('message:del', {id: this.id, type: 'note'});
                        WebToast('删除成功！', 500, _ => {
                            Events.Emit('message:close', this.id);
                            Events.Emit('note:close', this.id);
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
.note-edit {
    display: flex;
    height: 100%;
    flex-direction: column;
    padding: 8px;
}
.note-form {
    display: flex;
    flex: 1;
    flex-direction: column;
    height: 100%;
}
.note-form-item {
    margin-bottom: 8px;
}
.note-form-textarea {
    flex: 1;
}
.note-form-textarea .textarea {
    height: 100%;
    resize: none;
}
.note-form-button {
    display: flex;
    gap: 8px;
    justify-content: center;
}
</style>
