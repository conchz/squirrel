'use strict'

import Vue from 'vue'
import Router from 'vue-router'
import Hello from './components/Hello'
import Home from './components/HelloFromVux'
import NotFound from './components/NotFound'

Vue.use(Router)

const routes = [
  {
    path: '/',
    component: Home
  },
  {
    path: '/home',
    component: Hello
  },
  {
    path: '*',
    component: NotFound
  }
]

export default new Router({
  mode: 'history',
  routes: routes
})
