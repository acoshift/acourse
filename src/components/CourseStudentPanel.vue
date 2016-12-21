<template>
  <div class="ui segment">
    <div class="ui stackable equal width grid">
      <div class="column" v-if="course.canAttend">
        <div :class="{disabled: isAttended || !course.attend, loading: attending}" class="ui blue fluid button" @click="attend">Attend</div>
      </div>
      <div class="column" v-if="course.hasAssignment">
        <router-link class="ui teal fluid button" :to="`/course/${course.id}/assignment`">Assignments</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { Me, Document } from '../services'

export default {
  props: {
    course: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      attending: false
    }
  },
  subscriptions () {
    return {
      isAttended: Me.isAttendedCourse(this.course.id)
    }
  },
  methods: {
    attend () {
      this.attending = true
      Me.attendCourse(this.course.id)
        .finally(() => { this.attending = false })
        .subscribe(
          () => {
            Document.openSuccessModal('Success', 'You have attended to this section.')
          },
          (err) => {
            Document.openErrorModal('Attend Error', err && err.message || err)
          }
        )
    }
  }
}
</script>
