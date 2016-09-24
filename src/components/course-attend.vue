<template>
  <div class="ui basic segment" :class="{loading}">
    <div class="ui massive breadcrumb">
      <router-link class="section" to="/home">Courses</router-link>
      <i class="right chevron icon divider"></i>
      <router-link class="section" :to="`/course/${courseId}`">{{ course && course.title || courseId }}</router-link>
      <i class="right chevron icon divider"></i>
      <div class="active section">Attendants</div>
    </div>
    <div class="ui segment">
      <h3 class="ui header">Attendants <span v-if="students">({{ students.length }})</span></h3>
      <div v-for="x in students">
        <router-link :to="`/user/${x.id}`">
          <avatar :src="x.photo" size="tiny"></avatar>
          {{ x.name }} ({{ x.count }})
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
  import { Course } from '../services'
  import { Observable } from 'rxjs'
  import Avatar from './avatar'

  export default {
    components: {
      Avatar
    },
    data () {
      return {
        course: null,
        courseId: null,
        loading: false,
        students: null
      }
    },
    created () {
      this.loading = true
      this.courseId = this.$route.params.id

      Observable.combineLatest(
        Course.get(this.courseId),
        Course.attendUsers(this.courseId)
      )
        .subscribe(
          ([course, students]) => {
            this.loading = false
            this.course = course
            this.students = students
          },
          () => {
            this.loading = false
            this.$router.replace('/home')
          }
        )
    }
  }
</script>
