<script setup>
import LogoGoftp from './icons/LogoGoftp.vue'
</script>

<template>
  <header class="header">
    <div class="header__container">
      <div class="header__logo">
        <LogoGoftp />
        <div class="header__logo-divider"></div>
        <h1 class="header__title">Dashboard</h1>
      </div>
      <nav class="header__nav">
        <ul class="header__nav-list">
          <li class="header__nav-item">Overview</li>
        </ul>
      </nav>
    </div>
  </header>
</template>

<style scoped>
.header {
  height: 80px;
  background: white;
  box-shadow: 0px 4px 20px rgba(0, 0, 0, 0.25);
}
.header__container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1250px;
  margin: 0 auto;
}
.header__logo {
  display: flex;
  height: 100%;
  align-items: center;
  padding: 16px 0;
}
.header__title {
  color: var(--gray-color);
  font-size: 2.1rem;
  font-weight: bold;
}
.header__logo-divider {
  width: 2px;
  margin: 0 15px;
  height: 50px;
  background: var(--gray-color);
}

.header__nav {
  display: flex;
  justify-content: flex-end;
  width: 33%;
  height: 100%;
}
.header__nav-list {
  display: flex;
  list-style: none;
}
.header__nav-item {
  padding: 0 15px;
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--dark-color);
  cursor: pointer;
  border-bottom: 3px solid transparent;
  /* transition: 0.45s ease; */
}
.header__nav-item:hover {
  border-bottom: 3px solid var(--go-blue);
}
</style>
