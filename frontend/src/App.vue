<script setup>
  import styles from './style.module.css';
</script>

<template>
  <div :id="styles.wrapper">
    <router-view />
    <main :class="[styles.blockPadSm]">
      <router-view name="main" />
    </main>
  </div>
</template>

<style scoped></style>
