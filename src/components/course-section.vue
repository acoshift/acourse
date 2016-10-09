<template>
  <div>
    <div class="title">Section 1</div>
    <div class="content">{{ content.content }}</div>
  </div>
</template>

<script>
  export default {
    props: ['content']
  }
</script>
