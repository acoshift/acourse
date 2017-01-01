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
          div.video(v-if='x.video')
            .ui.embed(data-source='youtube', :data-id='x.video')
</template>

<script>
export default {
  props: ['contents'],
  mounted () {
    this.$nextTick(() => {
      $(this.$refs.accordion).accordion({
        onOpening () {
          const video = $(this).find('.video > .ui.embed')
          if (!video.hasClass('active')) {
            video.embed()
          }
        }
      })
    })
  }
}
</script>
