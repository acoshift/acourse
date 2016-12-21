<template>
  <div class="ui card">
    <router-link class="image" :to="`/course/${course.id}`">
      <img :src="course.photo">
    </router-link>
    <div class="content">
      <router-link class="header" :to="`/course/${course.id}`">{{ course.title | trim(40) }}</router-link>
      <div class="meta">
        <span class="date">{{ course.start | date('DD/MM/YYYY') }}</span>
      </div>
      <div class="description">
        {{ course.shortDescription }}
      </div>
    </div>
  </div>
</template>

<style scoped>
  .card img {
    object-fit: cover;
    object-position: center center;
    height: 180px !important;
  }
</style>

<script>
export default {
  props: {
    course: {
      type: Object,
      required: true
    }
  }
}
</script>
