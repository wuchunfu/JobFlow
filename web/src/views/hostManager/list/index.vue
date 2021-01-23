<template>
  <el-row class="row-box">
    <el-col class="col-box">
      <el-card class="card-box">
        <el-form :inline="true" :model="dataForm" @keyup.enter.native="getDataList()">
          <el-form-item>
            <el-input v-model="dataForm.hostName" placeholder="主机名" clearable/>
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
            label="主机别名"
            header-align="center"
            align="center"
            fixed
          >
            <template slot-scope="scope">
              {{ scope.row.hostAlias }}
            </template>
          </el-table-column>
          <el-table-column
            label="主机名"
            header-align="center"
            align="center"
            fixed
          >
            <template slot-scope="scope">
              {{ scope.row.hostName }}
            </template>
          </el-table-column>
          <el-table-column
            label="端口号"
            header-align="center"
            align="center"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.hostPort }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="备注"
            header-align="center"
            align="center"
          >
            <template slot-scope="scope">
              {{ scope.row.remark }}
            </template>
          </el-table-column>
          <el-table-column
            label="更新时间"
            header-align="center"
            align="center"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.updateTime }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="操作"
            header-align="center"
            align="center"
            width="235"
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
                @click="connectTestHandle(scope.row.hostId)"
              >
                连接测试
              </el-button>
              <el-button
                type="danger"
                size="small"
                @click="deleteHandle(scope.row.hostId)"
              >
                删除
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
    components: {
      AddOrUpdate
    },
    data() {
      return {
        loading: true,
        dataForm: {
          hostName: ''
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
      ...mapActions('hostManager', ['getList', 'getDeleteRow', 'getConnectTest']),
      getDataList() {
        this.loading = true
        const params = {
          pageIndex: this.pageIndex,
          pageSize: this.pageSize,
          hostName: this.dataForm.hostName
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
        this.addOrUpdateVisible = true
        this.$nextTick(() => {
          this.$refs.addOrUpdate.init(row)
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
          this.getDeleteRow({rowIds: rowIds}).then((res) => {
            // console.log(res.data)
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
      connectTestHandle(id) {
        const hostId = id !== undefined ? id || 0 : undefined
        this.getConnectTest({hostId: hostId}).then((res) => {
          // console.log(res.data)
          const result = res.data
          if (result !== undefined && result.code === 200) {
            this.$notify({
              title: '成功',
              showClose: true,
              message: '连接成功',
              type: 'success',
              duration: 3000
            })
          } else {
            this.$notify({
              title: '失败',
              showClose: true,
              message: '连接失败',
              type: 'error',
              duration: 3000
            })
          }
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
