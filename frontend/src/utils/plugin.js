import { GoGetPluginListWithState } from "../../bindings/PushMe/internal/services/pluginservice";
import {WebToast} from "../utils/common";

let plugin = null;
export const initPlugin = async(onMessage) => {
    const pluginFun = (msg) => {
        plugin && plugin.postMessage ? plugin.postMessage(msg) : onMessage(msg);
    }

    try {
        const list = await GoGetPluginListWithState(1);
        console.log('PluginEnabledCount', list.length);
        plugin && plugin.terminate && plugin.terminate();
        plugin = null
        if(!list.length) {
            return pluginFun;
        }

        let {default: pluginJs} = await import('../utils/plugin-tpl.js');
        let pluginStr = '';
        list.forEach(p => {
            pluginStr += `plugins.push(\n    ${p.content.replaceAll("\n", "\n    ")}\n);\n`;
        });
        pluginJs = pluginJs.replace('// plugin_list', pluginStr);

        const blob = new Blob([pluginJs], { type: 'application/javascript' });
        const url = URL.createObjectURL(blob);
        plugin = new Worker(url);
        plugin.onmessage = e => {
            onMessage(e.data)
        };

        URL.revokeObjectURL(url);
    } catch (err) {console.log(err);
        WebToast('插件初始化失败:' + err.message);
    }
    return pluginFun;
}