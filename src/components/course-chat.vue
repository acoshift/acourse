<template>
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
                  <img :src="x.m" onerror="this.src = '/static/icons/ic_insert_drive_file_black_48px.svg'" @load="scroll($event.target.height)">
                </div>
                <span v-else>{{ x.m }}</span>
              </a>
              <span v-else>{{ x.m }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div id="input" class="ui segment">
      <div class="ui grid">
        <div class="row">
          <div class="column">
            <div class="ui fluid input">
              <input ref="input" v-model="input" @keyup.13="send"></input>
              <div class="ui basic icon button" @click="upload"><i class="upload icon"></i></div>
              <div class="ui basic icon button" @click="send"><i class="send icon"></i></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
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
    height: calc(100vh - 17rem);
  }

  #input.segment {
    & input {
      border: none !important;
    }
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
  import { Course, Me, User, Loader, Document } from '../services'
  import Avatar from './avatar'
  import startsWith from 'lodash/fp/startsWith'

  export default {
    components: {
      Avatar
    },
    data () {
      return {
        course: null,
        courseId: null,
        input: '',
        messages: [],
        limit: 50,
        loadingTop: false,
        $course: null,
        $message: null,
        unread: 0,
        hidden: false,
        savedHeight: 0,
        paging: false
      }
    },
    created () {
      Loader.start('course')
      Loader.start('message')

      this.courseId = this.$route.params.id
      this.$course = Course.get(this.courseId)
        .subscribe(
          (course) => {
            Loader.stop('course')
            this.course = course
          },
          () => {
            this.$router.replace('/home')
          }
        )

      this.$visibilityChanged = Document.visibilityChanged()
        .subscribe(
          (hidden) => {
            this.hidden = hidden
            if (!hidden) {
              this.unread = 0
              Document.setTitle()
            }
          }
        )

      Course.lastMessage(this.courseId)
        .first()
        .subscribe(
          (message) => {
            if (!message) Loader.stop('message')
          }
        )
    },
    destroyed () {
      if (this.$message) this.$message.unsubscribe()
      this.$course.unsubscribe()
      this.$visibilityChanged.unsubscribe()
    },
    mounted () {
      $(this.$refs.chatBox)
        .off()
        .on('scroll', () => {
          const pos = this.$refs.chatBox.scrollTop
          if (pos <= 5 && !this.paging) {
            if (this.limit > this.messages.length) return
            this.limit += 30
            Loader.start('message')
            this.savedHeight = this.$refs.chatBox.scrollHeight
            this.paging = true
            this.initMessages()
          }
        })

      this.initMessages()
    },
    methods: {
      isUrl (text) {
        return /^https?:\/\/?[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]$/.test(text)
          ? (startsWith('https://firebasestorage.googleapis.com/v0/b/acourse-d9d0a.appspot.com')(text) ? 1 : 2)
          : 0
      },
      initMessages () {
        const messages = []
        let shouldScroll = false

        if (this.$message) this.$message.unsubscribe()

        this.$message = Course.messages(this.courseId, this.limit)
          .do((message) => {
            message.user = { id: message.u }
            message.h = this.isUrl(message.m)
          })
          .do((message) => User.inject(message.user))
          .do((message) => messages.push(message))
          .do(() => { shouldScroll = shouldScroll || this.shouldScroll() })
          .do(() => {
            if (this.hidden && !Loader.has('message')) {
              ++this.unread
              Document.setTitle(null, `(${this.unread})`)
            }
          })
          .debounceTime(200)
          .subscribe(
            () => {
              if (Loader.has('message')) Loader.stop('message')
              this.messages = messages
              if (this.paging) {
                this.$nextTick(() => {
                  this.paging = false
                  $(this.$refs.chatBox).scrollTop(this.$refs.chatBox.scrollHeight - this.savedHeight)
                })
              } else if (shouldScroll) {
                shouldScroll = false
                this.scroll(9999)
              }
            }
          )
      },
      shouldScroll () {
        return this.$refs.chatBox.scrollHeight - this.$refs.chatBox.scrollTop <= this.$refs.chatBox.clientHeight + 500
      },
      send () {
        const input = this.input
        if (!input || !input.trim()) return
        this.input = ''
        Me.sendMessage(this.courseId, input).subscribe()
      },
      upload () {
        Document.uploadModal.open()
          .flatMap((file) => Me.sendMessage(this.courseId, file.downloadURL))
          .subscribe(
            null,
            (err) => {
              Document.openErrorModal('Upload Error', err && err.message || err)
            }
          )
      },
      scroll (height) {
        this.$nextTick(() => {
          $(this.$refs.chatBox).scrollTop(this.$refs.chatBox.scrollTop + height)
        })
      }
    }
  }
</script>
