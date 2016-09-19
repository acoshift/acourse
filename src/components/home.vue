<template>
  <div class="ui basic segment">
    <div v-if="courses">
      <h1 class="text-center">All Courses</h1>
      <div class="ui four stackable cards">
        <course-card v-for="x in courses" :course="x"></course-card>
      </div>
    </div>
    <div v-else>
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
        courses: null
      }
    },
    created () {
      Course.list()
        .subscribe(
          (courses) => {
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
