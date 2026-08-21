<template>
<div class="setting" v-if="form">
    <div class="setting-menu">
        <div class="setting-menu-item" v-for="menu in formItem" :key="menu.key" @click="menuClick(menu.key)" :class="menu.key == current ? 'hover' : ''">{{menu.name}}</div>
    </div>
    <div class="setting-form" :scrollIntoView="current">
        <div class="form-section" v-for="menu in formItem" :key="menu.key" :id="menu.key">
            <div class="form-list">
                <div class="form-item" v-for="item in menu.items" :key="item.key">
                    <div class="form-label">{{item.label}} <span class="form-label-tips" v-if="getLabelTips(menu, item)">{{getLabelTips(menu, item)}}</span></div>
                    <div class="form-content">
                        <input class="input" v-if="item.type == 'input'" v-model="form[menu.key][item.key]" />
                        <select class="select" v-else-if="item.type == 'select'" v-model="form[menu.key][item.key]">
                            <template v-if="item.key == 'enable'">
                                <option v-for="val in enables" :key="val" :value="val">{{val ? '开启' : '关闭'}}</option>
                            </template>
                            <template v-else-if="item.key == 'offline_msg'">
                                <option v-for="val in offlines" :key="val" :value="val">{{val ? '接收' : '不接收'}}</option>
                            </template>
                            <template v-else-if="item.key == 'tls'">
                                <option v-for="val in tlsOptions" :key="val" :value="val">{{val}}</option>
                            </template>
                            <template v-else-if="item.key == 'html_js'">
                                <option v-for="val in enables" :key="val" :value="val">{{val ? '启用' : '禁用'}}</option>
                            </template>
                            <template v-else-if="item.key == 'method'">
                                <option v-for="val in methods" :key="val" :value="val">{{val}}</option>
                            </template>
                            <template v-else-if="item.key == 'duration'">
                                <option v-for="val in durations" :key="val" :value="val">{{val}}</option>
                            </template>
                            <template v-else-if="item.key == 'audio'">
                                <option v-for="val in audioList" :key="val" :value="val">{{val}}</option>
                            </template>
                            <template v-else-if="item.key == 'speak'">
                                <option v-for="val in speakList" :key="val" :value="val">{{val || '关闭'}}</option>
                            </template>
                            <template v-else-if="item.key == 'verify_key' || item.key == 'push_key'">
                                <option v-for="val in enables" :key="val" :value="val">{{val ? '开启' : '关闭'}}</option>
                            </template>
                        </select>
                        <div class="form-text" v-else-if="item.type == 'text'">
                            {{item.text}}
                        </div>
                        <div class="form-link" v-else-if="item.type == 'link'"><span @click="openUrl(item.url)">{{item.text || item.url}}</span></div>
                    </div>
                    <div class="form-tips" v-if="item.tips">{{item.tips}}</div>
                    <div class="form-tips cert-download" v-if="item.key == 'tls' && certUrl">默认下载地址：<span @click="openUrl(certUrl)" style="cursor: pointer;">{{certUrl}}</span></div>
                </div>

                <div class="section-tips" v-if="menu.key == 'api' && form && form.api.enable && form.api.port">
                    <div class="form-label">本机接口地址示例：</div>
                    <div class="form-link" v-for="ip in ips" :key="ip">http://{{formatIP(ip)}}:{{form.api.port}}</div>
                </div>
                <div class="section-tips" v-if="menu.key == 'other'">
                    <div class="form-label">消息总数：<span class="form-text">{{msgCount}}</span> <div class="button form-button" @click="clearMessage">清空</div></div>
                </div>
                <div v-if="menu.key == 'about'">
                    <div class="form-label">版本：<span class="form-text">{{version}}</span> <div class="button form-button" @click="checkVersion(true)">检查更新</div></div>
                    <div class="form-label">官网：<span class="form-link" @click="openUrl('https://push.i-i.me/')">https://push.i-i.me/</span></div>
                    <div class="form-label">仓库：<span class="form-link" @click="openUrl('https://github.com/yafoo/pushme-client')">GitHub</span> | <span class="form-link" @click="openUrl('https://gitee.com/yafu/pushme-client')">Gitee</span></div>
                </div>
            </div>
        </div>
    </div>
</div>

<div class="form-tools">
    <div class="button" @click="saveSetting">保存</div>
    <div class="button normal" @click="getDefault">默认值</div>
    <div class="button secondary" v-if="isRestart" @click="restart">重启</div>
</div>
</template>

<script>
import { GoGetIps, GoGetVersion, GoRestart, GoOpenBrowser } from "../../bindings/PushMe/internal/services/utilsservice";
import { GoGetSetting, GoSaveSetting, GoGetSettingDefault } from "../../bindings/PushMe/internal/services/settingservice";
import { GoClearMessage, GoGetMessageCount } from "../../bindings/PushMe/internal/services/messageservice";
import { GoOpenPlugin } from "../../bindings/PushMe/internal/services/appservice";
import { WebToast, WebConfirm, WebCheckVersion } from "../utils/common";
import { Events } from '@wailsio/runtime';

export default {
    data() {
        return {
            form: null,
            setstr: null,
            isRestart: false,
            formItem: [
                {name: '接口服务', key: 'api', items: [
                    {label: '状态', key: 'enable', type: 'select', tips: '默认：开启'},
                    {label: '监听IP', key: 'ip', type: 'input', tips: '默认：为空，支持本机所有IP'},
                    {label: '监听端口', key: 'port', type: 'input', tips: '默认：3010'},
                    {label: '验证push_key', key: 'verify_key', type: 'select', tips: '默认：关闭'},
                ]},
                {name: '自建服务', key: 'host', items: [
                    {label: '状态', key: 'enable', type: 'select', tips: '默认：关闭'},
                    {label: 'TLS/SSL', key: 'tls', type: 'select', tips: '提示：自签名证书需要导入浏览器中'},
                    {label: 'IP或域名', key: 'ip', type: 'input', tips: '提示：不加http或https，示例：192.168.1.1 或 www.example.com 或 [::1]'},
                    {label: '服务端口', key: 'port', type:'input', tips: 'PushMe Server v1.3.0+默认端口：3010，PushMe Server v2.0.0+默认端口：3100'},
                    {label: '离线消息', key: 'offline_msg', type: 'select', tips: '默认：不接收'},
                    {label: 'push_key', key: 'push_key', type: 'input', tips: 'PushMe APP 获取的push_key，并需在PushMe Server上配置'},
                ]},
                {name: '消息插件', key: 'plugin', items: []},
                {name: '消息转发', key: 'repost', items: [
                    {label: '状态', key: 'enable', type: 'select', tips: '默认：关闭'},
                    {label: '转发网址', key: 'url', type: 'input', tips: '完整URL，包含http前缀和端口号'},
                    {label: '请求方式', key: 'method', type: 'select', tips: '默认：POST/JSON'},
                    {label: '关键词限制', key: 'limit', type: 'input', tips: 'title包含关键词，多个以|隔开'},
                    {label: '关键词排除', key: 'omit', type: 'input', tips: 'title不含关键词，多个以|隔开'},
                    {label: '转发push_key', key: 'push_key', type: 'select', tips: '默认：关闭'},
                ]},
                {name: '桌面通知', key: 'notice', items: [
                    {label: '状态', key: 'enable', type: 'select', tips: '默认：开启'},
                    // {label: '持续时间', key: 'duration', type: 'select', tips: '默认：short'},
                    // {label: '提示音乐', key: 'audio', type: 'select', tips: '默认：default'},
                    {label: '语音通知', key: 'speak', type: 'select', tips: '默认：关闭'},
                ]},
                {name: '其他设置', key: 'other', items: [
                    {label: 'HTML消息JS支持', key: 'html_js', type: 'select', tips: '默认：禁用'},
                ]},
                {name: '系统设置', key: 'system', items: [
                    {label: '开机启动', key: 'enable', type: 'select', tips: '提示：目前仅支持windows系统，默认：关闭'},
                ]},
                {name: '关于我们', key: 'about', items: []},
            ],
            enables: [true, false],
            offlines: [false, true],
            tlsOptions: ['无证书', '公共证书', '自签名证书'],
            methods: ['GET', 'POST/JSON', 'POST/FORM'],
            durations: ['short', 'long'],
            audioList: ['default', 'im', 'mail', 'reminder', 'sms', 'loopingalarm', 'loopingalarm2', 'loopingalarm3', 'loopingalarm4', 'loopingalarm5', 'loopingalarm6', 'loopingalarm7', 'loopingalarm8', 'loopingalarm9', 'loopingalarm10', 'loopingcall', 'loopingcall2', 'loopingcall3', 'loopingcall4', 'loopingcall5', 'loopingcall6', 'loopingcall7', 'loopingcall8', 'loopingcall9', 'loopingcall10', 'silent'],
            speakList: ['', 'title', 'content', 'title+content'],
            current: 'api',
            ips: [],
            version: '获取中..',
            msgCount: '获取中..',
        }
    },
    computed: {
        certUrl() {
            if(this.form && this.form.host.tls == '自签名证书') {
                return `https://${this.form.host.ip + ':' + this.form.host.port}/certs/download`;
            }
            return '';
        }
    },
    watch: {
        form: {
            deep: true,
            handler(val) {
                if(val) {
                    if(JSON.stringify(val) != this.setstr) {
                        this.isRestart = true;
                    } else {
                        this.isRestart = false;
                    }
                }
            }
        }
    },
    created() {
        this.getIps();
        this.getSetting();
        this.getVersion();
        this.getMcount();
    },
    mounted() {
        this.checkVersion();
    },
    methods: {
        getSetting() {
            GoGetSetting().then(res => {
                this.form = res;
                this.setstr = JSON.stringify(res);
            });
        },
        getIps() {
            GoGetIps().then(res => {
                this.ips = [...res];
            });
        },
        formatIP(ip) {
            return ~ip.indexOf(':') ? `[${ip}]` : ip;
        },
        getVersion() {
            GoGetVersion().then(res => {
                this.version = res;
            });
        },
        getMcount() {
            GoGetMessageCount().then(res => {
                this.msgCount = res;
            });
        },
        getDefault() {
            GoGetSettingDefault().then(res => {
                this.form = res;
                WebToast('恢复默认值，保存后生效！');
            });
        },
        getLabelTips(menu, item) {
            let tips = '';
            if(menu.key == 'api' && item.key == 'ip') {
                if(this.form[menu.key].ip === '') {
                    tips = '置空，所有本机ip有效！';
                } else if(!~this.ips.indexOf(this.form[menu.key].ip.replace('[', '').replace(']', ''))) {
                    tips = '非本机IP，请确认！';
                }
            }
            return tips ? `(${tips})` : '';
        },
        saveSetting() {
            if(this.form.host.enable) {
                if(this.form.host.ip === '') {
                    this.menuClick('host');
                    return WebToast('请填写自建服务IP或域名！');
                } else if(this.form.host.ip.indexOf('http') > -1) {
                    this.menuClick('host');
                    return WebToast('自建服务IP或域名不能带http://或https://！');
                }
            }
            if(this.form.repost.enable) {
                if(!this.form.repost.url) {
                    this.menuClick('repost');
                    return WebToast('请填写接收消息转发的网址！');
                }
            }
            GoSaveSetting(this.form).then(res => {
                if (res === true) {
                    WebToast('保存成功，重启后生效！');
                } else {
                    WebToast('保存失败！');
                }
            });
        },
        menuClick(key) {
            if(key == 'plugin') {
                return GoOpenPlugin();
            }
            this.current = key;
            document.getElementById(key).scrollIntoView();
        },
        checkVersion(tips = false) {
            WebCheckVersion(tips);
        },
        clearMessage() {
            WebConfirm('该操作会清空通知消息，不会删除数据消息。', '确定清空？', action => {
                if(action != 'ok') {
                    return;
                }
                GoClearMessage().then(res => {
                    if(res == true) {
                        Events.Emit('message:del', {type: 'text'});
                        this.getMcount();
                        WebToast('清空成功！!');
                    } else {
                        WebToast('清空失败!');
                    }
                });
            });
        },
        openUrl(url) {
            GoOpenBrowser(url).then(res => {
                if(!res) {
                    WebToast('打开链接失败', 2000);
                }
            });
        },
        restart() {
            GoRestart();
        },
    }
}
</script>

<style scoped>
.setting {
    display: flex;
    height: 100%;
    width: 100%;
}
.setting-menu {
    width: 90px;
    padding-right: 1px;
    position: relative;
    text-align: center;
}
.setting-menu::after {
    position: absolute;
    right: 0;
    top: 0;
    height: 100%;
    width: 0;
    content: "";
    display: block;
    border-right: 1px solid #d2d2d2;
}
.setting-menu-item {
    padding: 8px 0;
    font-size: 12px;
    cursor: pointer;
    position: relative;
    transition: background-color .2s;
}
.setting-menu-item:hover,
.setting-menu-item.hover {
    background-color: var(--color-primary);
    color: var(--color-onPrimary);
}
.setting-menu-item:hover::after,
.setting-menu-item.hover::after {
    position: absolute;
    right: 0;
    top: 0;
    height: 100%;
    width: 0;
    content: "";
    display: block;
    border-right: 1px solid var(--color-primary);
    z-index: 1;
}

.setting-form {
    flex: 1;
    overflow-y: hidden;
    position: relative;
}
.form-section {
    height: 100vh;
    overflow-y: auto;
    overflow-y: overlay;
}
.form-list {
    padding: 5px calc(100% + 90px - 100vw + 8px) 50px 8px;
}
.form-item {
    margin-bottom: 8px;
}
.form-label-tips {
    font-size: 12px;
    color: red;
}
.form-tips {
    font-size: 10px;
    color: gray;
    line-height: 1.2;
    padding: 3px 0;
}
.form-text {
    color: var(--color-secondary);
}
.form-link {
    color: var(--color-primary);
    word-break: break-all;
    line-height: 1.2;
    padding: 3px;
    overflow: hidden;
}
.form-tools {
    position: absolute;
    bottom: 15px;
    right: 15px;
    z-index: 1;
    display: flex;
    gap: 8px;
    justify-content: center;
}
.section-tips {
    border-top: 1px dashed #ccc;
    padding-top: 10px;
}
#about {
    padding-top: 10px;
}
#about .form-label {
    padding-bottom: 8px;
    border-bottom: 1px solid #eee;
    margin-bottom: 8px;
}
#about .form-link {
    text-decoration: underline;
    cursor: pointer;
}

.form-button {
    box-sizing: content-box;
    font-size: 12px;
    height: 12px;
    line-height: 12px;
    padding: 4px 5px;
    float: right;
}
</style>