<template>
  <div class="ui basic segment" :class="{loading}">
    <div v-if="courses">
      <h1 class="text-center">All Courses</h1>
      <div class="ui three stackable cards">
        <course-card v-for="x in courses" :course="x"></course-card>
      </div>
    </div>
    <div v-if="!courses && !loading">
      <div class="ui message">No course available yet!</div>
    </div>
  </div>
</template>

<style>
  h1 {
    padding-bottom: 20px;
  }
</style>

<script>
  import CourseCard from './course-card'
  import { Course } from '../services'

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
      Course.list()
        .subscribe(
          (courses) => {
            this.loading = false
            if (courses.length === 0) {
              this.courses = null
            } else {
              this.courses = courses
            }
          }
        )
    }
  }
</script>
