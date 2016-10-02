<template>
  <div class="ui basic segment" :class="{loading: loading > 0}">
    <div class="ui massive breadcrumb">
      <router-link class="section" to="/home">Courses</router-link>
      <i class="right chevron icon divider"></i>
      <router-link class="section" :to="`/course/${courseId}`">{{ course && course.title || courseId }}</router-link>
      <i class="right chevron icon divider"></i>
      <div class="active section">Chat</div>
    </div>
    <div class="ui chat segment" ref="container">
      <div class="ui segment" id="chatBox" ref="chatBox">
        <div class="ui comments">
          <div v-for="x in messages" class="comment">
            <span class="avatar">
              <img :src="x.user && x.user.photo || '/static/icons/ic_face_black_48px.svg'">
            </span>
            <div class="content">
              <span class="author">{{ x.user && x.user.name || 'Anonymous' }}</span>
              <div class="metadata">
                <span class="date">{{ x.t | fromNow }}</span>
              </div>
              <div class="text">
                <a v-if="x.h" target="_blank" :href="x.m">
                  <div v-if="x.h === 1" class="ui small image">
                    <img :src="x.m" onerror="this.src = '/static/icons/ic_insert_drive_file_black_48px.svg'">
                  </div>
                  <span v-else>{{ x.m }}</span>
                </a>
                <span v-else>{{ x.m }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="ui segment">
        <div class="ui grid">
          <div class="row">
            <div class="column">
              <div class="ui fluid input">
                <input ref="input" v-model="input" @keyup.13="send"></input>
                <div class="ui basic icon button" @click="$refs.file.click()" :class="{'disabled loading': uploading}"><i class="upload icon"></i></div>
                <div class="ui basic icon button" @click="send"><i class="send icon"></i></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <input type="file" class="hidden" ref="file" @change="upload">
  </div>
</template>

<style scoped>
  .chat.segment {
    padding: 0;
    overflow: hidden;
  }

  .chat.segment > .segment {
    padding: 0;
    margin: 0;
  }

  #chatBox {
    padding: 10px;
    overflow-y: scroll;
    overflow-x: hidden;
  }

  .chat.segment .button {
    margin: 0;
  }

  .avatar img {
    border-radius: 50% !important;
    height: 35px !important;
    object-fit: cover;
    object-position: top;
  }
</style>

<script>
  import { Course, User } from '../services'
  import Avatar from './avatar'
  import startsWith from 'lodash/fp/startsWith'

  export default {
    components: {
      Avatar
    },
    data () {
      return {
        loading: false,
        course: null,
        courseId: null,
        input: '',
        messages: [],
        ob: [],
        limit: 50,
        loadingTop: false,
        uploading: false
      }
    },
    created () {
      this.courseId = this.$route.params.id
      this.loading = 2
      this.ob.push(Course.get(this.courseId)
        .subscribe(
          (course) => {
            --this.loading
            this.course = course
          },
          () => {
            this.loading = 0
            this.$router.replace('/home')
          }
        )
      )
    },
    destroyed () {
      this.ob.forEach((x) => x.unsubscribe())
    },
    mounted () {
      this.adjust()
      $(window).resize(this.adjust)
      $(this.$refs.chatBox)
        .off()
        .on('scroll', () => {
          let pos = $(this.$refs.chatBox).scrollTop()
          if (pos <= 5) {
            if (this.limit > this.messages.length) return
            this.limit += 30
            this.loading = 1
            this.initMessages()
            this.$nextTick(() => {
              $(this.$refs.chatBox).scrollTop(700)
            })
          }
        })

      this.initMessages()
    },
    methods: {
      isUrl (text) {
        if (/^https?:\/\/?[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]$/.test(text)) {
          if (startsWith('https://firebasestorage.googleapis.com/v0/b/acourse-d9d0a.appspot.com')(text)) {
            return 1
          }
          return 2
        }
        return 0
      },
      initMessages () {
        const messagesObservable = Course.messages(this.courseId, this.limit)
          .do((message) => {
            User.getOnce(message.u).subscribe((user) => {
              message.user = user
            })
          })

        const messages = []

        this.ob.push(messagesObservable
          .subscribe(
            (message) => {
              message.h = this.isUrl(message.m)
              messages.push(message)
            }
          )
        )

        let shouldScroll = false

        this.ob.push(messagesObservable
          .do(() => {
            if (this.$refs.chatBox.scrollHeight - this.$refs.chatBox.scrollTop <= this.$refs.chatBox.clientHeight + 500) {
              shouldScroll = true
            }
          })
          .debounceTime(200)
          .subscribe(
            () => {
              if (this.loading > 0) --this.loading
              this.messages = messages
              if (shouldScroll) {
                shouldScroll = false
                this.$nextTick(() => {
                  $(this.$refs.chatBox).scrollTop(99999)
                })
              }
            }
          )
        )
      },
      adjust () {
        let container = $(this.$refs.container)
        let box = $(this.$refs.chatBox)
        let h = window.innerHeight
        let input = $(this.$refs.input)
        container.height(() => h - container.offset().top - 20)
        box.height(() => h - box.offset().top - 80)
        box.scrollTop(99999)
        input.focus()
      },
      send () {
        const input = this.input
        if (!input || !input.trim()) return
        this.input = ''
        Course.sendMessage(this.courseId, input).subscribe()
      },
      upload () {
        if (this.uploading) return
        const f = this.$refs.file.files[0]
        if (!f) return
        this.uploading = true
        this.$refs.file.value = ''
        User.upload(f)
          .flatMap((file) => Course.sendMessage(this.courseId, file.downloadURL))
          .finally(() => { this.uploading = false })
          .subscribe(
            null,
            () => {
              window.alert('Please check file size.')
            }
          )
      }
    }
  }
</script>
