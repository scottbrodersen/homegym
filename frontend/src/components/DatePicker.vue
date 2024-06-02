<script setup>
  import { ref, computed } from 'vue';
  import * as styles from '../style.module.css';
  import { QTime, QInput, QPopupProxy, QDate, QBtn } from 'quasar';
  import * as dateUtils from '../modules/dateUtils';

  // expect dateValue to be seconds since epoch (utc)
  const props = defineProps({ dateValue: Number, hideTime: Boolean });
  const emit = defineEmits(['update']);

  // Sets to today if no props.dateValue
  const dateObj = dateUtils.dateFromSeconds(props.dateValue);

  const date = ref(dateUtils.formatDate(dateObj));

  const time = props.hideTime ? ref('') : ref(dateUtils.formatTime(dateObj));

  const dateTime = ref(dateUtils.formatDateTime(dateObj));

  const newDateTime = computed(() => {
    return `${date.value} ${time.value}`;
  });

  // emit date if it was set to now
  if (!props.dateValue) {
    emit('update', dateUtils.stringToEpoch(dateTime.value));
  }

  const updateTime = () => {
    emit('update', dateUtils.stringToEpoch(newDateTime.value));
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
                  label="OK"
                  color="primary"
                  dark
                  flat
                  @click="updateTime()"
                />
              </div>
            </q-date>
          </q-popup-proxy>
        </q-icon>
      </template>

      <template v-slot:append>
        <q-icon
          v-show="!props.hideTime"
          name="access_time"
          class="cursor-pointer"
        >
          <q-popup-proxy cover transition-show="scale" transition-hide="scale">
            <q-time v-model="time" mask="HH:mm" format24h dark>
              <div class="row items-center justify-end">
                <q-btn
                  v-close-popup
                  label="OK"
                  color="primary"
                  flat
                  dark
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
