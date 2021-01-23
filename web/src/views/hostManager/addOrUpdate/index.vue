<template>
  <el-dialog
    :title="!dataForm.hostId ? '新增' : '修改'"
    :close-on-click-modal="false"
    :visible.sync="visible">
    <el-form
      ref="dataForm"
      :model="dataForm"
      :rules="dataRules"
      label-width="80px"
      @keyup.enter.native="dataFormSubmit()"
    >
      <el-form-item label="主机别名" prop="hostAlias">
        <el-input v-model="dataForm.hostAlias" type="text" style="width: 90%;" placeholder="主机别名"/>
      </el-form-item>
      <el-form-item label="主机名" prop="hostName">
        <el-input v-model="dataForm.hostName" type="text" style="width: 90%;" placeholder="主机名"/>
      </el-form-item>
      <el-form-item label="端口号" prop="hostPort">
        <el-input v-model.number="dataForm.hostPort" type="text" style="width: 90%;" placeholder="端口号"/>
      </el-form-item>
      <el-form-item label="备注" prop="remark">
        <el-input v-model="dataForm.remark" type="text" style="width: 90%;" placeholder="备注"/>
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
      return {
        visible: false,
        dataForm: {
          hostId: 0,
          hostAlias: '',
          hostName: '',
          hostPort: 5921,
          remark: ''
        },
        dataRules: {
          hostAlias: [{required: true, message: '请输入主机别名', trigger: 'blur'}],
          hostName: [{required: true, message: '请输入主机名', trigger: 'blur'}],
          hostPort: [{required: true, message: '请输入端口号', trigger: 'blur'}]
        }
      }
    },
    methods: {
      ...mapActions('hostManager', ['getTableInfo', 'getTableSaveOrUpdate']),

      init(row) {
        this.dataForm.hostId = row !== undefined ? row.hostId || 0 : undefined
        this.visible = true
        if (this.dataForm.hostId) {
          this.getTableInfo({hostId: this.dataForm.hostId}).then((res) => {
            console.log(res.data)
            const result = res.data
            if (result !== undefined && result.code === 200) {
              this.dataForm.hostId = result.data.hostId
              this.dataForm.hostAlias = result.data.hostAlias
              this.dataForm.hostName = result.data.hostName
              this.dataForm.hostPort = result.data.hostPort
              this.dataForm.remark = result.data.remark
            }
          })
        } else {
          this.dataForm = {
            hostId: 0,
            hostAlias: '',
            hostName: '',
            hostPort: 5921,
            remark: ''
          }
        }
      },
      // 表单提交
      dataFormSubmit() {
        this.$refs['dataForm'].validate((valid) => {
          if (valid) {
            this.visible = true
            const params = {
              'hostId': this.dataForm.hostId || undefined,
              'hostAlias': this.dataForm.hostAlias,
              'hostName': this.dataForm.hostName,
              'hostPort': this.dataForm.hostPort,
              'remark': this.dataForm.remark
            }
            console.log(params)
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
                  hostId: 0,
                  hostAlias: '',
                  hostName: '',
                  hostPort: 5921,
                  remark: ''
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
                  hostId: 0,
                  hostAlias: '',
                  hostName: '',
                  hostPort: 5921,
                  remark: ''
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
