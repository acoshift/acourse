<template lang="pug">
  .ui.segment
    .ui.stackable.equal.width.grid
      .center.aligned.column
        .ui.blue.button(style="width: 200px;", :class="{loading: applying}", @click="apply") Enroll
    apply-modal(ref="applyModal", :course="course")
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
            Document.openSuccessModal('Success', 'You have enrolled to this course.')
          },
          () => {
            this.$refs.applyModal.show()
          }
        )
    }
  }
}
</script>
