<template lang="pug">
  div
    Avatar(:src='user.photo', size='tiny')
    | &nbsp;{{ user.name || 'Anonymous' }}
</template>

<script>
import Avatar from './Avatar'

export default {
  props: ['user'],
  components: {
    Avatar
  }
}
</script>
