<template lang="pug">
  .ui.small.modal
    .image.content
      .ui.centered.image
        i.huge.icons
          i.red.big.thin.circle.icon
          i.red.cancel.icon
    .description(style='text-align: center;')
      .ui.header {{ title }}
      p {{ description }}
      .ui.close.blue.basic.button(ref='closeButton') OK
</template>

<style scoped>
  .modal {
    padding-bottom: 30px;
  }

  .modal .close.button {
    width: 40%;
    min-width: 150px;
    margin-top: 10px;
  }
</style>

<script>
export default {
  data () {
    return {
      title: '',
      description: ''
    }
  },
  methods: {
    show (title, description) {
      this.title = title || 'Error'
      this.description = description || 'Unknown Error'
      this.$nextTick(() => {
        $(this.$el)
          .modal('attach events', this.$refs.closeButton, 'hide')
          .modal('show')
      })
    },
    close () {
      $(this.$el).modal('hide')
    }
  }
}
</script>
