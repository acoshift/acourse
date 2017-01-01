<template lang="pug">
  .ui.segment
    EmbedVideo(:src='src')
</template>

<script>
import EmbedVideo from './EmbedVideo'

export default {
  components: {
    EmbedVideo
  },
  props: ['src']
}
</script>
