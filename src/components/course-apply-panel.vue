<template>
  <div class="ui segment">
    <div class="ui blue button" style="width: 180px;" :class="{loading: applying}" @click="apply">Apply</div>
  </div>
</template>

<script>
  import { Course, Document } from '../services'

  export default {
    props: ['course'],
    data () {
      return {
        applying: false
      }
    },
    methods: {
      apply () {
        if (this.applying) return
        this.applying = true
        Course.join(this.course.id)
          .finally(() => { this.applying = false })
          .subscribe(
            () => {
              Document.openSuccessModal('Success', 'You have applied to this course.')
            }
          )
      }
    }
  }
</script>
