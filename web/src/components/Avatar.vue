<template lang="pug">
  img.ui.circular.image(:class="`a-${size}`", :src="src || '/static/icons/ic_face_black_48px.svg'")
</template>

<style lang="scss" scoped>
  img {
    object-fit: cover;
    object-position: top;
    display: inline-block !important;
    margin-top: 0 !important;
    margin-bottom: 0 !important;
    margin-left: 0 !important;

    &.a-tiny {
      width: 32px;
      height: 32px;
    }

    &.a-mini {
      width: 64px;
      height: 64px;
    }

    &.a-small {
      width: 120px;
      height: 120px;
    }

    &.a-medium {
      width: 200px;
      height: 200px;
    }
  }
</style>

<script>
export default {
  props: {
    src: {
      type: String
    },
    size: {
      type: String,
      required: true
    }
  }
}
</script>
