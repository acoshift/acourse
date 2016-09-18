<template>
  <img class="ui circular image" :class="`a-${size}`" :src="src">
</template>

<style scoped>
  img {
    object-fit: cover;
    object-position: top;
  }

  img.a-small {
    width: 120px;
    height: 120px;
  }
</style>

<script>
  export default {
    props: ['src', 'size']
  }
</script>
