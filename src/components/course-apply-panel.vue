<template>
  <div class="ui segment">
    <div class="ui stackable equal width grid">
      <div class="center aligned column">
        <div class="ui blue button" style="width: 200px;" :class="{loading: applying}" @click="apply">Apply</div>
      </div>
    </div>
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
