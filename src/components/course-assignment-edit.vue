<template>
  <div>
    <div class="ui segment" style="padding-bottom: 2rem;">
      <div class="ui teal button" @click="openAssignmentModal">Add Assignment</div>
      <div class="ui divider"></div>
      <div class="ui stackable grid">
        <div class="two column row" v-for="x in assignments">
          <div class="column">
            {{ x.title }} <span>({{ x.users.length }})</span>
          </div>
          <div class="action column">
            <div v-if="x.open" class="ui red button" @click="closeAssignment(x)">Close</div>
            <div v-else class="ui blue button" @click="openAssignment(x)">Open</div>
            <div class="ui green button">Edit</div>
          </div>
        </div>
      </div>
    </div>
    <div class="ui small modal" ref="assignmentModal">
      <div class="header">Add Assignment</div>
      <div class="content">
        <div class="ui form">
          <div class="field">
            <label>Title</label>
            <input v-model="assignmentCode">
          </div>
          <div class="ui fluid blue button" @click="submitAssignmentCode">OK</div>
        </div>
      </div>
    </div>
    <div class="ui segment">
      <div class="ui styled fluid accordion">
        <div v-for="x in assignments">
          <div class="title">
            <div class="dropdown icon"></div>
            {{ x.title }} <span>({{ x.users.length }})</span>
          </div>
          <div class="ui content">
            <div class="ui stackable grid">
              <div class="two column row" v-for="u in x.users">
                <div class="column">
                  <user-avatar :user="u"></user-avatar>
                </div>
                <div class="column">
                  <div v-for="(a, f) in u.files">
                    <span>{{ a.timestamp | date('YYYY/MM/DD HH:mm') }}</span>
                    <a :href="a.url" target="_blank">{{ f }}</a>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .action.column > .button {
    width: 120px;
  }
</style>

<script>
  import { Course, Loader, User } from '../services'
  import UserAvatar from './user-avatar'
  import flow from 'lodash/fp/flow'
  import map from 'lodash/fp/map'
  import keys from 'lodash/fp/keys'
  import forEach from 'lodash/fp/forEach'
  import isEmpty from 'lodash/fp/isEmpty'
  import filter from 'lodash/fp/filter'

  export default {
    components: {
      UserAvatar
    },
    data () {
      return {
        courseId: null,
        assignmentCode: '',
        assignments: null
      }
    },
    mounted () {
      this.courseId = this.$route.params.id
      Loader.start('assignment')
      Course.assignments(this.courseId)
        .subscribe(
          (assignments) => {
            Loader.stop('assignment')
            this.assignments = flow(
              keys,
              map((id) => ({
                id,
                ...assignments.code[id],
                users: flow(
                  keys,
                  map((x) => ({
                    id: x,
                    files: assignments.user[x][id]
                  })),
                  filter((x) => !isEmpty(x.files)),
                  forEach((x) => User.inject(x))
                )(assignments.user)
              }))
            )(assignments.code)
            setTimeout(() => {
              this.$nextTick(() => {
                $('.accordion').accordion()
              })
            }, 500)
          }
        )
    },
    updated () {
      this.$nextTick(() => {
        $('.accordion').accordion()
      })
    },
    methods: {
      openAssignmentModal () {
        $(this.$refs.assignmentModal).modal('show')
      },
      submitAssignmentCode () {
        Course.addAssignment(this.courseId, { title: this.assignmentCode })
          .subscribe(
            () => {
              $(this.$refs.assignmentModal).modal('hide')
            }
          )
      },
      openAssignment (x) {
        Course.openAssignment(this.courseId, x.id, true).subscribe()
      },
      closeAssignment (x) {
        Course.openAssignment(this.courseId, x.id, false).subscribe()
      }
    }
  }
</script>
