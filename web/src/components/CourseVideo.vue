<template lang="pug">
  .ui.segment
    Video(:src='src')
</template>

<script>
import Video from './Video'

export default {
  components: {
    Video
  },
  props: ['src']
}
</script>
