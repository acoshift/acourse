<template lang="pug">
.ui.basic.segment(:class='{loading: courses === false}')
  div(v-if='courses && courses.length')
    h1.text-center All Courses
    .ui.three.stackable.cards
      CourseCard(v-for='x in courses', :course='x')
  div(v-if='courses && !courses.length')
    .ui.message No course available yet!
</template>

<style scoped>
  .basic.segment {
    min-height: 200px;
  }
</style>

<script>
import { Course } from 'services'
import CourseCard from './CourseCard'

export default {
  components: {
    CourseCard
  },
  subscriptions () {
    return {
      courses: Course.list()
    }
  },
  created () {
    this.fetchCourses()
  },
  methods: {
    fetchCourses () {
      Course.fetchList()
    }
  }
}
</script>
