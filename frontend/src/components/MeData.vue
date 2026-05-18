<template>
<div class="data">{{ content }}</div>
</template>

<script>
export default {
    props: {
        content: {
            default: ''
        },
    },
}
</script>

<style scoped>
.data {
    white-space: pre-wrap;
    word-break: break-word;
    font-size: 18px;
    color: var(--color-primary);
    padding: 20px 0;
    line-height: 1.2;
}
</style>