<template>
  <div class="ui basic segment" :class="{loading}">
    <div class="ui massive breadcrumb">
      <div class="active section">My Courses</div>
    </div>
    <div class="ui segment">
      <h3 class="header">My Courses</h3>
      <router-link class="ui blue button" to="/course/new">Create new course</router-link>
      <div class="ui four stackable cards">
        <course-card v-for="x in courses" :course="x"></course-card>
      </div>
    </div>
  </div>
</template>

<script>
  import { User, Course } from '../services'
  import _ from 'lodash'
  import { Observable } from 'rxjs'
  import CourseCard from './course-card'

  export default {
    components: {
      CourseCard
    },
    data () {
      return {
        courses: null,
        loading: false
      }
    },
    created () {
      this.loading = true
      User.me()
        .map((user) => user.course)
        .flatMap((course) => _.keys(course))
        .flatMap((courses) => _.isArray(courses) ? Observable.from(courses) : Observable.of(courses))
        .flatMap(Course.get, (id, course) => ({id, ...course}))
        .toArray()
        .subscribe(
          (courses) => {
            this.loading = false
            this.courses = courses
          },
          () => {
            this.loading = false
          }
        )
    }
  }
</script>
