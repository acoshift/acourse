<template lang="pug">
  router-link.ui.card(:to='`/course/${url}`')
    .image(:to='`/course/${url}`')
      img(:src='course.photo')
    .content
      .header(:to='`/course/${url}`') {{ course.title | trim(45) }}
      .meta
        span(v-if="course.type === 'video'") Video
        span.date(v-if="course.type === 'live'") Live start at {{ course.start | date('DD/MM/YYYY') }}
      .description
        | {{ course.shortDescription }}
    .extra.content
      .right.floated
        i.user.icon
        | &nbsp;{{ course.student }}
      span.price(v-if='!hidePrice && !course.discount && course.price', :class='{line: course.discount}') ฿ {{ course.price | money }}
      span.discount.price(v-if='!hidePrice && course.discount') &nbsp;฿ {{ course.discountedPrice | money }}
</template>

<style scoped>
  .card img {
    object-fit: cover;
    object-position: center center;
    height: 180px !important;
  }

  .card > .content > .header {
    font-size: 1.2em !important;
  }

  .price {
    font-size: 1.3em;
  }

  .discount.price {
    color: red;
  }

  .price.line {
    text-decoration: line-through;
    font-size: 1.0em !important;
  }
</style>

<script>
export default {
  props: {
    course: {
      type: Object,
      required: true
    },
    hidePrice: {
      type: Boolean,
      required: false
    }
  },
  computed: {
    url () {
      if (!this.course) return ''
      return this.course.url || this.course.id
    }
  }
}
</script>
