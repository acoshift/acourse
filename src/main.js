import Vue from 'vue'
import VueRouter from './vue-router'
import App from './app'

import '!style!css!semantic-ui-css/semantic.min.css'
import '!script!jquery/dist/jquery.min.js'
import '!script!semantic-ui-css/semantic.min.js'

import {
  Firebase
} from './services'

Vue.use(VueRouter)

Firebase.init()

const router = new VueRouter({
  mode: 'history',
  scrollBehavior (to, from, savedPosition) {
    return { x: 0, y: 0 }
  },
  routes: []
})

/* eslint-disable no-new */
new Vue({
  router,
  render: (h) => h(App)
}).$mount('app')
