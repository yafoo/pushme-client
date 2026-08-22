<template>
<iframe ref="html" :srcdoc="htmlContent" style="border: none; width: 100%;" :style="{height}" @load="onLoad"></iframe>
</template>

<script>
import { GoGetHtmlJs } from "../../bindings/PushMe/internal/services/settingservice";
import { proxyImages } from "../utils/common";

export default {
    props: {
        content: {
            default: ''
        },
    },
    data() {
        return {
            height: 'auto',
            purify: null,
            enableJs: null,
            htmlContent: '',
        }
    },
    watch: {
        content: {
            immediate: true,
            handler() {
                this.render();
            }
        }
    },
    mounted() {
        window.addEventListener('resize', this.handleResize);
        window.resizeBody = () => {
            this.handleResize();
        }
    },
    unmounted() {
        console.log('html unmounted');
        window.removeEventListener('resize', this.handleResize);
        delete window.resizeBody;
    },
    methods: {
        async render() {
            let content = this.content;
            if(!(await this.getHtmlJs())) {
                content = await this.getDOMPurify().sanitize(this.content);
            }
            content = await proxyImages(content);
            if(~this.content.indexOf('<body') && ~this.content.indexOf('</body>')) {
                content = `<body>${content}</body>`;;
            }
            if(~this.content.indexOf('<body') && ~this.content.indexOf('</body>')) {
                this.htmlContent = `<style>html::-webkit-scrollbar{width:0;height:0;}html{overflow-y:hidden;}</style>${content}`;
            } else {
                this.htmlContent = `<style>html::-webkit-scrollbar{width:0;height:0;}html{overflow-y:hidden;}body{padding:0;margin:0;color:#333333;font:14px Helvetica Neue,Helvetica,PingFang SC,Microsoft YaHei,Tahoma,Arial,sans-serif;line-height:1.6;}</style><body>${content}</body>`;
            }
        },
        async getHtmlJs() {
            if(this.enableJs === null) {
                try {
                    this.enableJs = await GoGetHtmlJs();
                } catch(e) {
                    console.log('获取HtmlJs失败', e);
                }
            }
            return this.enableJs;
        },
        async getDOMPurify() {
            if(!this.purify) {
                try {
                    const {default: DOMPurify} = await import('dompurify');
                    this.purify = DOMPurify;
                } catch(e) {
                    console.log('DOMPurify初始化失败', e);
                }
            }
            return this.purify;
        },
        onLoad() {
            console.log('Iframe onLoad');
            this.handleResize();
        },
        handleResize() {
            const height = this.$refs.html.contentWindow.document.documentElement.scrollHeight;
            this.height = height + 'px';
            console.log('handleResize', height);
        },
    }
}
</script>