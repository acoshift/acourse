<template>
  <div>
    <navbar></navbar>
    <div id="view">
      <div class="ui container">
        <router-view></router-view>
      </div>
    </div>
  </div>
</template>

<style scoped>
  #view {
    padding-top: 80px;
  }
</style>

<script>
  import Navbar from './navbar'

  export default {
    components: {
      Navbar
    }
  }
</script>
