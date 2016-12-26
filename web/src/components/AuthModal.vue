<template>
  <div id="auth" class="ui fullscreen modal">
    <i class="close icon"></i>
    <div class="header">Sign In to Acourse.io</div>
    <div class="content">
      <div class="ui stackable grid">
        <div class="row">
          <div class="nine wide column">
            <div v-show="state === 0">
              <div class="ui fluid left aligned container">
                <form class="ui form" :class="{error}" @submit.prevent="signIn">
                  <div class="field">
                    <label>Email</label>
                    <input v-model="email" @focus="resetError">
                  </div>
                  <div class="field">
                    <label>Password</label>
                    <input v-model="password" type="password" @focus="resetError">
                  </div>
                  <div class="ui error message">
                    {{ error }}
                  </div>
                  <div class="ui right aligned fluid container" style="padding-bottom: 1rem;">
                    <a href="#" target="_self" @click.prevent="state = 1">Forget password ?</a>
                    /
                    <a href="#" target="_self" @click.prevent="state = 2">Register</a>
                  </div>
                  <button class="ui blue fluid submit button" :class="{loading}">Sign In</button>
                </form>
              </div>
            </div>
            <div v-show="state === 1">
              <h3>Reset Password</h3>
              <div class="ui fluid left aligned container">
                <form class="ui form" :class="{error}" @submit.prevent="forgot">
                  <div class="field">
                    <label>Email</label>
                    <input v-model="email" @focus="resetError">
                  </div>
                  <div class="ui error message">
                    {{ error }}
                  </div>
                  <div class="ui fluid container" style="padding-bottom: 1rem;">
                    <a href="#" target="_self" @click.prevent="state = 0"><i class="left chevron icon"></i> Sign In</a>
                  </div>
                  <button class="ui blue fluid submit button" :class="{loading}">Reset</button>
                </form>
              </div>
            </div>
            <div v-show="state === 2">
              <h3>Register</h3>
              <div class="ui fluid left aligned container">
                <form class="ui form" :class="{error}" @submit.prevent="signUp">
                  <div class="field">
                    <label>Email</label>
                    <input v-model="email" @focus="resetError">
                  </div>
                  <div class="field">
                    <label>Password</label>
                    <input v-model="password" type="password" @focus="resetError">
                  </div>
                  <div class="ui error message">
                    {{ error }}
                  </div>
                  <div class="ui fluid container" style="padding-bottom: 1rem;">
                    <a href="#" target="_self" @click.prevent="state = 0"><i class="left chevron icon"></i> Sign In</a>
                  </div>
                  <button class="ui blue fluid submit button" :class="{loading}">Sign Up</button>
                </form>
              </div>
            </div>
          </div>
          <div class="one wide column">
            <div class="ui vertical divider">Or</div>
          </div>
          <div class="six wide middle aligned column">
            <div class="ui facebook fluid button" @click="facebookSignIn" :class="{loading: facebookLoading}"><i class="facebook f icon"></i> Sign In with Facebook</div>
            <br>
            <div class="ui google plus fluid button" :class="{loading: googleLoading}" @click="googleSignIn"><i class="google plus icon"></i> Sign In with Google+</div>
            <br>
            <div class="ui black fluid button" :class="{loading: githubLoading}" @click="githubSignIn"><i class="github icon"></i> Sign In with Github</div>
          </div>
        </div>
      </div>
    </div>
  </div>
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
