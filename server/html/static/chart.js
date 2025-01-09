const Chart = {
    template: `<canvas ref="chart" style="width: 100%;"></canvas>`,
    props: {
        option: {
            type: Object,
            default: {}
        },
    },
    data() {
        return {}
    },
    watch: {
        option: {
            // deep: true,
            handler(val) {
                this.drawChart();
            }
        }
    },
    mounted() {
        console.log('chart mounted');
        this.init();
        window.addEventListener('resize', this.handleResize);
    },
    unmounted() {
        console.log('chart unmounted');
        window.removeEventListener('resize', this.handleResize);
    },
    methods: {
        init() {
            if(typeof(chart) === 'undefined') {
                if(document.getElementById('chart-js')) {
                    setTimeout(() => {
                        this.init();
                    }, 100);
                    return;
                }
                console.log('load chart js');
                const script = document.createElement('script');
                script.setAttribute('id', 'chart-js');
                script.src = '/static/chart.fun.js';
                script.onload = () => {
                    console.log('chart js loaded');
                    this.init();
                }
                document.body.appendChild(script);
                return;
            }
            this.handleResize();
        },
        handleResize() {
            if(!this.$refs.chart || typeof(chart) === 'undefined') {
                return;
            }
            const canvas = this.$refs.chart;
            const width = canvas.offsetWidth;
            canvas.style.height = width + "px";
            const dpr = window.devicePixelRatio || 1;
            canvas.width = width * dpr;
            canvas.height = width * dpr;
            this.drawChart();
        },
        drawChart() {
            if(!this.$refs.chart || typeof(chart) === 'undefined') {
                return;
            }
            const canvas = this.$refs.chart;
            const width = canvas.offsetWidth;
            canvas.getContext("2d").clearRect(0, 0, width, width);
            chart.drawChart(canvas, this.option);
        }
    }
};