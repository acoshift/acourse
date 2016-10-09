<template>
  <div class="ui basic segment" :class="{loading}" style="padding: 0; margin: 0; min-height: 100%;">
    <router-view></router-view>
    <success-modal ref="successModal"></success-modal>
    <error-modal ref="errorModal"></error-modal>
  </div>
</template>

<script>
  import { Loader, Document } from './services'
  import SuccessModal from './components/success-modal'
  import ErrorModal from './components/error-modal'

  export default {
    components: {
      SuccessModal,
      ErrorModal
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
      Document.$errorModal
        .subscribe(
          ({title, description}) => {
            this.$refs.errorModal.show(title, description)
          }
        )
    }
  }
</script>
