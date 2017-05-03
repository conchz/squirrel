'use strict'

// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import FastClick from 'fastclick'
import {AlertPlugin, ToastPlugin, AjaxPlugin, WechatPlugin} from 'vux'
import App from './App'
import router from './router'

FastClick.attach(document.body)

// Load plugins
Vue.use(AlertPlugin)
Vue.use(ToastPlugin)
Vue.use(AjaxPlugin)
Vue.use(WechatPlugin)

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  router,
  render: h => h(App)
}).$mount('#app-box')
