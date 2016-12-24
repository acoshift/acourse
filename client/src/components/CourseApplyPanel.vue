<template lang="pug">
  .ui.segment
    .ui.stackable.equal.width.grid
      .row(v-if="!purchased")
        .center.aligned.column
          span(v-if="price === 0")
            h2 FREE
          span(v-else)
            h2 ฿ {{ price }}
      .row
        .center.aligned.column
          .ui.blue.button(style="width: 200px;", :class="{loading: applying}", @click="apply", v-if="!purchased") Enroll
          .ui.blue.disabled.button(v-else) Wait for Approve
    apply-modal(ref="applyModal", :course="course")
</template>

<script>
import { Course, Document } from '../services'
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

      if (this.price === 0) {
        this.applying = true
        Course.enroll(this.course.id, {})
          .finally(() => { this.applying = false })
          .subscribe(
            () => {
              Document.openSuccessModal('Success', 'You have enrolled to this course.')
              // TODO: need reload course
            }
          )
      } else {
        this.$refs.applyModal.show()
      }
    }
  },
  computed: {
    price () {
      if (this.course.discount) return this.course.discountedPrice
      return this.course.price
    },
    purchased () {
      return this.course.purchaseStatus === 'waiting'
    }
  }
}
</script>
