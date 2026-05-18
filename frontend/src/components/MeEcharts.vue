<template>
<div ref="echarts" style="width: 100%; height:80px;" :style="{height: + height + 'px'}"></div>
</template>

<script>
import { markRaw } from 'vue';
import * as echarts from 'echarts';

export default {
    props: {
        content: {
            type: String,
            default: ''
        },
    },
    data() {
        return {
            height: 80,
            chart: null,
            handleResize: null,
        }
    },
    computed: {
        option() {
            return this.content ? JSON.parse(this.content) : {};
        }
    },
    mounted() {
        console.log('echarts mounted');
        this.clacSize();
        this.$nextTick(() => {
            this.init();
        });
    },
    unmounted() {
        console.log('echarts unmounted');
        this.chart && this.chart.dispose();
        this.chart = null;
        if(this.handleResize) {
            window.removeEventListener('resize', this.handleResize);
        }
    },
    methods: {
        init() {
            if(!this.content) {
                return;
            }
            if(!this.chart) {
                console.log('init echarts', echarts);
                this.chart = markRaw(echarts.init(this.$refs.echarts));
            }
            this.chart.setOption(this.option);
            if(!this.handleResize) {
                this.handleResize = () => {
                    if(!this.$refs.echarts) {
                        return;
                    }
                    this.clacSize();
                    this.$nextTick(() => {
                        this.chart.resize();
                    });
                }
                window.addEventListener('resize', this.handleResize);
            }
        },
        clacSize() {
            if(!this.$refs.echarts) {
                return;
            }
            let width = this.$refs.echarts.offsetWidth;
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
}
</script>