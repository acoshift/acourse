<template lang="pug">
  .ui.segment
    .ui.stackable.equal.width.grid
      .center.aligned.column
        .ui.blue.button(style="width: 200px;", :class="{loading: applying}", @click="apply") Apply
    apply-modal(ref="applyModal", @apply="applyWithCode")
</template>

<script>
import { Me, Document } from '../services'
import ApplyModal from './ApplyModal'

export default {
  components: {
    ApplyModal
  },
  props: {
    course: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      applying: false
    }
  },
  methods: {
    apply () {
      if (this.applying) return
      this.applying = true

      Me.applyCourse(this.course.id)
        .finally(() => { this.applying = false })
        .subscribe(
          () => {
            Document.openSuccessModal('Success', 'You have applied to this course.')
          },
          () => {
            this.$refs.applyModal.show()
          }
        )
    },
    applyWithCode (code) {
      if (this.applying) return
      this.applying = true
      Me.applyCourse(this.course.id, code)
        .finally(() => { this.applying = false })
        .subscribe(
          () => {
            Document.openSuccessModal('Success', 'You have applied to this course.')
          },
          (err) => {
            Document.openErrorModal('Error', 'Can not apply to this course. ' + err.message)
          }
        )
    }
  }
}
</script>
