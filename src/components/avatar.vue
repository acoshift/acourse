<template>
  <img class="ui circular image" :class="`a-${size}`" :src="src || '/static/icons/ic_face_black_48px.svg'">
</template>

<style scoped>
  img {
    object-fit: cover;
    object-position: top;
    display: inline-block !important;
    margin-top: 0 !important;
    margin-bottom: 0 !important;
    margin-left: 0 !important;
  }

  img.a-tiny {
    width: 32px;
    height: 32px;
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
