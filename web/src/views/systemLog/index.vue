<template>
  <el-row class="row-box">
    <el-col class="col-box">
      <el-card class="card-box">
        <div class="console-box">{{ log }}</div>
      </el-card>
    </el-col>
  </el-row>
</template>

<script>

  export default {
    name: 'Log',
    data() {
      return {
        log: ''
      }
    },
    created() {
      this.getDataList()
    },
    methods: {
      getDataList() {
        // 初始化 websocket
        const rootPath = process.env.VUE_APP_BASE_API
        const temp = rootPath.split('//')
        // 不指定 rootPath（为空），就使用当前 hostname
        const host = temp.length > 1 ? temp[temp.length - 1] : window.location.hostname
        const ws = new WebSocket('ws://' + host + '/websocket')
        // const ws = new WebSocket('ws://localhost:1060/websocket')

        // websocket 的响应函数
        ws.onopen = (evt) => {
          // console.log('Connection open ...')
          // console.log(evt)
          ws.send('Hello WebSockets!')
        }

        ws.onmessage = (evt) => {
          // console.log('Received Message: ' + evt.data)
          const oldLog = this.log
          const deltaLog = evt.data + '\n'
          let newLog = oldLog + deltaLog
          // 日志页的最大显示长度
          if (newLog.length > 10240) {
            // 如果超出，从头剔除超出部分
            newLog = newLog.slice(newLog.length - 10240, newLog.length)
          }
          // 写入到 web console 中
          this.log = newLog
        }

        ws.onclose = (evt) => {
          console.log('Connection closed.')
          console.log(evt)
        }
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
      padding: 24px;

      .el-pagination {
        margin-top: 20px;
        text-align: right;
      }
    }
  }

  .console-box {
    color: white;
    background: black;
    min-height: 360px;
    height: 720px;
    padding: 24px;
    margin: -24px;
    white-space: pre-wrap;
    border-radius: 8px;
    overflow: auto;
    font-size: 12px;
  }
</style>
