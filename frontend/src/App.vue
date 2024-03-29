<script setup>
  import InitState from './components/InitState.vue';
</script>

<template>
  <Suspense>
    <template #default>
      <InitState />
    </template>
    <template #fallback> Product is loading... </template>
  </Suspense>
</template>
