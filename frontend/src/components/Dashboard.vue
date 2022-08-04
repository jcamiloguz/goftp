<script setup>
import Card from './Card.vue'
</script>

<template>
  <div class="dashboard__container">
    <div class="dashboard__header">
      <Card title="Example" value="213" color="blue" />
      <Card title="Example" value="213" color="red" />
      <Card title="Example" value="213" color="green" />
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
