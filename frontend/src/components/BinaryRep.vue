<script setup>
  import * as styles from '../style.module.css';
  const props = defineProps({ success: Boolean });
  const repClass = props.success ? 'hgRepGood' : 'hgRepFail';
</script>
<template>
  <span :class="[styles[repClass], 'material-icons']">circle</span>
</template>
