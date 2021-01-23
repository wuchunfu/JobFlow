<template>
  <el-row class="row-box">
    <el-col class="col-box">
      <el-card class="card-box">
        <el-form :inline="true" :model="dataForm" @keyup.enter.native="getDataList()">
          <el-form-item>
            <el-input v-model="dataForm.username" placeholder="用户名" clearable/>
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
            label="用户名"
            header-align="center"
            align="center"
            fixed
          >
            <template slot-scope="scope">
              {{ scope.row.username }}
            </template>
          </el-table-column>
          <el-table-column
            label="邮箱"
            header-align="center"
            align="center"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.email }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="角色"
            header-align="center"
            align="center"
          >
            <template slot-scope="scope">
              {{ scope.row.isAdmin | roleConvert }}
            </template>
          </el-table-column>
          <el-table-column
            label="状态"
            class-name="status-col"
            header-align="center"
            align="center"
          >
            <template slot-scope="scope">
              <el-tag :type="scope.row.status | statusFilter">
                {{ scope.row.status | statusConvert }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            label="创建时间"
            header-align="center"
            align="center"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.createTime }}</span>
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
                @click="changePasswordHandle(scope.row)"
              >
                修改密码
              </el-button>
              <el-button
                v-loading="loading"
                type="danger"
                size="small"
                @click="deleteHandle(scope.row.userId)"
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
        <!-- 弹窗, 修改密码 -->
        <change-password
          v-if="changePasswordVisible"
          ref="changePassword"
          @refreshDataList="getDataList"
        />
      </el-card>
    </el-col>
  </el-row>
</template>

<script>
  import {mapActions} from 'vuex'
  import AddOrUpdate from '../addOrUpdate/index'
  import ChangePassword from '../changePassword/index'

  export default {
    filters: {
      statusFilter(status) {
        const statusMap = {
          0: 'danger',
          1: 'success'
        }
        return statusMap[status]
      },
      statusConvert(status) {
        const statusMap = {
          0: '禁用',
          1: '正常'
        }
        return statusMap[status]
      },
      roleConvert(status) {
        const statusMap = {
          0: '普通用户',
          1: '管理员'
        }
        return statusMap[status]
      }
    },
    components: {
      AddOrUpdate,
      ChangePassword
    },
    data() {
      return {
        loading: true,
        dataForm: {
          username: ''
        },
        dataListSelections: [],
        dataList: [],
        addOrUpdateVisible: false,
        changePasswordVisible: false,
        pageIndex: 1,
        pageSize: 10,
        totalPage: 0
      }
    },
    created() {
      this.getDataList()
    },
    methods: {
      ...mapActions('userManager', ['getList', 'getDeleteRow']),
      getDataList() {
        this.loading = true
        const params = {
          pageIndex: this.pageIndex,
          pageSize: this.pageSize,
          username: this.dataForm.username
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
      // 修改密码
      changePasswordHandle(row) {
        console.log(row)
        this.changePasswordVisible = true
        this.$nextTick(() => {
          this.$refs.changePassword.init(row)
        })
      },
      //  删除
      deleteHandle(id) {
        const rowIds = id ? [id] : this.dataListSelections.map(item => {
          return item.userId
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
