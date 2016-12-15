<template>
  <div class="ui basic segment">
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
      <div class="two column middle aligned row">
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
    </div>
  </div>
</template>

<style scoped>
  .segment {
    padding-top: 0;
  }
</style>

<script>
import Avatar from './Avatar'

export default {
  components: {
    Avatar
  },
  props: {
    course: {
      type: Object,
      required: true
    }
  }
}
</script>
