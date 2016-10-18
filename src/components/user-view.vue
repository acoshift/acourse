<template>
  <div>
    <div class="ui segment">
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
  import { User, Course, Loader } from '../services'
  import UserProfile from './user-profile'
  import CourseCard from './course-card'
  import isEmpty from 'lodash/fp/isEmpty'
  import filter from 'lodash/fp/filter'
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
        courses: null,
        $user: null,
        $ownCourse: null
      }
    },
    created () {
      this.init()
    },
    destroyed () {
      this.cleanup()
    },
    watch: {
      $route () {
        this.init()
      }
    },
    methods: {
      init () {
        this.cleanup()
        Loader.start('user')
        this.$user = User.get(this.$route.params.id)
          .flatMap((user) =>
            User.courses(this.$route.params.id)
              .flatMap((courseIds) =>
                isEmpty(courseIds)
                  ? Observable.of([])
                  : Observable.combineLatest(...courseIds.map((id) => Course.get(id)))
                    .map(filter((course) => !!course.public))
              ),
            (user, courses) => ([user, courses])
          )
          .subscribe(
            ([user, courses]) => {
              Loader.stop('user')
              this.user = user
              this.courses = isEmpty(courses) ? null : courses
            },
            () => {
              this.$router.replace('/home')
            }
          )
        this.$ownCourse = User.ownCourses(this.$route.params.id)
          .subscribe(
            (courses) => {
              this.ownCourses = isEmpty(courses) ? null : courses
            }
          )
      },
      cleanup () {
        if (this.$user) this.$user.unsubscribe()
        if (this.$ownCourse) this.$ownCourse.unsubscribe()
      }
    }
  }
</script>
