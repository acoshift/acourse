import Vue from 'vue'
import Vuex from 'vuex'
import find from 'lodash/find'

import { Auth, Course, Document, User } from 'services'

import router from './router'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    courses: false, // course tiny collection view
    course: {}, // course view
    authUser: false, // firebase user, uninit === false
    currentUser: false,
    ownCourses: false,
    myCourses: false
  },
  getters: {
    currentCourse (state) {
      const id = state.route.params.id
      const course = state.course[id] || find(state.course, { url: id })
      Document.setCourse(course)
      return course
    }
  },
  mutations: {
    updateAuthUser (state, user) {
      state.authUser = user
    },
    updateCurrentUser (state, user) {
      state.currentUser = user
    },
    updateCourses (state, courses) {
      state.courses = courses
    },
    updateOwnCourses (state, courses) {
      state.ownCourses = courses
    },
    updateMyCourses (state, courses) {
      state.myCourses = courses
    },
    updateCourse (state, course) {
      if (!course) return
      state.course = {
        ...state.course,
        [course.id]: course
      }
    }
  },
  actions: {
    patch (ctx, s) {
      if (!s) return
      if (s.courses) ctx.commit('updateCourses', s.courses)
      if (s.course) ctx.commit('updateCourse', s.course)
    },
    authStateChanged (ctx, user) {
      ctx.commit('updateAuthUser', user)
      ctx.dispatch('fetchMe')
    },
    signOut () {
      Auth.signOut().subscribe(() => {
        Vue.nextTick(() => {
          router.push('/')
        })
      })
    },
    fetchMe (ctx) {
      if (ctx.state.authUser) {
        User.get(ctx.state.authUser.uid).subscribe((user) => ctx.commit('updateCurrentUser', user))
      } else {
        ctx.commit('updateCurrentUser', null)
      }
    },
    fetchMeOwnCourses (ctx) {
      if (!ctx.state.authUser) return
      User.ownCourses(ctx.state.authUser.uid).subscribe((courses) => ctx.commit('updateOwnCourses', courses))
    },
    fetchMeMyCourses (ctx) {
      if (!ctx.state.authUser) return
      User.courses(ctx.state.authUser.uid).subscribe((courses) => ctx.commit('updateMyCourses', courses))
    },
    fetchCourses (ctx) {
      Course.list().subscribe((courses) => ctx.commit('updateCourses', courses))
    },
    fetchCourse (ctx, id) {
      Course.get(id)
        .subscribe(
          (course) => {
            ctx.commit('updateCourse', course)
          },
          () => {
            router.replace('/')
          }
        )
    },
    fetchCurrentCourse (ctx) {
      ctx.dispatch('fetchCourse', ctx.state.route.params.id)
    }
  }
})
