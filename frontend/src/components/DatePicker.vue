<script setup>
  import { ref, watch, computed, onBeforeMount } from 'vue';
  import styles from '../style.module.css';

  // stored epoch is in seconds utc
  const props = defineProps(['dateValue']);
  const emit = defineEmits(['update']);

  // javascript epoch is in milliseconds
  const epoch = !!props.dateValue
    ? props.dateValue * 1000
    : Date.now().valueOf();

  const dateObj = new Date(epoch);

  // prefix single-digit date numbers
  const prefixed = (datenumber) => {
    const prefix = datenumber < 10 ? '0' : '';
    return `${prefix}${datenumber}`;
  };

  // date format is YYYY-MM-DD
  const date = ref(
    `${dateObj.getFullYear()}-${prefixed(dateObj.getMonth() + 1)}-${prefixed(
      dateObj.getDate()
    )}`
  );

  // time format is HH:mm
  const time = ref(`${dateObj.getHours()}:${dateObj.getMinutes()}`);

  const dateTime = ref(`${date.value} ${time.value}`);

  const newDateTime = computed(() => {
    return `${date.value} ${time.value}`;
  });

  // transforms the date string to timestamp UTC in seconds
  const stringToEpoch = (dateString) => {
    const date = new Date(dateString);
    const milliseconds = date.valueOf();
    return Math.floor(milliseconds / 1000);
  };

  // emit date if it was set to now
  if (!!!props.dateValue) {
    emit('update', stringToEpoch(dateTime.value));
  }

  const updateTime = () => {
    emit('update', stringToEpoch(newDateTime.value));
    dateTime.value = newDateTime.value;
  };
</script>

<template>
  <div :class="styles.datePicker">
    <q-input outlined dark :model-value="dateTime">
      <template v-slot:prepend>
        <q-icon name="event" class="cursor-pointer">
          <q-popup-proxy cover transition-show="scale" transition-hide="scale">
            <q-date v-model="date" mask="YYYY-MM-DD" dark>
              <div class="row items-center justify-end">
                <q-btn
                  v-close-popup
                  label="Close"
                  color="primary"
                  text-color="dark"
                  flat
                  @click="updateTime()"
                />
              </div>
            </q-date>
          </q-popup-proxy>
        </q-icon>
      </template>

      <template v-slot:append>
        <q-icon name="access_time" class="cursor-pointer">
          <q-popup-proxy cover transition-show="scale" transition-hide="scale">
            <q-time v-model="time" mask="HH:mm" format24h dark>
              <div class="row items-center justify-end">
                <q-btn
                  v-close-popup
                  label="Close"
                  color="primary"
                  text-color="dark"
                  flat
                  @click="updateTime()"
                />
              </div>
            </q-time>
          </q-popup-proxy>
        </q-icon>
      </template>
    </q-input>
  </div>
</template>
