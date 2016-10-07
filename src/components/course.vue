<template>
  <div class="ui basic segment">
    <div class="ui huge breadcrumb" style="padding-bottom: 1.5rem;">
      <router-link class="section" to="/home">Courses</router-link>
      <i class="right chevron icon divider"></i>
      <router-link :to="`/course/${courseId}`" :tag="$route.name === 'courseView' ? 'div' : 'a'" class="section" active-class="active" exact>{{ course && course.title || courseId }}</router-link>
      <i v-show="$route.name !== 'courseView'" class="right chevron icon divider"></i>
      <div v-show="$route.name === 'courseEdit'" class="active section">Edit</div>
      <div v-show="$route.name === 'courseNew'" class="active section">New</div>
      <div v-show="$route.name === 'courseChat'" class="active section">Chat</div>
      <div v-show="$route.name === 'courseAssignment'" class="active section">Assignments</div>
      <div v-show="$route.name === 'courseAttend'" class="active section">Attendants</div>
    </div>
    <router-view></router-view>
  </div>
</template>

<script>
  import { Course, Loader } from '../services'

  export default {
    data () {
      return {
        courseId: null,
        course: null,
        $course: null
      }
    },
    beforeCreate () {
      Loader.start('course')
    },
    created () {
      this.loading = true
      this.courseId = this.$route.params.id
      this.$course = Course.get(this.courseId)
        .do(() => Loader.stop('course'))
        .subscribe(
          (course) => {
            this.course = course
          },
          () => {
            this.$router.replace('/home')
          }
        )
    },
    destroyed () {
      this.$course.unsubscribe()
    }
  }
</script>
