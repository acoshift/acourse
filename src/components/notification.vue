<template>
  <div class="noti">
    <div class="ui floating message" v-for="x in data" :class="x.color || ''" @click="close(x.id)">
      <i class="close icon"></i>
      {{ x.body }}
    </div>
  </div>
</template>

<style scoped>
  .noti {
    margin: 0;
    position: absolute;
    z-index: 150;
    top: 5rem;
    right: 1rem;
    width: 19rem;
    cursor: pointer;
  }

  .ui.message {
    margin: 0 0 0.6rem 0;
  }
</style>

<script>
import { Firebase } from '../services'
import uniqueId from 'lodash/uniqueId'
import findIndex from 'lodash/fp/findIndex'

export default {
  data () {
    return {
      data: []
    }
  },
  created () {
    Firebase.notification
      .subscribe(
        (payload) => {
          const id = uniqueId('notification_')
          if (payload.data) {
            this.data.push({
              id,
              body: payload.data.body || payload.notification.body,
              color: payload.data.color || 'info'
            })
            if (payload.data.timeout !== -1) {
              setTimeout(() => {
                this.close(id)
              }, payload.data.timeout || 5000)
            }
          } else {
            this.data.push({
              id,
              body: payload.notification.body,
              color: 'info'
            })
            setTimeout(() => {
              this.close(id)
            }, 5000)
          }
          if (this.data.length > 5) {
            this.data.splice(0, 1)
          }
        }
      )
  },
  methods: {
    close (id) {
      const i = findIndex({ id })(this.data)
      if (i >= 0) this.data.splice(i, 1)
    }
  }
}
</script>
