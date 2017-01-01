<template lang="pug">
  .ui.segment
    h4.ui.header Course Content
    .ui.styled.fluid.accordion(ref='accordion')
      div(v-for='(x, i) in contents')
        .title
          i.dropdown.icon
          | Section {{ i + 1 }}
          span(v-if='x.title') : {{ x.title }}
        .content
          p.description(v-html='x.description')
          div(v-if='x.video')
            EmbedVideo(:src='x.video')
</template>

<script>
import EmbedVideo from './EmbedVideo'

export default {
  components: {
    EmbedVideo
  },
  props: ['contents'],
  mounted () {
    this.$nextTick(() => {
      $(this.$refs.accordion).accordion()
    })
  }
}
</script>
