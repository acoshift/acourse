import Vue from 'vue'
import VueRouter from 'vue-router'
import VueRxJS from './vue-rxjs'
import App from './app'
import './filters'

import '!style!css!semantic-ui-css/semantic.min.css'
import '!script!jquery/dist/jquery.min.js'
import '!script!semantic-ui-css/semantic.min.js'
import '!style!css!./style.css'
// import serviceWorker from 'serviceworker!./sw.js'

// serviceWorker({ scope: '/' }).then((reg) => {
//   reg.pushManager.subscribe({
//     userVisibleOnly: true
//   }).then((sub) => {
//     console.log('endpoint:', sub.endpoint)
//   })
// })

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
