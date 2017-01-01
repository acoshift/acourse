<template lang="pug">
  .ui.segment(v-if='course.attend || course.assignment')
    .ui.stackable.equal.width.grid
      .column(v-if='course.attend')
        .ui.blue.fluid.button(:class='{disabled: isAttended || !course.attend, loading: attending}', @click='attend') Attend
      .column(v-if='course.assignment')
        router-link.ui.teal.fluid.button(:to='`/course/${course.id}/assignment`') Assignments
</template>

<script>
import { Me, Document } from 'services'

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
