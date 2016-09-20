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
  Course,
  CourseEditor,
  CourseView,
  UserView
} from './components'

Vue.use(VueRouter)

Firebase.init()

Vue.filter('date', (value, input) => {
  if (!value) return '-'
  return moment(value).format(input)
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
        { path: '/course', component: Course },
        { path: '/course/new', component: CourseEditor },
        { path: '/course/:id', component: CourseView },
        { path: '/course/:id/edit', component: CourseEditor },
        { path: '/user/:id', component: UserView }
      ],
      beforeEnter: redirectIfNotAuth
    },
    { path: '*', redirect: '/' }
  ]
})

function redirectIfAuth (to, redirect, next) {
  AuthService.currentUser
    .first()
    .subscribe(
      (user) => {
        if (user) {
          redirect('/home')
        } else {
          next()
        }
      }
    )
}

function redirectIfNotAuth (to, redirect, next) {
  AuthService.currentUser
    .first()
    .subscribe(
      (user) => {
        if (user) {
          next()
        } else {
          redirect('/')
        }
      }
    )
}

/* eslint-disable no-new */
new Vue({
  router,
  render: (h) => h(App)
}).$mount('app')
