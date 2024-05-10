<script setup>
  import InitState from './components/InitState.vue';
  import { Suspense } from 'vue';
</script>

<template>
  <Suspense>
    <template #default>
      <Suspense>
        <InitState />
      </Suspense>
    </template>
    <template #fallback> Product is loading... </template>
  </Suspense>
</template>
