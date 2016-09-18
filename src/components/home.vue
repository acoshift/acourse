<template>
  <div class="ui basic segment">
    <div v-if="courses" class="ui five stackable cards">
      <course-card v-for="x in courses"></course-card>
    </div>
    <div v-else>
      <div class="ui message">No course available yet!</div>
    </div>
  </div>
</template>

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
