<template>
  <div>
    <div class="ui segment" style="padding-bottom: 2rem;">
      <div class="ui teal button" @click="openAssignmentModal">Add Assignment</div>
      <div class="ui divider"></div>
      <div class="ui stackable grid">
        <div class="two column row" v-for="x in assignments">
          <div class="column">
            {{ x.title }}
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
  </div>
</template>

<style scoped>
  .action.column > .button {
    width: 120px;
  }
</style>

<script>
  import { Course, Loader } from '../services'

  export default {
    data () {
      return {
        courseId: null,
        assignmentCode: '',
        assignments: null
      }
    },
    created () {
      this.courseId = this.$route.params.id
      Loader.start('assignment')
      Course.getAssignments(this.courseId)
        .subscribe(
          (assignments) => {
            Loader.stop('assignment')
            this.assignments = assignments
          }
        )
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
