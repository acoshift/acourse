<template lang="pug">
  #auth.ui.fullscreen.modal
    i.close.icon
    .header Sign In to Acourse.io
    .content
      .ui.stackable.grid
        .row
          .nine.wide.column
            div(v-show='state === 0')
              .ui.fluid.left.aligned.container
                form.ui.form(:class='{error}', @submit.prevent='signIn')
                  .field
                    label Email
                    input(v-model='email', @focus='resetError')
                  .field
                    label Password
                    input(v-model='password', type='password', @focus='resetError')
                  .ui.error.message
                    | {{ error }}
                  .ui.right.aligned.fluid.container(style='padding-bottom: 1rem;')
                    a(href='#', target='_self', @click.prevent='state = 1') Forget password ?
                    | &nbsp;/&nbsp;
                    a(href='#', target='_self', @click.prevent='state = 2') Register
                  button.ui.blue.fluid.submit.button(:class='{loading}') Sign In
            div(v-show='state === 1')
              h3 Reset Password
              .ui.fluid.left.aligned.container
                form.ui.form(:class='{error}', @submit.prevent='forgot')
                  .field
                    label Email
                    input(v-model='email', @focus='resetError')
                  .ui.error.message
                    | {{ error }}
                  .ui.fluid.container(style='padding-bottom: 1rem;')
                    a(href='#', target='_self', @click.prevent='state = 0')
                      i.left.chevron.icon
                      | &nbsp;Sign In
                  button.ui.blue.fluid.submit.button(:class='{loading}') Reset
            div(v-show='state === 2')
              h3 Register
              .ui.fluid.left.aligned.container
                form.ui.form(:class='{error}', @submit.prevent='signUp')
                  .field
                    label Email
                    input(v-model='email', @focus='resetError')
                  .field
                    label Password
                    input(v-model='password', type='password', @focus='resetError')
                  .ui.error.message
                    | {{ error }}
                  .ui.fluid.container(style='padding-bottom: 1rem;')
                    a(href='#', target='_self', @click.prevent='state = 0')
                      i.left.chevron.icon
                      | Sign In
                  button.ui.blue.fluid.submit.button(:class='{loading}') Sign Up
          .one.wide.column
            .ui.vertical.divider Or
          .six.wide.middle.aligned.column
            .ui.facebook.fluid.button(@click='facebookSignIn', :class='{loading: facebookLoading}')
              i.facebook.f.icon
              | Sign In with Facebook
            br
            .ui.google.plus.fluid.button(:class='{loading: googleLoading}', @click='googleSignIn')
              i.google.plus.icon
              | Sign In with Google+
            br
            .ui.black.fluid.button(:class='{loading: githubLoading}', @click='githubSignIn')
              i.github.icon
              | Sign In with Github
</template>

<script>
import { Auth, Document } from 'services'

export default {
  data () {
    return {
      email: '',
      password: '',
      state: 0,
      error: '',
      loading: false,
      facebookLoading: false,
      googleLoading: false,
      githubLoading: false
    }
  },
  methods: {
    open () {
      this.state = 0
      this.resetError()
      $(this.$el).modal('show')
    },
    close () {
      this.email = ''
      this.password = ''
      $(this.$el).modal('hide')
    },
    signIn () {
      if (this.loading) return
      this.loading = true
      this.error = ''
      Auth.signIn(this.email, this.password)
        .finally(() => { this.loading = false })
        .subscribe(
          () => {
            this.close()
          },
          () => {
            this.error = 'Email or password wrong'
          }
        )
    },
    forgot () {
      if (this.loading) return
      this.loading = true
      this.error = ''
      Auth.resetPassword(this.email)
        .finally(() => { this.loading = false })
        .subscribe(
          () => {
            this.email = ''
            Document.openSuccessModal('Success', 'Please check email to reset password.')
          },
          (err) => {
            this.error = err.message
          }
        )
    },
    signUp () {
      if (this.loading) return
      this.loading = true
      this.error = ''
      Auth.signUp(this.email, this.password)
        .finally(() => { this.loading = false })
        .subscribe(
          (res) => {
            this.close()
          },
          (err) => {
            this.error = err.message
          }
        )
    },
    facebookSignIn () {
      if (this.facebookLoading) return
      this.facebookLoading = true
      Auth.signInWithFacebook()
        .finally(() => { this.facebookLoading = false })
        .subscribe(
          () => {
            this.close()
          },
          (err) => {
            Document.openErrorModal('Error', err.message)
          }
        )
    },
    googleSignIn () {
      if (this.googleLoading) return
      this.googleLoading = true
      Auth.signInWithGoogle()
        .finally(() => { this.googleLoading = false })
        .subscribe(
          () => {
            this.close()
          },
          (err) => {
            Document.openErrorModal('Error', err.message)
          }
        )
    },
    githubSignIn () {
      if (this.githubLoading) return
      this.githubLoading = true
      Auth.signInWithGithub()
        .finally(() => { this.githubLoading = false })
        .subscribe(
          () => {
            this.close()
          },
          (err) => {
            Document.openErrorModal('Error', err.message)
          }
        )
    },
    resetError () {
      this.error = ''
    }
  },
  watch: {
    state () {
      this.resetError()
    }
  }
}
</script>
