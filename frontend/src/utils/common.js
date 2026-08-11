import { GoNotification, GoCheckVersion } from "../../bindings/PushMe/internal/services/utilsservice";
import { calcTitleInfo } from "./message";

export const query = () => new URLSearchParams(window.location.search)

export function WebToast(msg, time, callback) {
    console.log('WebToast:', msg);
    msg === undefined && (msg = '');
    typeof msg === 'object' && (msg = JSON.stringify(msg));
    typeof time === 'function' && (callback = time, time = 3000)
    time === undefined && (time = 3000);
    WebToast.index === undefined && (WebToast.index = 1000);
    WebToast.index++;
    const toast = document.createElement('div');
    toast.className = 'toast';
    toast.style.display = 'block';
    toast.style.zIndex = WebToast.index;
    toast.innerText = msg;
    toast.style.animation = 'fadein 0.5s, fadeout 0.5s ' + (time / 1000) + 's forwards';
    document.body.append(toast);
    setTimeout(function(){
        toast.parentNode.removeChild(toast);
        typeof callback == 'function' && callback();
    }, time);
}

export function WebConfirm(content, title, callback) {
    if(typeof title == 'function') {
        callback = title;
        title = '';
    }
    console.log('WebConfirm:', content, title);

    const dom_bg = document.createElement('div');
    dom_bg.className = 'confirm-bg';

    const dom = document.createElement('div');
    dom.className = 'confirm';
    const dom_title = document.createElement('div');
    dom_title.className = 'confirm-title';
    dom_title.innerText = title;

    const dom_content = document.createElement('div');
    dom_content.className = 'confirm-content';
    dom_content.innerText = content;

    let buttons = '';
    buttons += '<div class="confirm-button-ok" data-action="ok">确定</div>';
    if(typeof callback == 'function') {
        buttons = '<div class="confirm-button-cancel" data-action="cancel">取消</div>' + buttons;
    }
    const dom_button = document.createElement('div');
    dom_button.className = 'confirm-button';
    dom_button.innerHTML = buttons;
    dom_button.childNodes.forEach(btn => {
        btn.onclick = function(){
            dom.parentNode.removeChild(dom_bg);
            dom.parentNode.removeChild(dom);
            typeof callback == 'function' && callback(btn.dataset.action);
        }
    });

    title && dom.append(dom_title);
    dom.append(dom_content);
    dom.append(dom_button);
    document.body.append(dom_bg);
    document.body.append(dom);
}

export const TextTypes = ['text', 'markdown', 'html', 'url', '']
export const DataTypes = ['data', 'markdata', 'chart', 'echarts', 'svg']

export function isTextMsg(msg) {
    return TextTypes.includes(msg.type);
}

export function isDataMsg(msg) {
    return DataTypes.includes(msg.type);
}

export function isMarkMsg(msg) {
    return msg.type == 'markdown' || msg.type == 'markdata';
}

export function isChartMsg(msg) {
    return msg.type == 'chart';
}

export function isEChartsMsg(msg) {
    return msg.type == 'echarts';
}

export function isHtmlMsg(msg) {
    return msg.type == 'html';
}

export function isUrlMsg(msg) {
    return msg.type == 'url';
}

export function isSvgMsg(msg) {
    return msg.type == 'svg';
}

export function getShortDate(input) {
    try {
        // 尝试直接解析日期
        let date = new Date(input);
        
        // 如果直接解析失败，尝试处理不同格式
        if (isNaN(date.getTime())) {
            // 替换常见格式分隔符
            let normalized = input.replace(/\//g, '-');
            
            // 处理只有日期没有时间的情况
            if (normalized.match(/^\d{4}-\d{1,2}-\d{1,2}$/)) {
                normalized += 'T00:00:00';
            }
            
            date = new Date(normalized);
            
            if (isNaN(date.getTime())) {
                return input;
            }
        }
        
        const today = new Date();
        const isToday = date.getFullYear() === today.getFullYear() &&
                        date.getMonth() === today.getMonth() &&
                        date.getDate() === today.getDate();
        
        const isSameYear = date.getFullYear() === today.getFullYear();
        
        const formatNumber = (num) => num.toString().padStart(2, '0');
        
        if (isToday) {
            return `${formatNumber(date.getHours())}:${formatNumber(date.getMinutes())}`;
        } else if (isSameYear) {
            return `${date.getMonth() + 1}月${formatNumber(date.getDate())}日`;
        } else {
            return `${date.getFullYear()}-${formatNumber(date.getMonth() + 1)}-${formatNumber(date.getDate())}`;
        }
    } catch (e) {
        return input;
    }
}

export function removeStyleScript(html) {
    return html.replace(/<(style|script)[^>]*>[\s\S]*?<\/\1>|<\s*(style|script)[^>]*\/>/gi, '');
}

export function WebCheckVersion(tips = false) {
    GoCheckVersion().then(res => {
        console.log('WebCheckVersion:', res);
        if(!res) {
            return tips && WebToast('当前已是最新版本！');
        }
        if(res.startsWith('{')) {
            res = JSON.parse(res);
            if(res && res.update) {
                WebConfirm(res.update + '\n点击确定，打开下载链接', '发现新版本' + res.version, action => {
                    action == "ok" && window.open('https://github.com/yafoo/pushme-client/releases');
                });
            } else {
                tips && WebToast('当前已是最新版本！');
            }
        } else {
            tips && WebToast('请求出错：' + res);
        }
    }).catch(res => {
        WebToast('请求出错：' + res);
    });
}

export function WebNotification(message) {
    const msg = {...message};
    const titleInfo = calcTitleInfo(msg.title);
    msg.title = ({'': '', i: '⬜️', s: '🟩', f: '🟥', w: '🟨'})[titleInfo.theme] + (titleInfo.user ? `[${titleInfo.user}]` : '') + titleInfo.title;
    if(msg.type == 'html') {
        msg.content = removeStyleScript(msg.content);
    }
    msg.content = msg.content.replace(/\s+/g, ' ').replace(/[#*]+\s/g, '').replace(/<[^>]+>|&[^>]+;/g, '');
    GoNotification(msg).then(() => {
        console.log('Notification Success', msg);
    }).catch(err => {
        WebToast('消息发送失败:' + err.message);
    });
}

export function WebSpeakMsg(message, type) {
    let text = [];
    const msg = {...message};
    if(type.indexOf('title') > -1) {
        const titleInfo = calcTitleInfo(msg.title);
        text.push((titleInfo.user ? `[${titleInfo.user}]` : '') + titleInfo.title);
    }
    if(type.indexOf('content') > -1) {
        const content = msg.type == 'html' ? removeStyleScript(msg.content) : msg.content;
        text.push(content.replace(/\s+/g, ' ').replace(/[#*]+\s/g, '').replace(/<[^>]+>|&[^>]+;/g, ''));
    }

    text.length && WebSpeakText(text.join('\n'));
}

export function WebSpeakText(text) {
    if(!window.speechSynthesis) {
        return WebToast('当前浏览器不支持语音播报');
    }
    // 为了防止队列堆积，先取消正在进行的播报
    window.speechSynthesis.cancel();

    // 创建 SpeechSynthesisUtterance 实例
    const utterance = new SpeechSynthesisUtterance(text);
    
    // 设置语言为中文
    utterance.lang = 'zh-CN';
    
    // 设置参数
    utterance.rate = 1.0;      // 语速 (0.1 - 10)
    utterance.pitch = 1.0;     // 音调 (0 - 2)
    utterance.volume = 1.0;     // 音量 (0 - 1)

    // 事件监听
    utterance.onstart = () => console.log("开始播报...");
    utterance.onend = () => console.log("播报结束");
    utterance.onerror = (event) => console.error("播报出错", event);

    // 开始播报
    window.speechSynthesis.speak(utterance);
}