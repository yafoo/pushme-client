<template>
<component :is="currentComponent" :content="message.content" v-if="currentComponent" />
</template>

<script>
import { markRaw } from 'vue'
import { isTextMsg, isMarkMsg, isHtmlMsg, isChartMsg, isEChartsMsg, isDataMsg } from "../utils/common";

export default {
    props: {
        message: {
            type: Object,
            default: {}
        },
    },
    data() {
        return {
            currentComponent: null,
        }
    },
    watch: {
        message: {
            deep: true,
            handler() {
                console.log('watch', this.message);
                this.render();
            },
            immediate: true,
        }
    },
    methods: {
        async render() {
            let component = 'MeText';
            if(isMarkMsg(this.message)) {
                component = 'MeMarkdown';
            } else if(isHtmlMsg(this.message)) {
                component = 'MeHtml';
            } else if(isChartMsg(this.message)) {
                component = 'MeChart';
            } else if(isEChartsMsg(this.message)) {
                component = 'MeEcharts';
            } else if(isDataMsg(this.message)) {
                component = 'MeData';
            }

            const module = await import(`../components/${component}.vue`)
            this.currentComponent = markRaw(module.default);
        },
    }
}
</script>