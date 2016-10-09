<template>
  <div class="ui basic segment" :class="{loading}" style="padding: 0; margin: 0; min-height: 100%;">
    <router-view></router-view>
    <success-modal ref="successModal"></success-modal>
  </div>
</template>

<script>
  import { Loader, Document } from './services'
  import SuccessModal from './components/success-modal'

  export default {
    components: {
      SuccessModal
    },
    data () {
      return {
        loader: Loader.state
      }
    },
    computed: {
      loading () {
        return !!this.loader.value
      }
    },
    created () {
      Document.$successModal
        .subscribe(
          ({title, description}) => {
            this.$refs.successModal.show(title, description)
          }
        )
    }
  }
</script>
