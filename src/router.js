import Vue from 'vue'
import VueRouter from 'vue-router'

import {
  Auth as AuthService,
  Loader,
  Document
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
  CourseChatHistory,
  CourseAttend,
  CourseAssignment,
  CourseAssignmentEdit
} from './components'

Vue.use(VueRouter)

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
            { path: 'chat/history', component: CourseChatHistory, name: 'courseChatHistory' },
            { path: 'edit', component: CourseEditor, name: 'courseEdit' },
            { path: 'attend', component: CourseAttend, name: 'courseAttend' },
            { path: 'assignment', component: CourseAssignment, name: 'courseAssignment' },
            { path: 'assignment/edit', component: CourseAssignmentEdit, name: 'courseAssignmentEdit' }
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
  Document.setTitle()
  Loader.reset()
  window.ga('set', 'page', to.path)
  window.ga('send', 'pageview')
})

function redirectIfAuth (to, from, next) {
  Loader.start('router')
  AuthService.currentUser()
    .first()
    .subscribe(
      (user) => {
        Loader.stop('router')
        if (user) {
          next('/home')
        } else {
          next()
        }
      }
    )
}

function redirectIfNotAuth (to, from, next) {
  Loader.start('router')
  AuthService.currentUser()
    .first()
    .subscribe(
      (user) => {
        Loader.stop('router')
        if (user) {
          next()
        } else {
          next('/')
        }
      }
    )
}

export default router
