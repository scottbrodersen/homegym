<script setup>
  /**
   * Indicates that a rep was either failed or successful.
   * Props:
   *  success: true for a successful rep, otherwise false.
   */
  import * as styles from '../style.module.css';
  const props = defineProps({ success: Boolean });
  const repClass = props.success ? 'hgRepGood' : 'hgRepFail';
</script>
<template>
  <span :class="[styles[repClass], 'material-icons']">circle</span>
</template>
