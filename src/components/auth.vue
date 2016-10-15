<template>
  <div id="auth" class="ui fluid container">
    <div class="ui center aligned grid">
      <div class="row">
        <div class="column">
          <img src="../assets/acourse.svg" class="ui centered image" style="width: 220px; height: 220px">
        </div>
      </div>
      <div class="row">
        <h1>Acourse</h1>
      </div>
      <div class="row">
        <div class="ui center aligned stacked segment">
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
                <div class="ui right aligned basic fluid container" style="padding-bottom: 1rem;">
                  <a href="#" @click="state = 1">Forget password ?</a>
                  /
                  <a href="#" @click="state = 2">Register</a>
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
                <div class="ui fluid basic segment">
                  <a href="#" @click="state = 0"><i class="left chevron icon"></i> Sign In</a>
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
                <div class="ui fluid basic segment">
                  <a href="#" @click="state = 0"><i class="left chevron icon"></i> Sign In</a>
                </div>
                <button class="ui blue fluid submit button" :class="{loading}">Sign Up</button>
              </form>
            </div>
          </div>
        </div>
      </div>
      <div class="row">
        <div class="ui basic segment" style="width: 60%; margin: 0; padding: 0">
          <div class="ui horizontal divider">Or</div>
        </div>
      </div>
      <div class="row">
        <div class="ui center aligned segment">
          <div class="ui facebook fluid button" @click="facebookSignIn" :class="{loading: facebookLoading}"><i class="facebook f icon"></i>Sign In with Facebook</div>
          <br>
          <div class="ui google plus fluid button" :class="{loading: googleLoading}" @click="googleSignIn"><i class="google plus icon"></i>Sign In with Google+</div>
          <br>
          <div class="ui black fluid button" :class="{loading: githubLoading}" @click="githubSignIn"><i class="github icon"></i>Sign In with Github</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  #auth {
    padding-top: 10px;
  }

  .segment {
    min-width: 320px;
  }

  @media only screen and (min-width: 500px) {
    .segment {
      min-width: 500px;
    }
  }

  @media only screen and (min-height: 850px) {
    #auth {
      padding-top: 5vh;
    }
  }

  .basic.segment {
    padding-top: 0;
    padding-bottom: 0;
    padding-left: 0;
    padding-right: 35px;
  }
</style>

<script>
  import { Auth, User, Document } from '../services'

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
      signIn () {
        if (this.loading) return
        this.loading = true
        this.error = ''
        Auth.signIn(this.email, this.password)
          .finally(() => { this.loading = false })
          .subscribe(
            () => {
              this.gotoHome()
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
              this.gotoHome()
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
          .flatMap((res) => User.saveAuthProfile(res.user), (x) => x)
          .finally(() => { this.facebookLoading = false })
          .subscribe(
            () => {
              this.gotoHome()
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
          .flatMap((res) => User.saveAuthProfile(res.user), (x) => x)
          .finally(() => { this.googleLoading = false })
          .subscribe(
            () => {
              this.gotoHome()
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
          .flatMap((res) => User.saveAuthProfile(res.user), (x) => x)
          .finally(() => { this.githubLoading = false })
          .subscribe(
            () => {
              this.gotoHome()
            },
            (err) => {
              Document.openErrorModal('Error', err.message)
            }
          )
      },
      resetError () {
        this.error = ''
      },
      gotoHome () {
        this.$nextTick(() => {
          this.$router.push('/home')
        })
      }
    },
    watch: {
      state () {
        this.resetError()
      }
    }
  }
</script>
