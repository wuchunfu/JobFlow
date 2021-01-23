<template>
  <el-dialog
    :title="'修改密码'"
    :close-on-click-modal="false"
    :visible.sync="visible">
    <el-form
      ref="dataForm"
      :model="dataForm"
      :rules="dataRules"
      label-width="80px"
      @keyup.enter.native="dataFormSubmit()"
    >
      <el-form-item label="用户名" prop="username">
        <el-input v-model="dataForm.username" type="text" disabled style="width: 90%;" placeholder="用户名"/>
      </el-form-item>
      <el-form-item label="原密码" prop="password">
        <el-input v-model="dataForm.password" type="password" autocomplete="off" style="width: 90%;" placeholder="原密码"/>
      </el-form-item>
      <el-form-item label="新密码" prop="newPassword">
        <el-input v-model="dataForm.newPassword" type="password" autocomplete="off" style="width: 90%;"
                  placeholder="新密码"/>
      </el-form-item>
      <el-form-item label="确认密码" prop="confirmPassword">
        <el-input v-model="dataForm.confirmPassword" type="password" autocomplete="off" style="width: 90%;"
                  placeholder="确认密码"/>
      </el-form-item>
    </el-form>
    <span slot="footer" class="dialog-footer">
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="dataFormSubmit()">保存</el-button>
    </span>
  </el-dialog>
</template>

<script>
  import {mapActions} from 'vuex'

  export default {
    data() {
      const validatePassword = (rule, value, callback) => {
        if (value === '' || value === undefined) {
          callback(new Error('请输入密码'))
        } else if (value.length < 6) {
          callback(new Error('密码不能少于6位数'))
        } else {
          callback()
        }
      }
      const validateNewPassword = (rule, value, callback) => {
        if (value === '' || value === undefined) {
          callback(new Error('请输入密码'))
        } else if (value.length < 6) {
          callback(new Error('密码不能少于6位数'))
        } else {
          if (this.dataForm.confirmPassword !== '') {
            this.$refs.dataForm.validateField('confirmPassword')
          }
          callback()
        }
      }
      const validateConfirmPassword = (rule, value, callback) => {
        if (value === '' || value === undefined) {
          callback(new Error('请再次输入密码'))
        } else if (value !== this.dataForm.newPassword) {
          callback(new Error('两次输入密码不一致!'))
        } else if (value.length < 6) {
          callback(new Error('密码不能少于6位数'))
        } else {
          callback()
        }
      }
      return {
        visible: false,
        dataForm: {
          userId: 0,
          username: '',
          password: '',
          newPassword: '',
          confirmPassword: ''
        },
        dataRules: {
          username: [{required: true, message: '请输入用户名', trigger: 'blur'}],
          password: [{required: true, trigger: 'blur', validator: validatePassword}],
          newPassword: [{required: true, trigger: ['blur'], validator: validateNewPassword}],
          confirmPassword: [{required: true, trigger: ['blur'], validator: validateConfirmPassword}]
        }
      }
    },
    methods: {
      ...mapActions('userManager', ['getTableInfo', 'getChangeLoginPassword']),

      init(id) {
        this.dataForm.userId = id !== undefined ? id || 0 : undefined
        this.visible = true
        if (this.dataForm.userId) {
          this.getTableInfo({userId: this.dataForm.userId}).then((res) => {
            console.log(11111)
            console.log(res.data)
            const result = res.data
            if (result !== undefined && result.code === 200) {
              this.dataForm.userId = result.data.userId
              this.dataForm.username = result.data.username
              this.dataForm.password = result.data.password
              this.dataForm.newPassword = result.data.newPassword
            }
          })
        } else {
          this.dataForm = {
            userId: 0,
            username: '',
            password: '',
            newPassword: '',
            confirmPassword: ''
          }
        }
      },
      // 表单提交
      dataFormSubmit() {
        this.$refs['dataForm'].validate((valid) => {
          if (valid) {
            this.visible = true
            const params = {
              'userId': this.dataForm.userId || undefined,
              'username': this.dataForm.username,
              'password': this.dataForm.password,
              'newPassword': this.dataForm.newPassword
            }
            this.getChangeLoginPassword(params).then((res) => {
              // console.log(res.data)
              const result = res.data
              if (result !== undefined && result.code === 200) {
                this.$notify({
                  title: '成功',
                  showClose: true,
                  message: '操作成功',
                  type: 'success',
                  duration: 3000
                })
                this.dataForm = {
                  userId: 0,
                  username: '',
                  password: '',
                  newPassword: '',
                  confirmPassword: ''
                }
                this.visible = false
                this.$emit('refreshDataList')
              } else {
                this.$notify({
                  title: '失败',
                  showClose: true,
                  message: '操作失败',
                  type: 'error',
                  duration: 3000
                })
                this.dataForm = {
                  userId: 0,
                  username: '',
                  password: '',
                  newPassword: '',
                  confirmPassword: ''
                }
                this.visible = false
                this.$emit('refreshDataList')
              }
            })
          }
        })
      }
    }
  }
</script>
