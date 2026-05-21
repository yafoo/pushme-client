import { GoGetHost } from "../../bindings/PushMe/internal/services/settingservice";

let client = null;

export const initHost = async(onMessage) => {
    let host = null;
    try {
        const res = await GoGetHost();
        console.log('GetSettingHost', res);
        host = res;
    } catch (err) {
        return console.log(err);
    }

    if(!host) {
        return;
    }

    if(!host.enable) {
        return console.log('自建服务未开启');
    }
    if(!host.ip || !host.port || !host.push_key) {
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
        keepalive: 300,
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
        client.subscribe(sub_topic, {qos: host.offline_msg ? 1 : 0});
    });

    client.on('message', (topic, payload) => {
        if(topic == sub_topic) {
            let data = payload.toString();
            console.log('Host NewMessage', data);
            try {
                data = JSON.parse(data);
                onMessage(data);
            } catch(e) {
                WebToast(e);
            }
        }
    });
}