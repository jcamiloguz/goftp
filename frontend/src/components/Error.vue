<script setup>
const props = defineProps({
  title: String,
  message: String
})
</script>
<template>
  <div class="modal--error">
    <div class="modal--error__content">
      <img class="modal--error__img" src="../assets/error.png" alt="" />
      <h1 class="modal--error__title">
        {{ props.title }}
      </h1>
      <p class="modal--error__message">
        {{ props.message }}
      </p>
    </div>
  </div>
</template>

<style>
.modal--error {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  z-index: 3;
  display: flex;
  justify-content: center;
  align-items: center;
}
.modal--error__content {
  background-color: var(--blank-color);
  padding: 20px;
  border-radius: 20px;
  width: 400px;
  margin: auto;
}
.modal--error__title {
  font-size: 2.25rem;
  text-align: center;
  font-weight: bold;
}
.modal--error__message {
  font-size: 1.5rem;
  text-align: center;
  font-weight: bold;
}
.modal--error__img {
  width: 100%;
  margin: 20px auto;
}
</style>
