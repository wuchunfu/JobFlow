export default {
  lang(state) {
    if (state.lang === 'en') {
      return 'English'
    } else if (state.lang === 'zh') {
      return '中文'
    } else {
      return state.lang
    }
  }
}
