<template>
  <div class="ui basic segment" :class="{loading}" style="padding: 0; margin: 0; min-height: 100%;">
    <router-view></router-view>
    <success-modal ref="successModal"></success-modal>
    <error-modal ref="errorModal"></error-modal>
    <upload-modal ref="uploadModal"></upload-modal>
    <notification></notification>
  </div>
</template>

<style>
  @import url('https://fonts.googleapis.com/css?family=Lato:400,700,400italic,700italic&subset=latin');

  @import '../node_modules/semantic-ui-css/components/reset.min.css';
  @import '../node_modules/semantic-ui-css/components/site.min.css';
  @import '../node_modules/semantic-ui-css/components/accordion.min.css';
  @import '../node_modules/semantic-ui-css/components/breadcrumb.min.css';
  @import '../node_modules/semantic-ui-css/components/button.min.css';
  @import '../node_modules/semantic-ui-css/components/card.min.css';
  @import '../node_modules/semantic-ui-css/components/checkbox.min.css';
  @import '../node_modules/semantic-ui-css/components/comment.min.css';
  @import '../node_modules/semantic-ui-css/components/container.min.css';
  @import '../node_modules/semantic-ui-css/components/dimmer.min.css';
  @import '../node_modules/semantic-ui-css/components/divider.min.css';
  @import '../node_modules/semantic-ui-css/components/dropdown.min.css';
  @import '../node_modules/semantic-ui-css/components/embed.min.css';
  @import '../node_modules/semantic-ui-css/components/form.min.css';
  @import '../node_modules/semantic-ui-css/components/grid.min.css';
  @import '../node_modules/semantic-ui-css/components/header.min.css';
  @import '../node_modules/semantic-ui-css/components/icon.min.css';
  @import '../node_modules/semantic-ui-css/components/image.min.css';
  @import '../node_modules/semantic-ui-css/components/input.min.css';
  @import '../node_modules/semantic-ui-css/components/item.min.css';
  @import '../node_modules/semantic-ui-css/components/list.min.css';
  @import '../node_modules/semantic-ui-css/components/menu.min.css';
  @import '../node_modules/semantic-ui-css/components/message.min.css';
  @import '../node_modules/semantic-ui-css/components/modal.min.css';
  @import '../node_modules/semantic-ui-css/components/segment.min.css';
  @import '../node_modules/semantic-ui-css/components/transition.min.css';
  @import '../node_modules/semantic-ui-css/components/progress.min.css';

  .hidden {
    display: none;
  }

  .text-center {
    text-align: center;
  }

  p.description {
    text-align: left;
    white-space: pre-wrap;
    word-break: break-word;
  }
</style>

<script>
  import { Loader, Document } from './services'
  import SuccessModal from './components/success-modal'
  import ErrorModal from './components/error-modal'
  import UploadModal from './components/upload-modal'
  import Notification from './components/notification'

  export default {
    components: {
      SuccessModal,
      ErrorModal,
      UploadModal,
      Notification
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
    },
    mounted () {
      Document.uploadModal = this.$refs.uploadModal
    }
  }
</script>
