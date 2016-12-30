<template>
  <div class="ui basic segment" :class="{loading: !course}">
    <div class="ui huge breadcrumb" style="padding-bottom: 1.5rem;">
      <router-link class="section" to="/home">Courses</router-link>
      <i class="right chevron icon divider"></i>
      <router-link :to="`/course/${courseId}`" :tag="$route.name === 'courseView' && 'div' || 'a'" class="section" active-class="active" exact>{{ course && course.title || courseId }}</router-link>
      <i v-show="$route.name !== 'courseView'" class="right chevron icon divider"></i>
      <div v-show="$route.name === 'courseEdit'" class="active section">Edit</div>
      <div v-show="$route.name === 'courseNew'" class="active section">New</div>
      <div v-show="$route.name === 'courseAssignment'" class="active section">Assignments</div>
      <div v-show="$route.name === 'courseAttend'" class="active section">Attendants</div>
      <div v-show="$route.name === 'courseAssignmentEdit'" class="active section">Edit Assignment</div>
    </div>
    <router-view></router-view>
  </div>
</template>

<style scoped>
  @media only screen and (max-width: 500px) {
    .breadcrumb {
      font-size: 1.05rem !important;
    }
  }
</style>

<script>
import { Course, Document } from 'services'

export default {
  data () {
    return {
      courseId: this.$route.params.id
    }
  },
  subscriptions () {
    return {
      course: Course.get(this.courseId)
        .do((course) => Document.setCourse(course))
    }
  },
  created () {
    Course.fetch(this.courseId)
      .subscribe(null, () => {
        this.$router.replace('/')
      })
  }
}
</script>
