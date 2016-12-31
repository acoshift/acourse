<template lang="pug">
  .ui.small.modal
    .image.content
      .ui.centered.image
        i.huge.icons
          i.green.big.thin.circle.icon
          i.green.check.icon
    .description(style='text-align: center;')
      .ui.header {{ title }}
      p {{ description }}
      .ui.close.button(ref='closeButton') OK
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
      this.title = title
      this.description = description
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
