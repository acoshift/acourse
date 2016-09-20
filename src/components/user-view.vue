<template>
  <div>
    <div class="ui segment" :class="{loading}">
      <user-profile :user="user" v-show="!loading"></user-profile>
    </div>
    <div class="ui segment" v-if="ownCourses">
      <h3 class="ui header">Courses own by <span v-if="user">{{ user.name }}</span><span v-else>this user</span></h3>
      <div class="four stackable cards" v-if="ownCourses">
        <course-card v-for="x in ownCourses" :course="x"></course-card>
      </div>
    </div>
    <div class="ui segment" v-if="courses">
      <h3 class="ui header"><span v-if="user">{{ user.name }}'s</span><span v-else>This user's</span> Courses</h3>
      <div class="four stackable cards">
        <course-card v-for="x in courses" :course="x"></course-card>
      </div>
    </div>
  </div>
</template>

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
        loading: false,
        ownCourses: null,
        courses: null
      }
    },
    created () {
      this.init()
    },
    watch: {
      $route () {
        this.init()
      }
    },
    methods: {
      init () {
        this.loading = true
        User.get(this.$route.params.id)
          .subscribe(
            (user) => {
              this.loading = false
              this.user = (user.name && user.photo) ? user : null
              Observable.of(user.course)
                .map(_.keys)
                .flatMap(Observable.from)
                .flatMap(Course.get, (id, course) => ({id, ...course}))
                .first()
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
            () => {
              this.loading = false
            }
          )
        Course.ownBy(this.$route.params.id)
          .subscribe(
            (courses) => {
              this.ownCourses = _.isEmpty(courses) ? null : courses
            }
          )
      }
    }
  }
</script>
