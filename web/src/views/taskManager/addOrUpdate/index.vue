<template>
  <div>
    <el-dialog
      :title="!dataForm.taskId ? '新增' : '修改'"
      :close-on-click-modal="false"
      :visible.sync="visible"
    >
      <el-form
        ref="dataForm"
        :model="dataForm"
        :rules="dataRules"
        label-width="165px"
        @keyup.enter.native="dataFormSubmit()"
      >
        <el-form-item label="任务名称" prop="taskName">
          <el-input v-model="dataForm.taskName" type="text" style="width: 90%;" placeholder="任务名称"/>
        </el-form-item>
        <el-form-item label="任务标签" prop="taskTag">
          <el-input v-model="dataForm.taskTag" type="text" style="width: 90%;" placeholder="通过标签将任务分组"/>
        </el-form-item>
        <el-form-item label="任务类型" prop="taskLevel">
          <el-select v-model="dataForm.taskLevel" :disabled="dataForm.taskId !== 0" class="filter-item"
                     style="width: 90%;" placeholder="任务类型">
            <el-option v-for="item in typeOptions" :key="item.value" :label="item.label" :value="item.value"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="dataForm.taskLevel === 1" label="子任务ID" prop="dependencyTaskId">
          <el-input v-model="dataForm.dependencyTaskId" type="text" style="width: 90%;" placeholder="多个ID逗号分隔"/>
        </el-form-item>
        <el-form-item v-if="dataForm.taskLevel === 1" label="依赖关系" prop="dependencyStatus">
          <el-select v-model="dataForm.dependencyStatus" class="filter-item" style="width: 90%;" placeholder="依赖关系">
            <el-option v-for="item in dependencyOptions" :key="item.value" :label="item.label" :value="item.value"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="dataForm.taskLevel === 1" label="crontab表达式" prop="cronExpression">
          <el-input v-model="dataForm.cronExpression" type="text" style="width: 76.5%;" placeholder="秒 分 时 天 月 周"/>
          <el-button
            class="cron-edit"
            type="primary"
            @click="onShowCronDialog"
          >
            Cron
          </el-button>
        </el-form-item>
        <el-form-item label="执行方式" prop="protocol">
          <el-select v-model="dataForm.protocol" class="filter-item" style="width: 90%;" placeholder="执行方式">
            <el-option v-for="item in protocolOptions" :key="item.value" :label="item.label" :value="item.value"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="dataForm.protocol !== 1" label="任务节点" prop="selectedHostByHostIds">
          <el-select v-model="selectedHostByHostIds" filterable multiple class="filter-item" style="width: 90%;"
                     placeholder="任务节点">
            <el-option v-for="item in hostsOptions" :key="item.hostId" :label="item.hostAlias + ' - ' + item.hostName"
                       :value="item.hostId"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="dataForm.protocol === 1" label="请求方法" prop="httpMethod">
          <el-select v-model="dataForm.httpMethod" class="filter-item" style="width: 90%;" placeholder="请求方法">
            <el-option v-for="item in httpMethodOptions" :key="item.value" :label="item.label" :value="item.value"/>
          </el-select>
        </el-form-item>
        <el-form-item label="命令" prop="command">
          <el-input v-model="dataForm.command" type="textarea" style="width: 90%;"
                    :placeholder="commandPlaceholder"/>
        </el-form-item>
        <el-form-item label="任务超时时间" prop="timeout">
          <el-input v-model.number.trim="dataForm.timeout" type="text" style="width: 90%;" placeholder="任务超时时间"/>
        </el-form-item>
        <el-form-item label="单实例运行" prop="isMultiInstance">
          <el-select v-model.trim="dataForm.isMultiInstance" class="filter-item" style="width: 90%;"
                     placeholder="单实例运行">
            <el-option v-for="item in isMultiInstanceOptions" :key="item.value" :label="item.label"
                       :value="item.value"/>
          </el-select>
        </el-form-item>
        <el-form-item label="任务失败重试次数" prop="retryTimes">
          <el-input v-model.number.trim="dataForm.retryTimes" type="text" style="width: 90%;"
                    placeholder="0 - 10, 默认0，不重试"/>
        </el-form-item>
        <el-form-item label="任务失败重试间隔时间" prop="retryInterval">
          <el-input v-model.number.trim="dataForm.retryInterval" type="text" style="width: 90%;"
                    placeholder="0 - 3600 (秒), 默认0，执行系统默认策略"/>
        </el-form-item>
        <el-form-item label="任务通知" prop="notifyStatus">
          <el-select v-model="dataForm.notifyStatus" class="filter-item" style="width: 90%;" placeholder="任务通知">
            <el-option v-for="item in notifyStatusOptions" :key="item.value" :label="item.label" :value="item.value"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="dataForm.notifyStatus !== 1" label="通知类型">
          <el-select v-model="dataForm.notifyType" class="filter-item" style="width: 90%;" placeholder="通知类型">
            <el-option v-for="item in notifyTypeOptions" :key="item.value" :label="item.label" :value="item.value"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="dataForm.notifyStatus !== 1 && dataForm.notifyType === 2" label="接收用户"
                      prop="selectedMailNotifyIds">
          <el-select v-model="selectedMailNotifyIds" filterable multiple class="filter-item" style="width: 90%;"
                     placeholder="接收用户">
            <el-option v-for="item in mailUsers" :key="item.id" :label="item.username" :value="item.id"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="dataForm.notifyStatus !== 1 && dataForm.notifyType === 3" label="发送Channel"
                      prop="selectedSlackNotifyIds">
          <el-select v-model="selectedSlackNotifyIds" filterable multiple class="filter-item" style="width: 90%;"
                     placeholder="发送Channel">
            <el-option v-for="item in slackChannels" :key="item.id" :label="item.username" :value="item.id"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="dataForm.notifyStatus === 4" label="任务执行输出关键字" prop="notifyKeyword">
          <el-input v-model="dataForm.notifyKeyword" type="textarea" style="width: 90%;"
                    placeholder="任务执行输出中包含此关键字将触发通知"/>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="dataForm.remark" type="textarea" style="width: 90%;" placeholder="备注"/>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="visible=false">取消</el-button>
        <el-button type="primary" @click="dataFormSubmit()">保存</el-button>
      </span>
    </el-dialog>

    <!--cron generation dialog-->
    <el-dialog title="生成 Cron" :visible.sync="cronDialogVisible">
      <vue-cron-linux
        ref="vue-cron-linux"
        :data="dataForm.cronExpression"
        :i18n="lang"
        @submit="onCronChange"
      />
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="cronDialogVisible = false">取消</el-button>
        <el-button size="small" type="primary" @click="onCronDialogSubmit">确认</el-button>
      </span>
    </el-dialog>
    <!--./cron generation dialog-->
  </div>
</template>

<script>
  import {mapActions} from 'vuex'
  import VueCronLinux from '@/components/Cron'

  export default {
    components: {
      VueCronLinux
    },
    data() {
      return {
        visible: false,
        dataForm: {
          taskId: 0,
          hostId: 0,
          taskName: '',
          taskTag: '',
          taskLevel: 1,
          dependencyTaskId: '',
          dependencyStatus: 1,
          cronExpression: '',
          protocol: 2,
          httpMethod: 1,
          command: '',
          timeout: 0,
          isMultiInstance: 2,
          retryTimes: 0,
          retryInterval: 0,
          notifyStatus: 1,
          notifyType: 2,
          notifyReceiverId: '',
          notifyKeyword: '',
          taskRemark: ''
        },
        typeOptions: [
          {
            label: '主任务',
            value: 1
          },
          {
            label: '子任务',
            value: 2
          }
        ],
        dependencyOptions: [
          {
            label: '强依赖',
            value: 1
          },
          {
            label: '弱依赖',
            value: 2
          }
        ],
        protocolOptions: [
          {
            label: 'http',
            value: 1
          },
          {
            label: 'shell',
            value: 2
          }
        ],
        httpMethodOptions: [
          {
            label: 'get',
            value: 1
          },
          {
            label: 'post',
            value: 2
          }
        ],
        isMultiInstanceOptions: [
          {
            label: '否',
            value: 1
          },
          {
            label: '是',
            value: 2
          }
        ],
        notifyStatusOptions: [
          {
            label: '不通知',
            value: 1
          },
          {
            label: '失败通知',
            value: 2
          },
          {
            label: '总是通知',
            value: 3
          },
          {
            label: '关键词匹配通知',
            value: 4
          }
        ],
        notifyTypeOptions: [
          {
            label: 'Mail',
            value: 2
          },
          {
            label: 'Slack',
            value: 3
          },
          {
            label: 'WebHook',
            value: 4
          }
        ],
        taskStatusOptions: [
          {
            label: '停止',
            value: 1
          },
          {
            label: '正常',
            value: 2
          }
        ],
        mailUsers: [],
        selectedMailNotifyIds: [],
        selectedHostByHostIds: [],
        hostsOptions: [],
        slackChannels: [],
        selectedSlackNotifyIds: [],
        cronDialogVisible: false,
        dataRules: {
          taskName: [{required: true, message: '请输入任务名称', trigger: 'blur'}],
          cronExpression: [{required: true, message: '请输入 crontab 表达式', trigger: 'blur'}],
          command: [{required: true, message: '请输入执行命令', trigger: 'blur'}],
          timeout: [{required: true, message: '请输入超时时间', trigger: 'blur'}],
          retryTimes: [{required: true, message: '请输入失败重试次数', trigger: 'blur'}],
          retryInterval: [{required: true, message: '请输入失败重试时间间隔', trigger: 'blur'}],
          notifyKeyword: [{required: true, message: '请输入关键字', trigger: 'blur'}]
        }
      }
    },
    computed: {
      lang() {
        const lang = this.$store.state.lang.lang || window.localStorage.getItem('lang')
        if (!lang) return 'cn'
        if (lang === 'zh') return 'cn'
        return 'en'
      },
      commandPlaceholder() {
        if (this.dataForm.protocol === 1) {
          return '请输入URL地址'
        }
        return '请输入shell命令'
      }
    },
    created() {
      this.getAllHostDataList()
    },
    methods: {
      ...mapActions('taskManager', ['getTableInfo', 'getTableSaveOrUpdate', 'getHostAllList']),

      init(row) {
        this.dataForm.taskId = row !== undefined ? row.taskId || 0 : undefined
        this.visible = true
        if (this.dataForm.taskId) {
          this.getTableInfo({taskId: this.dataForm.taskId}).then((res) => {
            console.log(res.data)
            const result = res.data
            if (result !== undefined && result.code === 200) {
              this.dataForm.taskId = result.data.taskId
              this.dataForm.hostId = result.data.hostId
              this.dataForm.taskName = result.data.taskName
              this.dataForm.taskLevel = result.data.taskLevel
              this.dataForm.dependencyTaskId = result.data.dependencyTaskId
              this.dataForm.dependencyStatus = result.data.dependencyStatus
              this.dataForm.cronExpression = result.data.cronExpression
              this.dataForm.protocol = result.data.protocol
              this.dataForm.httpMethod = result.data.httpMethod
              this.dataForm.command = result.data.command
              this.dataForm.timeout = result.data.timeout
              this.dataForm.isMultiInstance = result.data.isMultiInstance
              this.dataForm.retryTimes = result.data.retryTimes
              this.dataForm.notifyStatus = result.data.notifyStatus + 1
              this.dataForm.notifyType = result.data.notifyType + 1
              this.dataForm.notifyReceiverId = result.data.notifyReceiverId
              this.dataForm.notifyKeyword = result.data.notifyKeyword
              this.dataForm.taskTag = result.data.taskTag
              this.dataForm.taskRemark = result.data.taskRemark

              if (result.data.dependencyStatus) {
                this.dataForm.dependencyStatus = result.data.dependencyStatus
              }
              if (result.data.httpMethod) {
                this.dataForm.httpMethod = result.data.httpMethod
              }
              this.hostsOptions = result.data.hosts || []
              if (this.dataForm.protocol === 2) {
                this.hostsOptions.forEach((v) => {
                  this.selectedHostByHostIds.push(v.hostId)
                })
              }
              if (this.dataForm.notifyStatus > 1) {
                const notifyReceiverIds = this.dataForm.notifyReceiverId.split(',')
                if (this.dataForm.notifyType === 2) {
                  notifyReceiverIds.forEach((v) => {
                    this.selectedMailNotifyIds.push(parseInt(v))
                  })
                } else if (this.dataForm.notifyType === 3) {
                  notifyReceiverIds.forEach((v) => {
                    this.selectedSlackNotifyIds.push(parseInt(v))
                  })
                }
              }
              // notificationService.mail((data) => {
              //   this.mailUsers = data.mail_users
              // })
              //
              // notificationService.slack((data) => {
              //   this.slackChannels = data.channels
              // })
            }
          })
        } else {
          this.dataForm = {
            taskId: 0,
            hostId: 0,
            taskName: '',
            taskTag: '',
            taskLevel: 1,
            dependencyTaskId: '',
            dependencyStatus: 1,
            cronExpression: '',
            protocol: 2,
            httpMethod: 1,
            command: '',
            timeout: 0,
            isMultiInstance: 2,
            retryTimes: 0,
            retryInterval: 0,
            notifyStatus: 1,
            notifyType: 2,
            notifyReceiverId: '',
            notifyKeyword: '',
            taskRemark: ''
          }
        }
      },
      getAllHostDataList() {
        this.getHostAllList().then((res) => {
          const result = res.data
          // console.log(111111)
          // console.log(result)
          if (result !== undefined && result.code === 200) {
            this.hostsOptions = result.data || []
            console.log(this.dataForm.protocol)
            // if (this.dataForm.protocol === 2) {
            //   this.hostsOptions.forEach((v) => {
            //     this.selectedHostByHostIds.push(v.hostId)
            //   })
            // }
          } else {
            this.hostsOptions = []
            this.$notify({
              title: '失败',
              showClose: true,
              message: '获取数据失败',
              type: 'error',
              duration: 3000
            })
          }
        }).catch((res) => {
          console.log(res)
          this.$notify({
            title: '失败',
            showClose: true,
            message: '获取数据失败',
            type: 'error',
            duration: 3000
          })
        })
      },
      // 表单提交
      dataFormSubmit() {
        this.$refs['dataForm'].validate((valid) => {
          if (valid) {
            this.visible = true
            if (this.dataForm.protocol === 2 && this.selectedHostByHostIds.length === 0) {
              this.$message.error('请选择任务节点')
              return false
            }
            if (this.dataForm.notifyStatus > 1) {
              if (this.dataForm.notifyType === 2 && this.selectedMailNotifyIds.length === 0) {
                this.$message.error('请选择邮件接收用户')
                return false
              }
              if (this.dataForm.notifyType === 3 && this.selectedSlackNotifyIds.length === 0) {
                this.$message.error('请选择Slack Channel')
                return false
              }
            }

            if (this.dataForm.protocol === 2 && this.selectedHostByHostIds.length > 0) {
              this.dataForm.hostId = this.selectedHostByHostIds.join(',')
            }
            if (this.dataForm.notifyStatus > 1 && this.form.notifyType === 2) {
              this.dataForm.notifyReceiverId = this.selectedMailNotifyIds.join(',')
            }
            if (this.dataForm.notifyStatus > 1 && this.form.notifyType === 3) {
              this.dataForm.notifyReceiverId = this.selectedSlackNotifyIds.join(',')
            }

            const params = {
              'taskId': this.dataForm.taskId || undefined,
              'hostId': this.dataForm.hostId || undefined,
              'taskName': this.dataForm.taskName,
              'taskLevel': this.dataForm.taskLevel,
              'dependencyTaskId': this.dataForm.dependencyTaskId,
              'dependencyStatus': this.dataForm.dependencyStatus,
              'cronExpression': this.dataForm.cronExpression,
              'protocol': this.dataForm.protocol,
              'httpMethod': this.dataForm.httpMethod,
              'command': this.dataForm.command,
              'timeout': this.dataForm.timeout,
              'isMultiInstance': this.dataForm.isMultiInstance,
              'retryTimes': this.dataForm.retryTimes,
              'notifyStatus': this.dataForm.notifyStatus,
              'notifyType': this.dataForm.notifyType,
              'notifyReceiverId': this.dataForm.notifyReceiverId,
              'notifyKeyword': this.dataForm.notifyKeyword,
              'taskTag': this.dataForm.taskTag,
              'taskRemark': this.dataForm.taskRemark
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
                  taskId: 0,
                  hostId: 0,
                  taskName: '',
                  taskLevel: 1,
                  dependencyTaskId: '',
                  dependencyStatus: 1,
                  cronExpression: '',
                  protocol: 2,
                  httpMethod: 1,
                  command: '',
                  timeout: 0,
                  isMultiInstance: 2,
                  retryTimes: 0,
                  retryInterval: 0,
                  notifyStatus: 1,
                  notifyType: 2,
                  notifyReceiverId: '',
                  notifyKeyword: '',
                  taskTag: '',
                  taskRemark: ''
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
                  taskId: 0,
                  hostId: 0,
                  taskName: '',
                  taskLevel: 1,
                  dependencyTaskId: '',
                  dependencyStatus: 1,
                  cronExpression: '',
                  protocol: 2,
                  httpMethod: 1,
                  command: '',
                  timeout: 0,
                  isMultiInstance: 2,
                  retryTimes: 0,
                  retryInterval: 0,
                  notifyStatus: 1,
                  notifyType: 2,
                  notifyReceiverId: '',
                  notifyKeyword: '',
                  taskTag: '',
                  taskRemark: ''
                }
                this.visible = false
                this.$emit('refreshDataList')
              }
            })
          }
        })
      },
      onShowCronDialog() {
        this.cronDialogVisible = true
      },
      onCronChange(value) {
        console.log(value)
        this.dataForm.cronExpression = value
      },
      onCronDialogSubmit() {
        const valid = this.$refs['vue-cron-linux'].submit()
        if (valid) {
          this.cronDialogVisible = false
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
  .cron-edit {
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
  }
</style>
