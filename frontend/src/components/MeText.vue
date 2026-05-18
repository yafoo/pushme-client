<template>
<div class="text">{{ content }}</div>
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
.text {
    white-space: pre-wrap;
    word-break: break-word;
}
</style>