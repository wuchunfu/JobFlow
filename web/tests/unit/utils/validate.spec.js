import {validUsername, isExternal} from '@/utils/validate.js'

describe('Utils:validate', () => {
  it('validUsername', () => {
    expect(validUsername('admin')).toBe(true)
    expect(validUsername('xx@xx')).toBe(false)
  })
  it('isExternal', () => {
    expect(isExternal('https://github.com/wuchunfu/vue-admin')).toBe(true)
    expect(isExternal('http://github.com/wuchunfu/vue-admin')).toBe(true)
    expect(isExternal('github.com/wuchunfu/vue-admin')).toBe(false)
    expect(isExternal('/dashboard')).toBe(false)
    expect(isExternal('./dashboard')).toBe(false)
    expect(isExternal('dashboard')).toBe(false)
  })
})
