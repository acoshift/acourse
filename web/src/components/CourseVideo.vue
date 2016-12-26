<template>
  <div class="ui segment">
    <div ref="video" class="ui embed" data-source="youtube" :data-id="src"></div>
  </div>
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
