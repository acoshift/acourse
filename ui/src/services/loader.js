import keys from 'lodash/fp/keys'
import flow from 'lodash/fp/flow'
import filter from 'lodash/fp/filter'

export default {
  state: {
    value: 0
  },
  index: {},
  reIndex () {
    this.state.value = flow(
      filter((x) => x > 0),
      keys
    )(this.index).length
  },
  start (val) {
    this.index[val] = (this.index[val] || 0) + 1
    if (this.index[val] === 0) delete this.index[val]
    this.reIndex()
  },
  stop (val) {
    this.index[val] = (this.index[val] || 0) - 1
    if (this.index[val] === 0) delete this.index[val]
    this.reIndex()
  },
  reset () {
    this.index = {}
    this.reIndex()
  },
  has (val) {
    return !!this.index[val]
  }
}
