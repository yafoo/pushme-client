<template>
    <div class="avatar" :style="avatarStyle">
        <img v-if="isImageAvatar" :src="face" class="avatar__image" :style="imageStyle" />

        <div v-else class="avatar__text">
            <div :style="textStyle">{{ displayText }}</div>
        </div>
    </div>
</template>

<script>
export default {
    name: 'Avatar',
    props: {
        user: {
            type: String,
            default: ''
        },
        face: {
            type: String,
            default: ''
        },
        size: {
            type: [Number, String],
            default: 42
        },
        textSize: {
            type: [Number, String],
            default: null
        },
    },
    computed: {
        isImageAvatar() {
            return this.face && this.face.startsWith('http')
        },
        bgColor() {
            if (this.isImageAvatar) {
                return 'transparent'
            }
            return this.generateColorFromUser(this.user)
        },
        avatarStyle() {
            return {
                width: `${this.size}px`,
                height: `${this.size}px`,
                backgroundColor: this.bgColor,
                borderRadius: '50%',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                overflow: 'hidden',
                flexShrink: 0
            }
        },
        imageStyle() {
            return {
                width: '100%',
                height: '100%',
                objectFit: 'cover'
            }
        },
        displayText() {
            if (this.face && !this.face.startsWith('http')) {
                return this.face
            }
            return this.getFirstChar(this.user)
        },
        textLength() {
            if (!this.displayText) return 0
            return this.getStringLength(this.displayText)
        },
        baseSize() {
            if (this.textSize) {
                return typeof this.textSize === 'number' ? `${this.textSize}px` : this.textSize
            }
            return (this.size / 2.4) + 'px'
        },
        fontSize() {
            const baseSize = this.baseSize.replace(/px/, '')
            const length = this.textLength

            if (length > 1 && length <= 4) {
                return `${baseSize / 2}px`
            }
            if (length > 4) {
                return `${baseSize / 3}px`
            }
            return `${baseSize}px`
        },
        textStyle() {
            const length = this.textLength
            let maxLines = 1
            if (length > 1 && length <= 4) {
                maxLines = 2
            } else if (length > 4) {
                maxLines = 3
            }

            return {
                color: '#FFFFFF',
                fontSize: this.fontSize,
                lineHeight: this.fontSize,
                textAlign: 'center',
                wordBreak: 'break-word',
                display: '-webkit-box',
                WebkitLineClamp: maxLines,
                WebkitBoxOrient: 'vertical',
                overflow: 'hidden',
                maxWidth: `calc(${this.baseSize} + ${this.fontSize} / 2)`,
                marginTop: '1px',
            }
        }
    },
    methods: {
        getStringLength(str) {
            if (typeof str !== 'string') return 0;
            
            if (typeof Intl.Segmenter === 'function') {
                const segmenter = new Intl.Segmenter(undefined, { granularity: 'grapheme' });
                return [...segmenter.segment(str)].length;
            }
            
            // 降级方案
            return [...str].length;
        },
        getFirstChar(str) {
            if (typeof str !== 'string' || str.length === 0) return '';

            if (typeof Intl.Segmenter === 'function') {
                const segmenter = new Intl.Segmenter(undefined, { granularity: 'grapheme' });
                const segments = segmenter.segment(str);

                for (const { segment } of segments) {
                    return segment;
                }
                return '';
            }

            // 降级方案
            return [...str][0] || '';
        },
        generateColorFromUser(user) {
            const colors = [
                // 粉色系 - 加深
                '#EC407A', // 玫红
                '#D81B60', // 深粉
                '#C2185B', // 酒红

                // 紫色系 - 加深
                '#AB47BC', // 紫罗兰
                '#8E24AA', // 深紫
                '#7B1FA2', // 皇家紫

                // 蓝色系 - 加深
                '#42A5F5', // 亮蓝
                '#1E88E5', // 钴蓝
                '#1976D2', // 宝蓝
                '#0097A7', // 青蓝

                // 绿色系 - 加深
                '#66BB6A', // 草绿
                '#43A047', // 森林绿
                '#388E3C', // 深绿
                '#00897B', // 翡翠绿

                // 暖色调 - 加深
                '#FFA726', // 橙黄
                '#FB8C00', // 橙色
                '#F57C00', // 深橙

                // 其他清新色调 - 加深
                '#7E57C2', // 薰衣草紫
                '#26C6DA', // 湖蓝
                '#26A69A', // 碧绿
                '#9CCC65', // 黄绿
                '#FFCA28', // 金黄
                '#8D6E63', // 咖啡棕

                // 更多加深的清新色调
                '#5C6BC0', // 矢车菊蓝
                '#EF5350', // 珊瑚红
                '#78909C', // 石板灰（柔和的深色）
                '#7CB342', // 苹果绿
                '#FFB300', // 琥珀色
                '#546E7A'  // 钢蓝
            ]

            // 复杂的哈希算法确保均匀分布
            let hash = 0
            for (let i = 0; i < user.length; i++) {
                hash = user.charCodeAt(i) + ((hash << 5) - hash) + (i * 31)
                // 防止溢出，限制在合理范围内
                hash = hash & 0x7FFFFFFF
            }

            const index = Math.abs(hash % colors.length)
            return colors[index]
        }
    }
}
</script>

<style scoped>
.avatar {
    position: relative;
    transition: all 0.2s ease;
}

.avatar__image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 50%;
}

.avatar__text {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>