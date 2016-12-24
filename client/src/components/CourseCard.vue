<template>
  <router-link class="ui card" :to="`/course/${url}`">
    <div class="image" :to="`/course/${url}`">
      <img :src="course.photo">
    </div>
    <div class="content">
      <div class="header" :to="`/course/${url}`">{{ course.title | trim(45) }}</div>
      <div class="meta">
        <span v-if="course.type === 'video'">Video</span>
        <span v-if="course.type === 'live'" class="date">Live start at {{ course.start | date('DD/MM/YYYY') }}</span>
      </div>
      <div class="description">
        {{ course.shortDescription }}
      </div>
    </div>
    <div class="extra content">
      <div class="right floated">
        <i class="user icon"></i> {{ course.student }}
      </div>
      <span v-if="course.price" class="price">฿ {{ course.price | money }}</span>
    </div>
  </router-link>
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
</style>

<script>
export default {
  props: {
    course: {
      type: Object,
      required: true
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
