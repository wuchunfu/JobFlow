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
            label="任务ID"
            header-align="center"
            align="center"
            fixed
            width="70"
          >
            <template slot-scope="scope">
              {{ scope.row.taskId }}
            </template>
          </el-table-column>
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
            label="crontab表达式"
            header-align="center"
            align="center"
            fixed
            width="80"
          >
            <template slot-scope="scope">
              {{ scope.row.cronExpression }}
            </template>
          </el-table-column>
          <el-table-column
            label="执行方式"
            header-align="center"
            align="center"
            width="80"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.protocol | protocolConvert }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="命令"
            header-align="center"
            align="center"
          >
            <template slot-scope="scope">
              {{ scope.row.command }}
            </template>
          </el-table-column>
          <el-table-column
            label="超时时间(秒)"
            header-align="center"
            align="center"
            width="80"
          >
            <template slot-scope="scope">
              {{ scope.row.timeout }}
            </template>
          </el-table-column>
          <el-table-column
            label="重试次数"
            header-align="center"
            align="center"
            width="80"
          >
            <template slot-scope="scope">
              {{ scope.row.retryTimes }}
            </template>
          </el-table-column>
          <el-table-column
            label="开始时间"
            header-align="center"
            align="center"
            width="155"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.startTime }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="结束时间"
            header-align="center"
            align="center"
            width="155"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.endTime }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="总耗时(秒)"
            header-align="center"
            align="center"
            width="70"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.totalTime }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="状态"
            header-align="center"
            align="center"
            width="50"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.status | statusConvert }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="结果"
            header-align="center"
            align="center"
            fixed="right"
            width="300"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.result }}</span>
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
      </el-card>
    </el-col>
  </el-row>
</template>

<script>
  import {mapActions} from 'vuex'

  export default {
    filters: {
      protocolConvert(status) {
        const statusMap = {
          1: 'http',
          2: 'shell'
        }
        return statusMap[status]
      },
      statusConvert(status) {
        const statusMap = {
          0: '失败',
          1: '执行中',
          2: '成功',
          3: '取消'
        }
        return statusMap[status]
      }
    },
    components: {},
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
      ...mapActions('taskLogManager', ['getList', 'getDeleteRow']),
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
      //  删除
      deleteHandle(id) {
        const rowIds = id ? [id] : this.dataListSelections.map(item => {
          return item.hostId
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
