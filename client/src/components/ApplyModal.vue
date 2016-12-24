<template lang="pug">
  .ui.small.modal
    .header Enroll
    .content
      div(v-if="course.canQueueEnroll")
        h4 Direct Transfer
        p.description {{ course.queueEnrollDetail }}
        .ui.fluid.green.button(@click="upload") Upload Slip
      h4.ui.horizontal.divider.header OR
      form.ui.form(@submit.prevent="apply")
        .field
          label Enter Code
          input(v-model="code")
        button.ui.fluid.submit.blue.button(:class="{loading: applying}") Enroll
</template>

<script>
import { Document, Me } from '../services'

export default {
  props: {
    course: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      code: '',
      applying: false
    }
  },
  methods: {
    show () {
      this.code = ''
      $(this.$el).modal('show')
    },
    apply () {
      $(this.$el).modal('hide')
      if (this.applying) return
      this.applying = true
      Me.applyCourse(this.course.id, this.code)
        .finally(() => { this.applying = false })
        .subscribe(
          () => {
            Document.openSuccessModal('Success', 'You have enrolled to this course.')
          },
          (err) => {
            Document.openErrorModal('Error', 'Can not enroll to this course. ' + err.message)
          }
        )
    },
    upload () {
      Document.uploadModal.open('image/*')
        .flatMap((file) => Me.submitCourseQueueEnroll(this.course.id, file.downloadURL))
        .subscribe(
          () => {
            Document.openSuccessModal('Success', 'Your enroll request success!.')
          },
          (err) => {
            Document.openErrorModal('Upload Error', err && err.message || err)
          }
        )
    }
  }
}
</script>
