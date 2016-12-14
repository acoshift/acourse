<template lang="pug">
  div
    .ui.segment(style="padding-bottom: 2rem;")
      .ui.teal.button(@click="openAssignmentModal") Add Assignment
      .ui.divider
      .ui.stackable.grid
        .two.column.row(v-for="x in assignments")
          .column
            | {{ x.title }}
            span ({{ x.users.length }})
          .action.column
            .ui.red.button(v-if="x.open", @click="closeAssignment(x)") Close
            .ui.blue.button(v-else, @click="openAssignment(x)") Open
            .ui.green.button Edit
    .ui.small.modal(ref="assignmentModal")
      .header Add Assignment
      .content
        .ui.form
          .field
            label Title
            input(v-model="assignmentCode")
          .ui.fluid.blue.button(@click="submitAssignmentCode") OK
    .ui.segment
      .ui.styled.fluid.accordion
        div(v-for="x in assignments")
          .title
            .dropdown.icon
            | {{ x.title }}
            span ({{ x.users.length }})
          .ui.content
            .ui.stackable.grid
              .two.column.row(v-for="u in x.users")
                .column
                  user-avatar(:user="u")
                .column
                  div(v-for="(a, f) in u.files")
                    span {{ a.timestamp | date('YYYY/MM/DD HH:mm') }}
                    a(:href="a.url", target="_blank") {{ f }}
</template>

<style scoped>
  .action.column > .button {
    width: 120px;
  }
</style>

<script>
  import { Loader, User, Assignment } from '../services'
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
        courseId: this.$route.params.id,
        assignmentCode: '',
        assignments: null
      }
    },
    created () {
      Loader.start('assignment')
      this.$subscribeTo(Assignment.get(this.courseId),
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
        })
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
        Assignment.addCode(this.courseId, { title: this.assignmentCode })
          .subscribe(
            () => {
              $(this.$refs.assignmentModal).modal('hide')
            }
          )
      },
      openAssignment (x) {
        Assignment.open(this.courseId, x.id, true).subscribe()
      },
      closeAssignment (x) {
        Assignment.open(this.courseId, x.id, false).subscribe()
      }
    }
  }
</script>
