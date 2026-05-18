import { createApp } from 'vue'

const loadPage = async () => {
  const urlParams = new URLSearchParams(window.location.search)
  const pageName = urlParams.get('page') || 'App'
  const module = await import(`./pages/${pageName}.vue`)
  const component = module.default

  createApp(component).mount('#app')
}

loadPage()