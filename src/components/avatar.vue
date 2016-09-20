<template>
  <img class="ui circular image" :class="`a-${size}`" :src="src">
</template>

<style scoped>
  img {
    object-fit: cover;
    object-position: top;
    display: inline-block !important;
  }

  img.a-mini {
    width: 64px;
    height: 64px;
  }

  img.a-small {
    width: 120px;
    height: 120px;
  }

  img.a-medium {
    width: 200px;
    height: 200px;
  }
</style>

<script>
  export default {
    props: ['src', 'size']
  }
</script>
