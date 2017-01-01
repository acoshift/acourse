<template lang="pug">
  .ui.segment
    .ui.stackable.equal.width.grid
      .column
        router-link.ui.green.fluid.button(:to='`/course/${course.id}/edit`') Edit
      .column(v-if='course.canAttend && !course.attend')
        .ui.teal.fluid.button(@click='openAttendModal') Open Attend
      .column(v-if='course.canAttend && course.attend')
        .ui.red.fluid.button(@click='closeAttend', :class='{loading: removingCode}') Close Attend
      .column(v-if='course.hasAssignment')
        router-link.ui.blue.fluid.button(:to='`/course/${course.id}/assignment/edit`') Assignments
      .column(v-if='course.canAttend')
        router-link.ui.blue.fluid.button(:to='`/course/${course.id}/attend`') Attendants
    .ui.small.modal(ref='attendModal')
      .header
        | Set Attend Code
      .content
        .ui.form
          .field
            label Enter Code
            input(v-model='attendCode')
          .ui.red.message(v-if='attendError') {{ attendError }}
          .ui.fluid.blue.button(@click='submitAttend', :class='{loading: submitingAttendCode}') OK
</template>

<script>
import { Course } from 'services'
import moment from 'moment'

export default {
  props: {
    course: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      attendCode: '',
      attendError: '',
      submitingAttendCode: false,
      removingCode: false
    }
  },
  methods: {
    openAttendModal () {
      this.attendError = ''
      this.attendCode = moment().format('DDMMYYYY')
      $(this.$refs.attendModal).modal('show')
    },
    submitAttend () {
      this.attendError = ''
      this.submitingAttendCode = true
      Course.setAttendCode(this.course.id, this.attendCode)
        .finally(() => { this.submitingAttendCode = false })
        .subscribe(
          () => {
            this.attendCode = ''
            $(this.$refs.attendModal).modal('hide')
          },
          (err) => {
            this.attendError = err.message
          }
        )
    },
    closeAttend () {
      this.removingCode = true
      Course.removeAttendCode(this.course.id)
        .finally(() => { this.removingCode = false })
        .subscribe()
    }
  }
}
</script>
