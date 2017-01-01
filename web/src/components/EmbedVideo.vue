<template lang="pug">
  .ui.embed(ref='video', data-source='youtube', :data-id='src')
</template>

<script>
export default {
  props: ['src'],
  mounted () {
    this.embed()
  },
  updated () {
    this.embed()
  },
  methods: {
    embed () {
      this.$nextTick(() => {
        $(this.$refs.video).embed()
      })
    }
  }
}
</script>
