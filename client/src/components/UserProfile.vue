<template>
  <div>
    <div class="ui center aligned grid">
      <div class="row" style="padding-bottom: 0;">
        <avatar :src="user && user.photo" size="medium"></avatar>
      </div>
      <div v-if="user && user.username" class="row" style="padding-bottom: 0;">
        <h2>@{{ user.username }}</h1>
      </div>
      <div class="row" style="padding-bottom: 0;">
        <h1>{{ user && user.name || 'Anonymous' }}</h1>
      </div>
      <div class="row">
        <h3>{{ user && user.aboutMe }}</h3>
      </div>
    </div>
  </div>
</template>

<script>
import Avatar from './Avatar'

export default {
  components: {
    Avatar
  },
  props: ['user']
}
</script>
