export const bc = new BroadcastChannel('pushme');
export const query = () => new URLSearchParams(window.location.search)

export function WebToast(msg, time, callback) {
    console.log('WebToast:', msg);
    msg === undefined && (msg = '');
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

export function isTextMsg(msg) {
    return !msg.type || msg.type == 'text' || msg.type == 'markdown' || msg.type == 'html';
}

export function isDataMsg(msg) {
    return msg.type == 'data' || msg.type == 'markdata' || msg.type == 'chart' || msg.type == 'echarts';
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

export function parseTitle(title='') {
    const reg = /^\[([iswf])\]/;
    const res = reg.exec(title);
    if(res) {
        return {theme: res[1], title: title.replace(reg, '')};
    } else {
        return {theme: '', title};
    }
}

export function WebCheckVersion(res, tips=false) {
    console.log('WebCheckVersion:', res);
    if(res && res.update) {
        WebConfirm(res.update + '\n点击确定，打开下载链接', '发现新版本' + res.version, action => {
            action == "ok" && window.open('https://github.com/yafoo/pushme-client/releases');
        });
    } else {
        tips && WebToast('当前已是最新版本！');
    }
}
