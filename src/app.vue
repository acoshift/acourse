<template>
  <div class="ui basic segment" :class="{loading}" style="padding: 0; margin: 0; min-height: 100%;">
    <router-view></router-view>
  </div>
</template>

<script>
  import { Loader } from './services'

  export default {
    data () {
      return {
        loader: Loader.state
      }
    },
    computed: {
      loading () {
        return !!this.loader.value
      }
    }
  }
</script>
