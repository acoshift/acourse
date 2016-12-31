<template lang="pug">
  .ui.segment(style='padding-bottom: 2rem;')
    h3.ui.header
      | Attendants
      span(v-if='students') ({{ students.length }})
    .ui.stackable.three.column.grid
      .column(v-for='x in students')
        span(:to='`/user/${x.id}`')
          avatar(:src='x.photo', size='tiny')
          | &nbsp;{{ x.name || 'Anonymous' }} ({{ x.count }})
</template>

<script>
import { Course, User, Loader } from 'services'
import Avatar from './Avatar'
import forEach from 'lodash/fp/forEach'

export default {
  components: {
    Avatar
  },
  data () {
    return {
      courseId: this.$route.params.id
    }
  },
  subscriptions () {
    Loader.start('attend-course')
    Loader.start('attend-students')
    return {
      course: Course.get(this.courseId).do(() => Loader.stop('attend-course')).catch(() => { this.$router.replace(`/course/${this.courseId}`) }),
      students: Course.attendUsers(this.courseId).do(forEach(User.inject.bind(User))).do(() => Loader.stop('attend-students'))
    }
  }
}
</script>
