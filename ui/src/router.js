import Vue from 'vue'
import VueRouter from 'vue-router'

import {
  Auth,
  Loader,
  Document,
  Me
} from 'services'

import {
  Layout,
  Home,
  Profile,
  ProfileEdit,
  Course,
  CourseEditor,
  CourseView,
  // CourseAttend,
  CourseAssignments,
  // CourseAssignmentEdit,
  AdminCourse,
  AdminPayment,
  AdminPaymentHistory,
  Privacy
} from 'components'

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
        { path: '/course/new', component: CourseEditor, name: 'courseNew', beforeEnter: isRole('instructor') },
        {
          path: '/course/:id',
          component: Course,
          children: [
            { path: '', component: CourseView, name: 'courseView' },
            { path: 'edit', component: CourseEditor, name: 'courseEdit', beforeEnter: redirectIfNotAuth },
            // { path: 'attend', component: CourseAttend, name: 'courseAttend', beforeEnter: redirectIfNotAuth },
            { path: 'assignments', component: CourseAssignments, name: 'courseAssignments' }
            // { path: 'assignment/edit', component: CourseAssignmentEdit, name: 'courseAssignmentEdit', beforeEnter: redirectIfNotAuth }
          ]
        },
        { path: '/admin/course', component: AdminCourse, beforeEnter: isRole('admin') },
        { path: '/admin/payment', component: AdminPayment, beforeEnter: isRole('admin') },
        { path: '/admin/payment/history', component: AdminPaymentHistory, beforeEnter: isRole('admin') },
        { path: '/privacy', component: Privacy }
      ]
    },
    { path: '*', redirect: '/' }
  ]
})

router.afterEach((to) => {
  Document.reset()
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
//           next('/')
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

function isRole (role) {
  return (to, from, next) => {
    Loader.start('router')
    Me.get()
      .subscribe(
        (user) => {
          Loader.stop('router')
          if (user && user.role[role]) {
            next()
          } else {
            next('/')
          }
        }
      )
  }
}

export default router
