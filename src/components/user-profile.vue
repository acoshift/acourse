<template>
  <div>
    <div v-if="user" class="ui center aligned grid">
      <div class="row" style="padding-bottom: 0;">
        <avatar :src="user.photo" size="medium"></avatar>
      </div>
      <div class="row" style="padding-bottom: 0;">
        <h1>{{ user.name }}</h1>
      </div>
      <div class="row">
        <h3>{{ user.aboutMe }}</h3>
      </div>
    </div>
    <div v-if="!user">
      <div class="ui yellow message">No Profile Data</div>
    </div>
  </div>
</template>

<script>
  import Avatar from './avatar'

  export default {
    components: {
      Avatar
    },
    props: ['user']
  }
</script>
