const Chart = {
    template: `<div ref="chart" style="width: 100%; height:80px;" :style="{height: + height + 'px'}"></div>`,
    props: {
        option: {
            type: Object,
            default: {}
        },
    },
    data() {
        return {
            height: 80,
            chart: null,
        }
    },
    watch: {
        option: {
            // deep: true,
            handler(val) {
                if(!this.chart) {
                    return;
                }
                this.chart.setOption(val);
            }
        }
    },
    mounted() {
        console.log('chart mounted');
        this.clacSize();
        this.$nextTick(() => {
            this.init();
        });
    },
    unmounted() {
        console.log('chart unmounted');
        this.chart && this.chart.dispose();
        this.chart = null;
    },
    methods: {
        init() {
            if(typeof(echarts) === 'undefined') {
                if(document.getElementById('echarts')) {
                    setTimeout(() => {
                        this.init();
                    }, 100);
                    return;
                }
                console.log('load echarts');
                const script = document.createElement('script');
                script.setAttribute('id', 'echarts');
                script.src = '/static/echarts.min.js';
                script.onload = () => {
                    console.log('echarts loaded');
                    this.init();
                }
                document.body.appendChild(script);
                return;
            }
            if(!echarts.init) {
                setTimeout(() => {
                    this.init();
                }, 100);
                return;
            }
            this.chart = Vue.markRaw(echarts.init(this.$refs.chart));
            this.chart.setOption(this.option);
            window.addEventListener('resize', () => {
                if(!this.$refs.chart) {
                    return;
                }
                this.clacSize();
                this.$nextTick(() => {
                    this.chart.resize();
                });
            });
        },
        clacSize() {
            if(!this.$refs.chart) {
                return;
            }
            let width = this.$refs.chart.offsetWidth;
            let height = Math.round(width * 0.8);
            if(this.option.height && this.option.width) {
                height = Math.round(this.option.height / this.option.width * width);
                if(width > this.option.width && height > this.option.height) {
                    width = this.option.width;
                    height = this.option.height;
                }
            }
            if(height > window.innerHeight) {
                height = window.innerHeight;
            }
            this.width = width;
            this.height = height;
        },
    }
};