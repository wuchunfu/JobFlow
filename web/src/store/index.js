import Vue from 'vue'
import Vuex from 'vuex'
import app from './app'
import settings from './settings'
import user from './user'
import userManager from './userManager'
import hostManager from './hostManager'
import taskManager from './taskManager'
import taskLogManager from './taskLogManager'
import loginLogManager from './loginLogManager'
import lang from './lang'

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    app,
    settings,
    user,
    userManager,
    hostManager,
    taskManager,
    taskLogManager,
    loginLogManager,
    lang
  }
})

export default store
