import i18n from '@/i18n'

/**
 * download file
 */
const downloadFile = ($url, $obj) => {
  const param = {
    url: $url,
    obj: $obj
  }

  if (!param.url) {
    this.$message.warning(`${i18n.$t('Unable to download without proper url')}`)
    return
  }

  const generatorInput = function (obj) {
    let result = ''
    const keyArr = Object.keys(obj)
    keyArr.forEach(function (key) {
      result += "<input type='hidden' name = '" + key + "' value='" + obj[key] + "'>"
    })
    return result
  }
  $(`<form action="${param.url}" method="get">${generatorInput(param.obj)}</form>`).appendTo('body').submit().remove()
}

export {downloadFile}
