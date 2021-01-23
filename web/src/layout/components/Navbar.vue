<template>
  <div class="navbar">
    <hamburger
      :is-active="sidebar.opened"
      class="hamburger-container"
      @toggleClick="toggleSideBar"
    />

    <breadcrumb class="breadcrumb-container"/>

    <div class="right-menu">
      <el-dropdown class="lang-list" trigger="click">
        <span class="el-dropdown-link">
          {{ $t($store.getters['lang/lang']) }}
          <i class="el-icon-arrow-down el-icon--right" />
        </span>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item @click.native="setLang('zh')">
            <span>中文</span>
          </el-dropdown-item>
          <el-dropdown-item @click.native="setLang('en')">
            <span>English</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>

      <el-dropdown class="avatar-container" trigger="click">
        <div class="avatar-wrapper">
          <img :src="'/avatar.jpg?imageView2/1/w/80/h/80'" class="user-avatar">
          <i class="el-icon-arrow-down el-icon--right" />
        </div>
        <el-dropdown-menu slot="dropdown" class="user-dropdown">
          <router-link to="/">
            <el-dropdown-item>
              {{ $t('Home') }}
            </el-dropdown-item>
          </router-link>
          <el-dropdown-item>
            <span
              style="display:block;"
              @click="changePasswordHandle(userId)"
            >
              {{ $t('Change Password') }}
            </span>
          </el-dropdown-item>
          <el-dropdown-item divided @click.native="logout">
            <span style="display:block;">{{ $t('Logout') }}</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>

      <!-- 弹窗, 修改密码 -->
      <change-password
        v-if="changePasswordVisible"
        ref="changePassword"
      />
    </div>
  </div>
</template>

<script>
  import {mapState} from 'vuex'
  import Breadcrumb from '@/components/Breadcrumb'
  import Hamburger from '@/components/Hamburger'
  import ChangePassword from './changePassword/index'

  export default {
    components: {
      Breadcrumb,
      Hamburger,
      ChangePassword
    },
    data() {
      return {
        changePasswordVisible: false
      }
    },
    computed: {
      ...mapState('app', ['sidebar']),
      ...mapState('user', ['userId'])
    },
    methods: {
      toggleSideBar() {
        this.$store.dispatch('app/toggleSideBar')
      },
      async logout() {
        await this.$store.dispatch('user/logout')
        await this.$router.push('/login')
      },
      setLang(lang) {
        window.localStorage.setItem('lang', lang)
        this.$i18n.locale = lang
        this.$store.commit('lang/SET_LANG', lang)
      },
      // 修改密码
      changePasswordHandle(row) {
        console.log(row)
        this.changePasswordVisible = true
        this.$nextTick(() => {
          this.$refs.changePassword.init(row)
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
  .navbar {
    height: 50px;
    overflow: hidden;
    position: relative;
    background: #fff;
    box-shadow: 0 1px 4px rgba(0, 21, 41, .08);

    .hamburger-container {
      line-height: 46px;
      height: 100%;
      float: left;
      cursor: pointer;
      transition: background .3s;
      -webkit-tap-highlight-color: transparent;

      &:hover {
        background: rgba(0, 0, 0, .025)
      }
    }

    .breadcrumb-container {
      float: left;
    }

    .right-menu {
      float: right;
      height: 100%;
      line-height: 50px;

      &:focus {
        outline: none;
      }

      .right-menu-item {
        display: inline-block;
        padding: 0 8px;
        height: 100%;
        font-size: 18px;
        color: #5a5e66;
        vertical-align: text-bottom;

        &.hover-effect {
          cursor: pointer;
          transition: background .3s;

          &:hover {
            background: rgba(0, 0, 0, .025)
          }
        }
      }

      .lang-list {
        cursor: pointer;
        display: inline-block;
        margin-right: 30px;
      }

      .avatar-container {
        margin-right: 30px;

        .avatar-wrapper {
          margin-top: 5px;
          position: relative;

          .user-avatar {
            cursor: pointer;
            width: 40px;
            height: 40px;
            border-radius: 10px;
          }

          .el-icon-caret-bottom {
            cursor: pointer;
            position: absolute;
            right: -20px;
            top: 25px;
            font-size: 12px;
          }
        }
      }
    }
  }
</style>
