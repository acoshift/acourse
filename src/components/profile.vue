<template>
  <div>
    <div class="ui segment" :class="{loading}">
      <user-profile :user="user" v-show="!loading"></user-profile>
      <div class="ui right aligned basic segment">
        <router-link class="ui green edit button" to="/profile/edit">Edit</router-link>
      </div>
    </div>
    <div class="ui segment" v-if="user && user.instructor">
      <h3 class="ui header">My Own Courses</h3>
      <router-link class="ui blue button" to="/course/new">Create new course</router-link>
      <div class="ui four stackable cards" v-if="ownCourses">
        <course-card v-for="x in ownCourses" :course="x"></course-card>
      </div>
    </div>
    <div class="ui segment" v-if="courses">
      <h3 class="ui header">My Courses</h3>
      <div class="ui four stackable cards">
        <course-card v-for="x in courses" :course="x"></course-card>
      </div>
    </div>
  </div>
</template>

<style>
  .cards {
    padding-top: 30px;
  }

  .edit.button {
    width: 140px;
  }
</style>

<script>
  import { Auth, User, Course } from '../services'
  import UserProfile from './user-profile'
  import CourseCard from './course-card'
  import _ from 'lodash'
  import { Observable } from 'rxjs'

  export default {
    components: {
      UserProfile,
      CourseCard
    },
    data () {
      return {
        user: null,
        loading: false,
        ownCourses: null,
        courses: null,
        ob: []
      }
    },
    created () {
      this.loading = true
      this.ob.push(User.me()
        .subscribe(
          (user) => {
            this.loading = false
            this.user = user
          },
          () => {
            this.loading = false
          }
        )
      )
      Auth.currentUser()
        .first()
        .flatMap((user) => Course.ownBy(user.uid))
        .subscribe(
          (courses) => {
            this.ownCourses = _.isEmpty(courses) ? null : courses
          }
        )
      User.me()
        .first()
        .map((user) => user.course)
        .map(_.keys)
        .flatMap(Observable.from)
        .flatMap((id) => Course.get(id).first())
        .toArray()
        .subscribe(
          (courses) => {
            this.courses = courses
          },
          () => {
            this.courses = null
          }
        )
    },
    destroyed () {
      _.forEach(this.ob, (x) => x.unsubscribe())
    }
  }
</script>
