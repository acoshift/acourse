<template lang="pug">
.ui.segment
  h3.ui.header Edit Profile
  form.ui.form(@submit.prevent='submit')
    .ui.red.message(v-if='error') {{ error }}
    .field
      label Photo
      avatar(v-show='user.photo', :src='user.photo', size='small')
      .ui.green.button(@click='selectPhoto') Select Photo
    .field
      label Username
      input(v-model='user.username', maxlength='25')
    .field
      label Name
      input(v-model='user.name', maxlength='45')
    .field
      label About me
      input(v-model='user.aboutMe', maxlength='40')
    button.ui.blue.submit.button(:class='{loading: saving}') Save
    router-link.ui.red.button(to='/profile') Cancel
</template>

<style>
  img.circular.image {
    margin: 10px;
  }
</style>

<script>
import { Observable } from 'rxjs/Observable'
import { Loader, Me, Document } from 'services'
import Avatar from './Avatar'
import pick from 'lodash/fp/pick'
import keys from 'lodash/fp/keys'

export default {
  components: {
    Avatar
  },
  data () {
    return {
      user: {
        photo: '',
        name: '',
        username: '',
        aboutMe: ''
      },
      saving: false,
      error: ''
    }
  },
  created () {
    Observable.of({})
      .do(() => Loader.start('user'))
      .finally(() => Loader.stop('user'))
      .flatMap(() => Me.get())
      .first()
      .subscribe(
        (user) => {
          this.user = pick(keys(this.user))(user)
        }
      )
  },
  methods: {
    submit () {
      this.error = ''
      if (this.saving) return
      this.saving = true
      Me.update(this.user)
        .finally(() => { this.saving = false })
        .subscribe(
          () => {
            Me.fetch()
            this.$router.push('/profile')
          }
        )
    },
    selectPhoto () {
      Document.uploadModal.open('image/*')
        .subscribe(
          (f) => {
            this.user.photo = f.downloadURL
          },
          (err) => {
            Document.openErrorModal('Upload Error', err.message)
          }
        )
    }
  }
}
</script>
