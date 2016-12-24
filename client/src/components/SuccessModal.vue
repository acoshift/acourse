<template>
  <div class="ui small modal">
    <div class="image content">
      <div class="ui centered image">
        <i class="huge icons">
          <i class="green big thin circle icon"></i>
          <i class="green check icon"></i>
        </i>
      </div>
    </div>
    <div class="description" style="text-align: center;">
      <div class="ui header">{{ title }}</div>
      <p>{{ description }}</p>
      <div ref="closeButton" class="ui close button">OK</div>
    </div>
  </div>
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
