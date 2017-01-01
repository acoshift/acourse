<template lang="pug">
  .ui.segment
    .ui.embed(ref='video', data-source='youtube', :data-id='src')
</template>

<script>
export default {
  props: ['src'],
  mounted () {
    this.$nextTick(() => {
      $(this.$refs.video).embed()
    })
  }
}
</script>
