<template>
  <div>
    <div class="ui massive breadcrumb">
      <router-link class="section" to="/course">My Courses</router-link>
      <i class="right chevron icon divider"></i>
      <div class="active section">{{ courseId }}</div>
    </div>
    <div class="ui segment" v-if="course">
      <div class="ui center aligned grid">
        <div class="row">
          <div class="column">
            <img :src="course.photo" class="ui centered big image">
          </div>
        </div>
        <div class="row">
          <div class="column">
            <h1>{{ course.title }}</h1>
          </div>
        </div>
        <div class="row">
          <div class="column">
            <p>{{ course.description }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import { Course } from '../services'

  export default {
    data () {
      return {
        courseId: '',
        course: null
      }
    },
    created () {
      this.init()
    },
    watch: {
      $route () {
        this.init()
      }
    },
    methods: {
      init () {
        this.courseId = this.$route.params.id
        Course.get(this.courseId)
          .subscribe(
            (course) => {
              this.course = course
            },
            () => {
              // not found
              this.$router.replace('/home')
            }
          )
      }
    }
  }
</script>
