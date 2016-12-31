<template lang="pug">
  .ui.small.modal
    .image.content
      .ui.centered.image
        i.huge.cloud.upload.icon
    .description(style='text-align: center;')
      h2.ui.header Upload
      .ui.yellow.message Limit file size to 5 MiB
      .ui.indicating.progress(v-show='uploading', ref='progress')
        .bar
          .progress
      input.hidden(type='file', ref='file', :accept='accept', @change='upload')
      .ui.green.button(:class='{disabled: uploading}', @click='$refs.file.click()') Select File
      .ui.red.button(@click='cancel') Cancel
</template>

<style scoped>
  .modal {
    padding-bottom: 30px;
  }

  .description {
    margin-left: 2rem;
    margin-right: 2rem;
  }

  .ui.progress {
    margin-bottom: 1rem;
  }

  .button {
    width: 180px;
  }
</style>

<script>
import { Subject } from 'rxjs/Subject'
import { Me } from 'services'
import firebase from 'firebase'

export default {
  data () {
    return {
      uploading: false,
      o: null,
      task: null,
      accept: ''
    }
  },
  methods: {
    open (accept) {
      this.accept = accept || ''
      this.o = null
      this.task = null
      this.uploading = false
      this.$nextTick(() => {
        $(this.$el)
          .modal({
            closable: false
          })
          .modal('show')
        $(this.$refs.progress).progress()
      })
      this.o = new Subject()
      return this.o.asObservable()
    },
    close () {
      $(this.$el).modal('hide')
    },
    cancel () {
      if (this.task) this.task.cancel()
      this.close()
    },
    upload () {
      if (this.uploading) return
      const f = this.$refs.file.files[0]
      if (!f) return
      this.uploading = true
      this.$refs.file.value = ''
      Me.upload(f)
        .subscribe(
          (task) => {
            this.task = task
            task.then((snapshot) => {
              this.close()
              this.o.next(snapshot)
              this.o.complete()
            }, (err) => {
              this.close()
              this.o.error(err)
            })
            task.on(
              firebase.storage.TaskEvent.STATE_CHANGED,
              (snapshot) => {
                const percent = snapshot.bytesTransferred / snapshot.totalBytes * 100
                $(this.$refs.progress).progress('set progress', percent)
              }
            )
          }
        )
    }
  }
}
</script>
