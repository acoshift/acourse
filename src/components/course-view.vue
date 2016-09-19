<template>
  <div>
    <div class="ui massive breadcrumb">
      <router-link v-if="!isMyCourse" class="section" to="/home">Courses</router-link>
      <router-link v-else class="section" to="/course">My Courses</router-link>
      <i class="right chevron icon divider"></i>
      <div class="active section">{{ course && course.title || courseId }}</div>
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
        <div v-if="isMyCourse" class="right aligned row">
          <div class="column">
            <router-link class="ui green button" :to="`/course/${courseId}/edit`">Edit</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import { User, Course } from '../services'
  import { Observable } from 'rxjs'
  import _ from 'lodash'

  export default {
    data () {
      return {
        courseId: '',
        course: null,
        isMyCourse: false
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
        Observable.forkJoin(
          User.me(),
          Course.get(this.courseId)
        )
          .subscribe(
            ([user, course]) => {
              this.course = course
              if (_.get(user.course, this.courseId)) this.isMyCourse = true
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
