import { GoGetHost } from "../../bindings/PushMe/internal/services/settingservice";
import { WebToast } from "../utils/common";

let client = null;
let hostStatus = { enable: false, connected: false };

export const getHostStatus = () => hostStatus;

export const initHost = async(onMessage, onStatusChange) => {
    let host = null;
    try {
        const res = await GoGetHost();
        console.log('GetSettingHost', res);
        host = res;
    } catch (err) {
        return console.log(err);
    }

    if(!host) {
        hostStatus = { enable: false, connected: false };
        onStatusChange && onStatusChange(hostStatus);
        return;
    }

    hostStatus.enable = host.enable;

    if(!host.enable) {
        hostStatus.connected = false;
        onStatusChange && onStatusChange(hostStatus);
        // 关闭旧的 MQTT 连接
        client && client.end();
        client = null;
        return console.log('自建服务未开启');
    }
    if(!host.ip || !host.port || !host.push_key) {
        hostStatus.connected = false;
        onStatusChange && onStatusChange(hostStatus);
        return WebToast('自建服务参数配置不全', {...host});
    }

    let ip = host.ip + '';
    if(~ip.indexOf(':')) {
        !~ip.indexOf('[') && (ip = `[${ip}]`);
        !~ip.indexOf('[[') && (ip = `[${ip}]`);
    }
    const protocol = !host.tls || host.tls == "无证书" ? "ws" : "wss";
    const url = `${protocol}://${ip}:${host.port}`;
    const sub_topic = host.push_key;
    const client_id = 'pc_' + sub_topic;
    const options = {
        clientId: client_id,
        keepalive: 120,
    };
    if(host.offline_msg) {
        options.clean = false;
    }

    let mqtt = null;
    try {
        const res = await import('mqtt');
        mqtt = res.default;
    } catch (err) {
        return console.log(err);
    }
    if(!mqtt) {
        return WebToast('MQTT 模块加载失败');
    }

    client && client.end();
    client = mqtt.connect(url, options);

    client.on('connect', res => {
        console.log('Host Connected');
        hostStatus.connected = true;
        onStatusChange && onStatusChange(hostStatus);
        client.subscribe(sub_topic, {qos: host.offline_msg ? 1 : 0});
    });
    client.on('disconnect', res => {
        console.log('Host Disconnected', res);
        hostStatus.connected = false;
        onStatusChange && onStatusChange(hostStatus);
    });
    client.on('close', res => {
        hostStatus.connected = false;
        onStatusChange && onStatusChange(hostStatus);
    });
    client.on('error', res => {
        hostStatus.connected = false;
        onStatusChange && onStatusChange(hostStatus);
    });

    client.on('message', (topic, payload) => {
        if(topic == sub_topic) {
            let data = payload.toString();
            try {
                data = JSON.parse(data);
                onMessage(data);
            } catch(e) {
                WebToast(e);
            }
        }
    });
}