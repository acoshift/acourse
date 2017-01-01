<template lang="pug">
  .ui.center.aligned.grid
    .row(style='padding-bottom: 0;')
      Avatar(:src='user && user.photo', size='medium')
    .row(v-if='user && user.username', style='padding-bottom: 0;')
      h2 @{{ user.username }}
    .row(style='padding-bottom: 0;')
      h1 {{ user && user.name || 'Anonymous' }}
    .row
      h3 {{ user && user.aboutMe }}
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
