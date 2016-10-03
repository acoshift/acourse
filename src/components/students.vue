<template>
  <div class="ui styled fluid accordion">
    <div class="title">
      <i class="dropdown icon"></i>
      Students <span v-if="users">({{ users.length }})</span>
    </div>
    <div class="ui content">
      <div class="ui stackable three column grid">
        <div class="column" v-for="x in users">
          <router-link :to="`/user/${x.id}`">
            <avatar :src="x.photo" size="tiny"></avatar>
            {{ x.name || 'Anonymous' }}
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .accordion {
    .grid {
      .column {
        padding: .6rem;
      }
    }
  }
</style>

<script>
  import Avatar from './avatar'

  export default {
    props: ['users'],
    components: {
      Avatar
    },
    mounted () {
      $(this.$el).accordion()
    }
  }
</script>
