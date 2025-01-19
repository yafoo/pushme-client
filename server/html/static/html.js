const Html = {
    template: `<iframe ref="html" :srcdoc="doc" style="border: none; width: 100%;" :style="{height}" @load="onLoad"></iframe>`,
    props: {
        content: {
            default: ''
        },
    },
    computed: {
        doc() {
            if(~this.content.indexOf('<body') && ~this.content.indexOf('</body>')) {
                return this.content;
            } else {
                return `<style>body::-webkit-scrollbar{width:0;}body{padding:0;margin:0;color:#333333;font:14px Helvetica Neue,Helvetica,PingFang SC,Microsoft YaHei,Tahoma,Arial,sans-serif;line-height:1.6;}body>div{overflow:hidden;}</style><body><div>${this.content}</div></body>`;
            }
        }
    },
    data() {
        return {
            height: 'auto',
        }
    },
    mounted() {
        console.log('html mounted');
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
        onLoad() {
            console.log('onLoad');
            this.handleResize();
        },
        handleResize() {
            const height = this.$refs.html.contentWindow.document.body.scrollHeight;
            this.height = height + 'px';
            console.log('handleResize', height);
        },
    }
};