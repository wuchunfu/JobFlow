export default {
  CHANGE_SETTING: (state, {key, value}) => {
    if (Object.prototype.hasOwnProperty.call(state, key)) {
      state[key] = value
    }
  }
}
