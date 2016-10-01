<template>
  <div>
    <div class="ui segment" :class="{loading: !user}">
      <user-profile :user="user" v-show="user"></user-profile>
    </div>
    <div class="ui segment" v-if="ownCourses">
      <h3 class="ui header">Courses own by {{ user && user.name || 'Anonymous' }}</h3>
      <div class="ui four stackable cards" v-if="ownCourses">
        <course-card v-for="x in ownCourses" :course="x"></course-card>
      </div>
    </div>
    <div class="ui segment" v-if="courses">
      <h3 class="ui header">{{ user && user.name || 'Anonymous' }}'s Courses</h3>
      <div class="ui four stackable cards">
        <course-card v-for="x in courses" :course="x"></course-card>
      </div>
    </div>
  </div>
</template>

<script>
  import { User, Course } from '../services'
  import UserProfile from './user-profile'
  import CourseCard from './course-card'
  import { keys, isEmpty } from 'lodash'
  import { Observable } from 'rxjs'

  export default {
    components: {
      UserProfile,
      CourseCard
    },
    data () {
      return {
        user: null,
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
        User.get(this.$route.params.id)
          .subscribe(
            (user) => {
              this.user = user
              Observable.of(user.course)
                .map(keys)
                .flatMap(Observable.from)
                .flatMap((id) => Course.get(id).first())
                .filter((course) => course.open)
                .toArray()
                .subscribe(
                  (courses) => {
                    this.courses = courses
                  },
                  () => {
                    this.courses = null
                  }
                )
            }
          )
        Course.ownBy(this.$route.params.id)
          .subscribe(
            (courses) => {
              this.ownCourses = isEmpty(courses) ? null : courses
            }
          )
      }
    }
  }
</script>
