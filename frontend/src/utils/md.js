import markdownit from 'markdown-it'
let md = null

export function getMd() {
    if(!md) {
        md = markdownit({html: true, linkify: false})
    }
    return md
}