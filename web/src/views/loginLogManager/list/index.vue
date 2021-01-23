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
            label="登陆IP"
            header-align="center"
            align="center"
            fixed
          >
            <template slot-scope="scope">
              {{ scope.row.ip }}
            </template>
          </el-table-column>
          <el-table-column
            label="登陆时间"
            header-align="center"
            align="center"
            fixed
          >
            <template slot-scope="scope">
              {{ scope.row.createTime }}
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
    components: {},
    data() {
      return {
        loading: true,
        dataForm: {
          username: ''
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
      ...mapActions('loginLogManager', ['getList']),
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
