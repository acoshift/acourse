<template>
  <div class="ui card">
    <router-link class="image" :to="`/course/${course.id}`">
      <img :src="course.photo">
    </router-link>
    <div class="content">
      <router-link class="header" :to="`/course/${course.id}`">{{ course.title | trim(30) }}</router-link>
      <div class="meta">
        <span class="date">{{ course.timestamp | date('DD/MM/YYYY') }}</span>
      </div>
      <div class="description">
        {{ course.description | trim(40) }}
      </div>
    </div>
    <div class="extra content">
      <div>
        <span class="right floated">
          <i class="user icon"></i>
          0
        </span>
        <span>
          <i class="heart outline link icon"></i>
          0
        </span>
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
    props: ['course']
  }
</script>
