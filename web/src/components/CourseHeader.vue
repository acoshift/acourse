<template lang="pug">
  .ui.basic.segment
    .ui.center.aligned.grid
      .row
        .column
          img.ui.centered.big.image(:src='course.photo')
      .row(style='padding-top: 0;')
        .column
          h1 {{ course.title }}
      .two.column.middle.aligned.row
        .right.aligned.column(style='padding-right: 2px;')
          span(:to='`/user/${course.owner.id}`')
            avatar(:src='course.owner.photo', size='mini')
        .left.aligned.column(style='padding-left: 2px;')
          span(:to='`/user/${course.owner.id}`')
            h3 {{ course.owner.name || 'Anonymous' }}
</template>

<style scoped>
  .segment {
    padding-top: 0;
  }
</style>

<script>
import Avatar from './Avatar'

export default {
  components: {
    Avatar
  },
  props: {
    course: {
      type: Object,
      required: true
    }
  }
}
</script>
