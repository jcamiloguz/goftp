<script setup>
import Card from './Card.vue'
import { reactive, watch } from 'vue'
import Error from './Error.vue'
import Spinner from './Spinner.vue'

const data = reactive({
  loading: false,
  error: null,
  channels: [],
  fileCount: '0',
  fileSize: '0',
  totalSubscribers: '0'
})

const socket = new WebSocket('ws://localhost:8080/socket')
socket.onopen = () => {
  console.log('connected')
  data.loading = true
}
socket.onmessage = (event) => {
  const payload = JSON.parse(event.data)
  console.log(payload)
  data.channels = payload.channels
  data.loading = false
}
socket.onerror = (event) => {
  console.log(event)
  // TODO manage error
  data.error = { title: 'Error', message: 'Server error' }
  data.loading = false
}
socket.onclose = () => {
  console.log('disconnected')
  data.loading = false
  data.error = {
    title: 'Disconnected',
    message: 'Please refresh the page'
  }
}

const formatPayload = (rawPayload) => {
  return new Promise((resolve, reject) => {
    const formatPayload = {}
    formatPayload.fileCount = rawPayload.reduce(
      (acc, channel) => acc + channel.files.length,
      0
    )
    const totalSize =
      rawPayload.reduce(
        (acc, channel) =>
          acc + channel.files.reduce((acc, file) => acc + file.size, 0),
        0
      ) + ''
    formatPayload.totalSubscribers =
      rawPayload.reduce((acc, channel) => acc + channel.subscribers.length, 0) +
      ''
    const totalSizeKB = totalSize / 1024
    const totalSizeMB = totalSizeKB / 1024
    const totalSizeGB = totalSizeMB / 1024
    if (totalSizeGB > 1) {
      formatPayload.fileSize = `${totalSizeGB.toFixed(2)} MB`
      return
    }
    if (totalSizeMB > 1) {
      formatPayload.fileSize = `${totalSizeMB.toFixed(2)} KB`
      return
    }
    if (totalSizeKB > 1) {
      formatPayload.fileSize = `${totalSizeKB.toFixed(2)} KB`
      return
    }
    formatPayload.fileSize = `${totalSize} bytes`
    resolve(formatPayload)
  })
}

watch(
  () => data.channels,
  async (newChannels) => {
    try {
      const { fileSize, totalSubscribers, fileCount } = await formatPayload(
        newChannels
      )
      data.fileCount = fileCount
      data.fileSize = fileSize
      data.totalSubscribers = totalSubscribers
    } catch (e) {
      console.log(e)
      data.error = {
        title: 'Error',
        message: 'Data integrity error'
      }
    }
  }
)
</script>

<template>
  <div class="dashboard__container">
    <div class="dashboard__header">
      <Card title="Files Transfered" :value="data.fileCount" :color="'blue'" />
      <Card title="Total Size Transfered" :value="data.fileSize" color="red" />
      <Card title="Subscribers" :value="data.totalSubscribers" color="green" />
      <Spinner v-show="data.loading" />
      <Error
        v-show="data.error"
        :title="data.error?.title"
        :message="data.error?.message"
      />
    </div>
  </div>
</template>

<style scoped>
.dashboard__container {
  max-width: 1250px;
  margin: 0 auto;
}
.dashboard__header {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
}
</style>
