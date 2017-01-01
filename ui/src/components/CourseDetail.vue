<template lang="pug">
  .ui.segment
    h4.ui.header Course Detail
    span Start: {{ course.start | date('DD/MM/YYYY') }}
    .ui.divider
    p.description(v-html='course.description')
</template>

<script>
export default {
  props: {
    course: {
      type: Object,
      required: true
    }
  }
}
</script>
