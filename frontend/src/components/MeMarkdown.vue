<template>
<div class="markdown" v-html="htmlContent"></div>
</template>

<script>
import { getMd } from "../utils/md";
import { proxyImages } from "../utils/common";
import "../../public/markdown.css";

export default {
    props: {
        content: {
            default: ''
        },
    },
    data() {
        return {
            htmlContent: '',
            purify: null,
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
    methods: {
        async render() {
            if(!this.content) {
                return this.htmlContent = this.content;
            }
            let html = getMd().parse(this.content);
            html = await proxyImages(html);
            this.htmlContent = (await this.getDOMPurify()).sanitize(html);
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
    },

}
</script>

<style scoped>

</style>