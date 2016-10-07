<template>
  <div>
    <div class="ui segment" style="padding-bottom: 2rem;">
      <h3 class="ui header">Attendants <span v-if="students">({{ students.length }})</span></h3>
      <div class="ui stackable three column grid">
        <div class="column" v-for="x in students">
          <router-link :to="`/user/${x.id}`">
            <avatar :src="x.photo" size="tiny"></avatar>
            {{ x.name || 'Anonymous' }} ({{ x.count }})
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import { Course, Loader } from '../services'
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
        students: null,
        $attend: null
      }
    },
    beforeCreate () {
      Loader.start('attend')
    },
    created () {
      this.courseId = this.$route.params.id

      this.$attend = Observable.combineLatest(
        Course.get(this.courseId),
        Course.attendUsers(this.courseId)
      )
        .subscribe(
          ([course, students]) => {
            Loader.stop('attend')
            this.course = course
            this.students = students
          },
          () => {
            this.$router.replace(`/course/${this.courseId}`)
          }
        )
    },
    destroyed () {
      this.$attend.unsubscribe()
    }
  }
</script>
