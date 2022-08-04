<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: String,
  value: String,
  color: String
})
const cardStyle = computed(() => ({
  backgroundColor: `var(--go-${props.color})`,
  color: props.color ? 'white' : 'var(--dark-color)'
}))
</script>

<template>
  <div class="card" :style="cardStyle">
    <h2 class="card__title">
      {{ props.title }}
    </h2>
    <p class="card__value">{{ props.value }}</p>
  </div>
</template>

<style scoped>
.card {
  width: 100%;
  height: 100%;
  max-width: 350px;
  min-width: 164px;
  margin: 16px auto;
  padding: 20px 0;
  color: white;
  box-shadow: 0px 4px 20px rgba(0, 0, 0, 0.25);
  border-radius: 20px;
}
.card__title {
  font-size: 2.25rem;
  text-align: center;
  font-weight: bold;
}
.card__value {
  font-size: 4rem;
  text-align: center;
  font-weight: bold;
}
</style>
