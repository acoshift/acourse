import Vue from 'vue'
import VueRouter from 'vue-router'

import {
  Auth,
  Loader,
  Document,
  Me
} from './services'

import {
  Layout,
  Home,
  Profile,
  ProfileEdit,
  Course,
  CourseEditor,
  CourseView,
  CourseAttend,
  CourseAssignment,
  CourseAssignmentEdit,
  QueueEnroll
} from './components'

Vue.use(VueRouter)

const router = new VueRouter({
  mode: 'history',
  scrollBehavior (to, from, savedPosition) {
    return { x: 0, y: 0 }
  },
  routes: [
    {
      path: '/',
      component: Layout,
      children: [
        { path: '', component: Home },
        { path: '/profile', component: Profile, beforeEnter: redirectIfNotAuth },
        { path: '/profile/edit', component: ProfileEdit, beforeEnter: redirectIfNotAuth },
        { path: '/course/new', component: CourseEditor, name: 'courseNew', beforeEnter: isInstructor },
        {
          path: '/course/:id',
          component: Course,
          children: [
            { path: '', component: CourseView, name: 'courseView' },
            { path: 'edit', component: CourseEditor, name: 'courseEdit', beforeEnter: redirectIfNotAuth },
            { path: 'attend', component: CourseAttend, name: 'courseAttend', beforeEnter: redirectIfNotAuth },
            { path: 'assignment', component: CourseAssignment, name: 'courseAssignment' },
            { path: 'assignment/edit', component: CourseAssignmentEdit, name: 'courseAssignmentEdit', beforeEnter: redirectIfNotAuth }
          ]
        },
        { path: '/queue-enroll', component: QueueEnroll, beforeEnter: redirectIfNotAuth }
      ]
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

// function redirectIfAuth (to, from, next) {
//   Loader.start('router')
//   Auth.currentUser()
//     .first()
//     .subscribe(
//       (user) => {
//         Loader.stop('router')
//         if (user) {
//           next('/home')
//         } else {
//           next()
//         }
//       }
//     )
// }

function redirectIfNotAuth (to, from, next) {
  Loader.start('router')
  Auth.currentUser()
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

function isInstructor (to, from, next) {
  Loader.start('router')
  Me.get()
    .subscribe(
      (user) => {
        Loader.stop('router')
        if (user && user.instructor) {
          next()
        } else {
          next('/')
        }
      }
    )
}

export default router
