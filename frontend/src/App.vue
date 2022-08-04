<script setup>
import './assets/main.css'
import Header from './components/Header.vue'
import Dashboard from './components/Dashboard.vue'
import Spinner from './components/Spinner.vue'
import Error from './components/Error.vue'
import { ref } from 'vue'

const loading = ref(false)
const error = ref(null)

const socket = new WebSocket('ws://localhost:8080/socket')
socket.onopen = () => {
  console.log('connected')
  loading.value = true
}
socket.onmessage = (event) => {
  console.log(event.data)
  loading.value = false
}
socket.onclose = () => {
  console.log('disconnected')
  loading.value = false
  error.value = {
    title: 'Disconnected',
    message: 'Please refresh the page'
  }
}
</script>

<template>
  <Header />
  <main>
    <Dashboard />
    <Spinner v-if="loading" />
    <Error v-if="error" title="Error" message="Something went wrong" />
  </main>
</template>

<style scoped></style>
