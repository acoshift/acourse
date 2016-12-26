<template>
  <div class="ui basic segment" :class="{loading: courses === false}">
    <div v-if="courses && courses.length">
      <h1 class="text-center">All Courses</h1>
      <div class="ui three stackable cards">
        <course-card v-for="x in courses" :course="x"></course-card>
      </div>
    </div>
    <div v-if="courses && !courses.length">
      <div class="ui message">No course available yet!</div>
    </div>
  </div>
</template>

<style scoped>
  .basic.segment {
    min-height: 200px;
  }
</style>

<script>
import { mapState, mapActions } from 'vuex'
import CourseCard from './CourseCard'

export default {
  components: {
    CourseCard
  },
  computed: {
    ...mapState(['courses'])
  },
  created () {
    this.fetchCourses()
  },
  methods: {
    ...mapActions(['fetchCourses'])
  }
}
</script>
