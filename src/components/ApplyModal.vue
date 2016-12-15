<template lang="pug">
  .ui.small.modal
    .header Code
    .content
      form.ui.form(@submit.prevent="submit")
        .field
          label Enter Code
          input(v-model="code")
        button.ui.fluid.submit.blue.button Apply
</template>

<script>
export default {
  data () {
    return {
      code: ''
    }
  },
  methods: {
    show () {
      this.code = ''
      $(this.$el).modal('show')
    },
    submit () {
      $(this.$el).modal('hide')
      this.$emit('apply', this.code)
    }
  }
}
</script>
