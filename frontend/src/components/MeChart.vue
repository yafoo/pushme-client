<template>
<canvas ref="chart" style="width: 100%;"></canvas>
</template>

<script>
import chart from '../utils/chart.js';
export default {
    props: {
        content: {
            type: String,
            default: ''
        },
    },
    data() {
        return {}
    },
    watch: {
        content: {
            handler(val) {
                val && this.drawChart();
            }
        }
    },
    mounted() {
        console.log('chart mounted');
        this.init();
        window.addEventListener('resize', this.handleResize);

        this.drawChart();
    },
    unmounted() {
        console.log('chart unmounted');
        window.removeEventListener('resize', this.handleResize);
    },
    methods: {
        init() {
            this.handleResize();
        },
        handleResize() {
            if(!this.$refs.chart) {
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
            const canvas = this.$refs.chart;
            const width = canvas.offsetWidth;
            canvas.getContext("2d").clearRect(0, 0, width, width);
            chart.drawChart(canvas, JSON.parse(this.content));
        }
    }
}
</script>