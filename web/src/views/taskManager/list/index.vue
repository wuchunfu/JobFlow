<template>
  <el-row class="row-box">
    <el-col class="col-box">
      <el-card class="card-box">
        <el-form :inline="true" :model="dataForm" @keyup.enter.native="getDataList()">
          <el-form-item>
            <el-input v-model="dataForm.taskName" placeholder="任务名" clearable/>
          </el-form-item>
          <el-form-item>
            <el-button @click="getDataList()">查询</el-button>
            <el-button type="primary" @click="addOrUpdateHandle()">新增</el-button>
            <el-button
              type="danger"
              :disabled="dataListSelections.length <= 0"
              @click="deleteHandle()"
            >
              批量删除
            </el-button>
          </el-form-item>
        </el-form>
        <el-table
          v-loading="loading"
          :data="dataList"
          element-loading-text="Loading"
          border
          fit
          style="width: 100%;"
          height="620"
          highlight-current-row
          @selection-change="selectionChangeHandle"
        >
          <el-table-column
            type="selection"
            header-align="center"
            align="center"
            width="50"
          />
          <el-table-column
            label="任务名"
            header-align="center"
            align="center"
            fixed
          >
            <template slot-scope="scope">
              {{ scope.row.taskName }}
            </template>
          </el-table-column>
          <el-table-column
            label="标签"
            header-align="center"
            align="center"
            fixed
          >
            <template slot-scope="scope">
              {{ scope.row.taskTag }}
            </template>
          </el-table-column>
          <el-table-column
            label="表达式"
            header-align="center"
            align="center"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.cronExpression }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="下次执行时间"
            header-align="center"
            align="center"
            width="160"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.nextRunTime }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="执行方式"
            header-align="center"
            align="center"
          >
            <template slot-scope="scope">
              {{ scope.row.protocol | protocolConvert }}
            </template>
          </el-table-column>
          <el-table-column
            label="重试次数"
            header-align="center"
            align="center"
          >
            <template slot-scope="scope">
              {{ scope.row.retryTimes }}
            </template>
          </el-table-column>
          <el-table-column
            label="超时时间"
            header-align="center"
            align="center"
          >
            <template slot-scope="scope">
              {{ scope.row.timeout }}
            </template>
          </el-table-column>
          <el-table-column
            label="激活/停止"
            header-align="center"
            align="center"
          >
            <template slot-scope="scope">
              <el-switch
                v-model="scope.row.taskStatus"
                :active-value="1"
                :inactive-vlaue="0"
                active-color="#13ce66"
                inactive-color="#ff4949"
                @change="changeStatus(scope.row)"
              />
            </template>
          </el-table-column>
          <el-table-column
            label="更新时间"
            header-align="center"
            align="center"
            width="155"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.updateTime }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="操作"
            header-align="center"
            align="center"
            width="180"
            fixed="right"
          >
            <template slot-scope="scope">
              <el-button
                type="primary"
                size="small"
                @click="addOrUpdateHandle(scope.row)"
              >
                修改
              </el-button>
              <el-button
                type="primary"
                size="small"
                @click="manualExecHandle(scope.row.taskId)"
              >
                手动执行
              </el-button>
              <el-button
                type="primary"
                size="small"
                @click="stopHandle(scope.row.taskId)"
              >
                手动停止
              </el-button>
              <el-button
                v-loading="loading"
                type="danger"
                size="small"
                @click="deleteHandle(scope.row.taskId)"
              >
                删除
              </el-button>
              <el-button
                type="primary"
                size="small"
                @click="viewLogHandle(scope.row.taskId)"
              >
                查看日志
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          :total="totalPage"
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="pageIndex"
          :page-sizes="[10, 20, 50, 100]"
          :page-size="pageSize"
          @current-change="currentChangeHandle"
          @size-change="sizeChangeHandle"
        />
        <!-- 弹窗, 新增 / 修改 -->
        <add-or-update
          v-if="addOrUpdateVisible"
          ref="addOrUpdate"
          @refreshDataList="getDataList"
        />
      </el-card>
    </el-col>
  </el-row>
</template>

<script>
  import {mapActions} from 'vuex'
  import AddOrUpdate from '../addOrUpdate/index'

  export default {
    filters: {
      protocolConvert(status) {
        const statusMap = {
          1: 'http',
          2: 'shell'
        }
        return statusMap[status]
      }
    },
    components: {
      AddOrUpdate
    },
    data() {
      return {
        loading: true,
        dataForm: {
          taskName: ''
        },
        dataListSelections: [],
        dataList: [],
        addOrUpdateVisible: false,
        pageIndex: 1,
        pageSize: 10,
        totalPage: 0
      }
    },
    created() {
      this.getDataList()
    },
    methods: {
      ...mapActions('taskManager', ['getList', 'getDeleteRow', 'getRun', 'getStop', 'getEnable', 'getDisable']),
      getDataList() {
        this.loading = true
        const params = {
          pageIndex: this.pageIndex,
          pageSize: this.pageSize,
          taskName: this.dataForm.taskName
        }
        this.getList(params).then((res) => {
          // console.log(res.data)
          const result = res.data
          if (result !== undefined && result.code === 200) {
            this.dataList = result.data
            this.totalPage = result.total
          } else {
            this.dataList = []
            this.totalPage = 0
            this.$notify({
              title: '失败',
              showClose: true,
              message: '获取数据失败',
              type: 'error',
              duration: 3000
            })
          }
          this.loading = false
        }).catch((res) => {
          this.$notify({
            title: '失败',
            showClose: true,
            message: '获取数据失败',
            type: 'error',
            duration: 3000
          })
          this.loading = false
        })
      },
      // 新增 / 修改
      addOrUpdateHandle(row) {
        console.log(row)
        this.addOrUpdateVisible = true
        this.$nextTick(() => {
          this.$refs.addOrUpdate.init(row)
        })
      },
      //  删除
      deleteHandle(id) {
        const rowIds = id ? [id] : this.dataListSelections.map(item => {
          return item.taskId
        })
        this.$confirm(`确定要进行 [${id ? '删除' : '批量删除'}] 操作?`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          this.loading = true
          console.log(rowIds)
          this.getDeleteRow({rowIds: rowIds}).then((res) => {
            console.log(res.data)
            const result = res.data
            if (result !== undefined && result.code === 200) {
              this.$notify({
                title: '成功',
                showClose: true,
                message: '删除成功',
                type: 'success',
                duration: 3000
              })
              this.getDataList()
              this.loading = false
            } else {
              this.$notify({
                title: '失败',
                showClose: true,
                message: '删除失败',
                type: 'error',
                duration: 3000
              })
              this.getDataList()
              this.loading = false
            }
          })
        })
      },
      manualExecHandle(id) {
        const taskId = id !== undefined ? id || 0 : undefined
        // this.loading = true
        this.loading = true
        this.getRun({taskId: taskId}).then((res) => {
          // console.log(res.data)
          const result = res.data
          if (result !== undefined && result.code === 200) {
            this.$notify({
              title: '成功',
              showClose: true,
              message: '运行成功',
              type: 'success',
              duration: 3000
            })
            this.getDataList()
            this.loading = false
          } else {
            this.$notify({
              title: '失败',
              showClose: true,
              message: '运行失败',
              type: 'error',
              duration: 3000
            })
            this.getDataList()
            this.loading = false
          }
        })
      },
      stopHandle(id) {
        const taskId = id !== undefined ? id || 0 : undefined
        // this.loading = true
        this.loading = true
        this.getStop({taskId: taskId}).then((res) => {
          // console.log(res.data)
          const result = res.data
          if (result !== undefined && result.code === 200) {
            this.$notify({
              title: '成功',
              showClose: true,
              message: '运行成功',
              type: 'success',
              duration: 3000
            })
            this.getDataList()
            this.loading = false
          } else {
            this.$notify({
              title: '失败',
              showClose: true,
              message: '运行失败',
              type: 'error',
              duration: 3000
            })
            this.getDataList()
            this.loading = false
          }
        })
      },
      changeStatus(item) {
        console.log(item)
        if (item.taskStatus) {
          this.getEnable({taskId: item.taskId})
        } else {
          this.getDisable({taskId: item.taskId})
        }
        this.getDataList()
      },
      viewLogHandle(id) {
        const taskId = id !== undefined ? id || 0 : undefined
        // this.loading = true
        console.log(taskId)
      },
      // 多选
      selectionChangeHandle(val) {
        this.dataListSelections = val
      },
      // 当前页
      currentChangeHandle(val) {
        this.pageIndex = val
        this.getDataList()
      },
      // 每页数
      sizeChangeHandle(val) {
        this.pageSize = val
        this.pageIndex = 1
        this.getDataList()
      }
    }
  }
</script>

<style lang="scss" scoped>
  .row-box {
    height: 100%;

    .col-box {
      width: 100%;
      height: 100%;
    }

    .card-box {
      width: 100%;
      height: 100%;

      .el-pagination {
        margin-top: 20px;
        text-align: right;
      }
    }
  }
</style>
