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
                return `<style>html::-webkit-scrollbar{width:0;height:0;}html{overflow-y:hidden;}</style>${this.content}`;
            } else {
                return `<style>html::-webkit-scrollbar{width:0;height:0;}html{overflow-y:hidden;}body{padding:0;margin:0;color:#333333;font:14px Helvetica Neue,Helvetica,PingFang SC,Microsoft YaHei,Tahoma,Arial,sans-serif;line-height:1.6;}</style><body>${this.content}</body>`;
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
            const height = this.$refs.html.contentWindow.document.documentElement.scrollHeight;
            this.height = height + 'px';
            console.log('handleResize', height);
        },
    }
};