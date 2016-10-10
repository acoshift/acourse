<template>
  <div class="ui segment">
    <h4 class="ui header">Course Content</h4>
    <div class="ui styled fluid accordion" ref="accordion">
      <div v-for="(x, i) in contents">
        <div class="title"><i class="dropdown icon"></i>Section {{ i + 1 }}</div>
        <div class="content">
          <p class="description" v-html="x.content"></p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    props: ['contents'],
    mounted () {
      this.$nextTick(() => {
        $(this.$refs.accordion).accordion()
      })
    }
  }
</script>
