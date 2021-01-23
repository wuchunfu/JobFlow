<template>
  <el-dialog
    :title="!dataForm.userId ? '新增' : '修改'"
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
        <el-input v-model="dataForm.username" type="text" style="width: 90%;" placeholder="用户名"/>
      </el-form-item>
      <el-form-item v-if="!dataForm.userId" label="密码" prop="password">
        <el-input v-model="dataForm.password" type="password" autocomplete="off" style="width: 90%;" placeholder="密码"/>
      </el-form-item>
      <el-form-item v-if="!dataForm.userId" label="确认密码" prop="confirmPassword">
        <el-input v-model="dataForm.confirmPassword" type="password" autocomplete="off" style="width: 90%;"
                  placeholder="确认密码"/>
      </el-form-item>
      <el-form-item label="邮箱" prop="email">
        <el-input v-model="dataForm.email" type="text" style="width: 90%;" placeholder="邮箱"/>
      </el-form-item>
      <el-form-item label="角色" prop="isAdmin">
        <el-select v-model="dataForm.isAdmin" class="filter-item" style="width: 90%;" placeholder="请选择">
          <el-option v-for="item in roleOptions" :key="item.value" :label="item.label" :value="item.value"/>
        </el-select>
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="dataForm.status" class="filter-item" style="width: 90%;" placeholder="请选择">
          <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value"/>
        </el-select>
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
        if (value === '') {
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
        if (value === '') {
          callback(new Error('请再次输入密码'))
        } else if (value !== this.dataForm.password) {
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
          confirmPassword: '',
          email: '',
          isAdmin: 0,
          status: 1
        },
        roleOptions: [
          {
            label: '普通用户',
            value: 0
          },
          {
            label: '管理员',
            value: 1
          }
        ],
        statusOptions: [
          {
            label: '禁用',
            value: 0
          },
          {
            label: '正常',
            value: 1
          }
        ],
        dataRules: {
          username: [{required: true, message: '请输入用户名', trigger: 'blur'}],
          password: [{required: true, trigger: 'blur', validator: validatePassword}],
          confirmPassword: [{type: 'password', required: true, trigger: ['blur'], validator: validateConfirmPassword}],
          email: [{type: 'email', required: true, message: '请输入正确的邮箱地址', trigger: ['blur']}],
          isAdmin: [{required: true, message: '角色不能为空', trigger: 'blur'}],
          status: [{required: true, message: '状态不能为空', trigger: 'blur'}]
        }
      }
    },
    methods: {
      ...mapActions('userManager', ['getTableInfo', 'getTableSaveOrUpdate']),

      init(row) {
        this.dataForm.userId = row !== undefined ? row.userId || 0 : undefined
        this.visible = true
        if (this.dataForm.userId) {
          this.getTableInfo({userId: this.dataForm.userId}).then((res) => {
            console.log(res.data)
            const result = res.data
            if (result !== undefined && result.code === 200) {
              this.dataForm.userId = result.data.userId
              this.dataForm.username = result.data.username
              this.dataForm.password = result.data.password
              this.dataForm.email = result.data.email
              this.dataForm.isAdmin = result.data.isAdmin
              this.dataForm.status = result.data.status
            }
          })
        } else {
          this.dataForm = {
            userId: 0,
            username: '',
            password: '',
            confirmPassword: '',
            email: '',
            isAdmin: 0,
            status: 1
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
              'email': this.dataForm.email,
              'isAdmin': this.dataForm.isAdmin,
              'status': this.dataForm.status
            }
            this.getTableSaveOrUpdate(params).then((res) => {
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
                  confirmPassword: '',
                  email: '',
                  isAdmin: 0,
                  status: 1
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
                  confirmPassword: '',
                  email: '',
                  isAdmin: 0,
                  status: 1
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
