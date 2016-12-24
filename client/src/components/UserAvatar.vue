<template>
  <div>
    <avatar :src="user.photo" size="tiny"></avatar>
    {{ user.name || 'Anonymous' }}
  </div>
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
