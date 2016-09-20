<template>
  <div>
    <div class="ui segment" :class="{loading}">
      <user-profile :user="user" v-show="!loading"></user-profile>
      <div class="ui right aligned basic segment">
        <router-link class="ui green edit button" to="/profile/edit">Edit</router-link>
      </div>
    </div>
    <div class="ui segment">
      <h3 class="ui header">My Courses</h3>
      <router-link class="ui blue button" to="/course/new">Create new course</router-link>
      <div class="four stackable cards" v-if="user && user.courses">
        <course-card v-for="x in user.courses" :course="x"></course-card>
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
  import { User, Course } from '../services'
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
        loading: false
      }
    },
    created () {
      this.loading = true
      User.me()
        .flatMap((user) =>
          Observable.of(user.course)
            .map(_.keys)
            .flatMap(Observable.from)
            .flatMap(Course.get, (id, course) => ({id, ...course}))
            .first()
            .toArray(),
          (user, courses) => ({...user, courses})
        )
        .subscribe(
          (user) => {
            this.loading = false
            this.user = !_.isEmpty(user) ? user : null
          },
          () => {
            this.loading = false
          }
        )
    },
    methods: {
    }
  }
</script>
