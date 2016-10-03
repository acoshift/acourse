import Vue from 'vue'
import VueRouter from 'vue-router'
import VueRxJS from './vue-rxjs'
import Raven from 'raven-js'
import App from './app'
import './filters'

import 'semantic-ui-css/components/accordion.min.js'
import 'semantic-ui-css/components/checkbox.min.js'
import 'semantic-ui-css/components/dimmer.min.js'
import 'semantic-ui-css/components/modal.min.js'
import 'semantic-ui-css/components/transition.min.js'

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
  UserView,
  CourseChat,
  CourseAttend,
  CourseAssignment
} from './components'

Raven
  .config('https://fda9f1b21cd04a4585b9f9051b37a466@sentry.io/103020')
  .install()

Vue.use(VueRxJS)
Vue.use(VueRouter)

Firebase.init()

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
        { path: '/course/new', component: CourseEditor, name: 'courseNew' },
        {
          path: '/course/:id',
          component: Course,
          children: [
            { path: '', component: CourseView, name: 'courseView' },
            { path: 'chat', component: CourseChat, name: 'courseChat' },
            { path: 'edit', component: CourseEditor, name: 'courseEdit' },
            { path: 'attend', component: CourseAttend, name: 'courseAttend' },
            { path: 'assignment', component: CourseAssignment, name: 'courseAssignment' }
          ]
        },
        { path: '/user/:id', component: UserView }
      ],
      beforeEnter: redirectIfNotAuth
    },
    { path: '*', redirect: '/' }
  ]
})

router.afterEach((to) => {
  window.ga('set', 'page', to.path)
  window.ga('send', 'pageview')
})

function redirectIfAuth (to, from, next) {
  AuthService.currentUser()
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
  AuthService.currentUser()
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
