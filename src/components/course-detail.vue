<template>
  <div class="ui segment">
    <h4 class="ui header">Course Detail</h4>
    <span>Start: {{ course.start | date('DD/MM/YYYY') }}</span>
    <div class="ui divider"></div>
    <p class="description">{{ course.description }}</p>
  </div>
</template>

<style scoped>
  p.description {
    text-align: left;
    white-space: pre-wrap;
    word-break: break-word;
  }
</style>

<script>
  export default {
    props: ['course']
  }
</script>
