<template lang="pug">
  div
    Navbar
    #view
      .ui.container
        router-view
</template>

<style scoped>
  #view {
    padding-top: 70px;
    padding-bottom: 30px;
  }
</style>

<script>
import Navbar from './Navbar'

export default {
  components: {
    Navbar
  }
}
</script>
