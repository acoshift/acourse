<template>
  <div class="ui segment">
    <div class="ui center aligned grid">
      <div class="row">
        <div class="column">
          <img :src="course.photo" class="ui centered big image">
        </div>
      </div>
      <div class="row" style="padding-top: 0;">
        <div class="column">
          <h1>{{ course.title }}</h1>
        </div>
      </div>
      <div class="row" style="margin-top: -2rem; margin-bottom: 1rem;">
        <div class="column">
          <i>{{ course.start | date('DD/MM/YYYY') }}</i>
        </div>
      </div>
      <div class="two column middle aligned row" style="margin-top: -30px !important;">
        <div class="right aligned column" style="padding-right: 2px;">
          <router-link :to="`/user/${course.owner.id}`">
            <avatar :src="course.owner.photo" size="mini"></avatar>
          </router-link>
        </div>
        <div class="left aligned column" style="padding-left: 2px;">
          <router-link :to="`/user/${course.owner.id}`">
            <h3>{{ course.owner.name || 'Anonymous' }}</h3>
          </router-link>
        </div>
      </div>
      <div class="row">
        <div class="column">
          <p class="description">{{ course.description }}</p>
        </div>
      </div>
    </div>
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
  import Avatar from './avatar'

  export default {
    props: ['course'],
    components: {
      Avatar
    }
  }
</script>
