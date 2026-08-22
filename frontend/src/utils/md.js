import { marked } from 'marked'

marked.use({
    breaks: false,
    gfm: true
})

export function getMd() {
    return marked
}