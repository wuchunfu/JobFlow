/**
 * @param {string} path
 * @returns {Boolean}
 */
export function isExternal(path) {
  return /^(https?:|mailto:|tel:)/.test(path)
}

/**
 * @param {string} str
 * @returns {Boolean}
 */
export function validUsername(str) {
  const arr = ['&', '\\', '/', '*', '>', '<', '@', '!', '#']
  if (str.trim().length >= 0) {
    for (let i = 0; i < arr.length; i++) {
      for (let j = 0; j < str.length; j++) {
        if (arr[i] === str.charAt(j)) {
          return false
        }
      }
    }
  }
  return true
}
