<template>
  <div>
    <div class="ui segment" style="padding-bottom: 2rem;">
      <h3 class="ui header">Attendants <span v-if="students">({{ students.length }})</span></h3>
      <div class="ui stackable three column grid">
        <div class="column" v-for="x in students">
          <span :to="`/user/${x.id}`">
            <avatar :src="x.photo" size="tiny"></avatar>
            {{ x.name || 'Anonymous' }} ({{ x.count }})
          </span>
        </div>
      </div>
    </div>
  </div>
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
