import Vue from 'vue'
import VueRouter from 'vue-router'
import App from './app'
import moment from 'moment'

import '!style!css!semantic-ui-css/semantic.min.css'
import '!script!jquery/dist/jquery.min.js'
import '!script!semantic-ui-css/semantic.min.js'
import '!style!css!./style.css'

import {
  Firebase,
  Auth as AuthService
} from './services'

import {
  Auth,
  Layout,
  Home,
  Profile,
  ProfileEdit,
  CourseEditor,
  CourseView,
  UserView,
  CourseChat,
  CourseAttend,
  CourseAssignment
} from './components'

Vue.use(VueRouter)

Firebase.init()

const time = new Vue({
  data () {
    return {
      now: null
    }
  },
  created () {
    setInterval(() => {
      this.now = Date.now()
    }, 60000)
  }
})

Vue.filter('date', (value, input) => {
  time.now
  if (!value) return '-'
  return moment(value).format(input)
})

Vue.filter('fromNow', (value) => {
  time.now
  if (!value) return '-'
  return moment(value).fromNow()
})

Vue.filter('trim', (value, input) => {
  value = value || ''
  if (value.length <= input) return value
  return value.substr(0, input) + '...'
})

const router = new VueRouter({
  mode: 'history',
  scrollBehavior (to, from, savedPosition) {
    return { x: 0, y: 0 }
  },
  routes: [
    { path: '/', component: Auth, beforeEnter: redirectIfAuth },
    {
      path: '/home',
      component: Layout,
      children: [
        { path: '', component: Home },
        { path: '/profile', component: Profile },
        { path: '/profile/edit', component: ProfileEdit },
        { path: '/course/new', component: CourseEditor },
        { path: '/course/:id', component: CourseView },
        { path: '/course/:id/chat', component: CourseChat },
        { path: '/course/:id/edit', component: CourseEditor },
        { path: '/course/:id/attend', component: CourseAttend },
        { path: '/course/:id/assignment', component: CourseAssignment },
        { path: '/user/:id', component: UserView }
      ],
      beforeEnter: redirectIfNotAuth
    },
    { path: '*', redirect: '/' }
  ]
})

router.beforeEach((to, from, next) => {
  window.ga('set', 'page', to.path)
  window.ga('send', 'pageview')
  next()
})

function redirectIfAuth (to, from, next) {
  AuthService.currentUser
    .first()
    .subscribe(
      (user) => {
        if (user) {
          next('/home')
        } else {
          next()
        }
      }
    )
}

function redirectIfNotAuth (to, from, next) {
  AuthService.currentUser
    .first()
    .subscribe(
      (user) => {
        if (user) {
          next()
        } else {
          next('/')
        }
      }
    )
}

/* eslint-disable no-new */
new Vue({
  router,
  render: (h) => h(App)
}).$mount('app')
