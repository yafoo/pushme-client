<template>
<div class="svg" v-html="svgContent"></div>
</template>

<script>
export default {
    props: {
        content: {
            default: ''
        },
    },
    data() {
        return {
            svgContent: '',
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
            this.svgContent = (await this.getDOMPurify()).sanitize(this.content);
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
    }
}
</script>

<style>
.svg>svg {
    display: block;
    width: 100%;
    height: auto;
}
</style>