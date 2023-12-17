<script setup>
  import styles from '../style.module.css';
  import { QInput } from 'quasar';
  import { ref } from 'vue';

  const props = defineProps({ writable: Boolean, time: Number });
  const emit = defineEmits(['update']);

  // We store the value as a number (in seconds) but present it as a string (h:mm:ss)
  const time = ref('');
  const hours = ref('');
  const minutes = ref('');
  const seconds = ref('');

  if (!!props.time) {
    hours.value = Math.floor(props.time / 3600).toString();

    minutes.value = Math.floor((props.time % 3600) / 60).toString();
    if (minutes.value.length == 1) {
      minutes.value = `0${minutes.value}`;
    }

    seconds.value = Math.floor((props.time % 360) % 60).toString();
    if (seconds.value.length == 1) {
      seconds.value = `0${seconds.value}`;
    }
    time.value = `${hours.value}:${minutes.value}:${seconds.value}`;
  }

  const validate = (timeValue) => {
    const regex = new RegExp('^[0-9]:[0-5][0-9]:[0-5][0-9]$');
    return regex.test(timeValue);
  };

  const update = (newTime) => {
    if (validate(newTime)) {
      const timeParts = newTime.split(':');

      const hours = Number(timeParts[0]);
      const mins = Number(timeParts[1]);
      const secs = Number(timeParts[2]);

      emit('update', hours * 3600 + mins * 60 + secs);
    }
  };
</script>
<template>
  <div v-if="props.writable" :class="[styles.horiz]">
    <q-input
      v-model="time"
      :class="[styles.timeInput]"
      focus
      select
      dark
      mask="#:##:##"
      label="h:mm:ss"
      stack-label
      :rules="[validate]"
      lazy-rules
      @update:model-value="update"
    />
  </div>
  <div v-else>{{ time }}</div>
</template>
