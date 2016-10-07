<template>
  <div class="ui basic segment">
    <div class="ui container"></div>
    <div v-if="courses">
      <div v-if="courses.length">
        <h1 class="text-center">All Courses</h1>
        <div class="ui three stackable cards">
          <course-card v-for="x in courses" :course="x"></course-card>
        </div>
      </div>
      <div v-else>
        <div class="ui message">No course available yet!</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .basic.segment {
    min-height: 200px;
  }

  h1 {
    padding-bottom: 20px;
  }
</style>

<script>
  import CourseCard from './course-card'
  import { Course, Loader } from '../services'

  export default {
    components: {
      CourseCard
    },
    data () {
      return {
        courses: Course.list()
          .do(() => Loader.stop('courses'))
      }
    },
    beforeCreate () {
      Loader.start('courses')
    }
  }
</script>
