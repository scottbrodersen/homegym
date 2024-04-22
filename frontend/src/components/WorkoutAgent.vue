<script setup>
  import WorkoutStatus from './WorkoutStatus.vue';
  import { provide, ref } from 'vue';
  import { programInstanceStore } from '../modules/state';
  import { QCarousel, QCarouselSlide } from 'quasar';
  import { updateProgramInstance } from './../modules/utils';
  import styles from '../style.module.css';
  import { useRouter } from 'vue-router';
  import * as utils from '../modules/programUtils';
  import CarouselNav from './CarouselNav.vue';

  const router = useRouter();

  const props = defineProps({
    activityID: String,
    workoutCoords: Array,
    dayIndex: Number,
  });

  const activeInstance = programInstanceStore.getActive(props.activityID);

  if (!activeInstance.events) {
    activeInstance['events'] = {};
  }

  const slides = ref();

  slides.value = utils.getWorkouts(props.activityID);

  const slide = ref(props.dayIndex);

  const currentIndex = ref(props.dayIndex);
  provide('current', currentIndex);

  const statusColourStyles = [];
  for (let i = 0; i < slides.value.length; i++) {
    const status = utils.getWorkoutStatus(
      activeInstance.events[i],
      i,
      props.dayIndex,
      slides.value[i].restDay
    );
    statusColourStyles.push(utils.getColorStyle(status));
  }

  const startWorkout = (slideIndex) => {
    // ensure we can only start today's workout
    if (props.dayIndex == slideIndex) {
      if (slides.value[slideIndex].restDay) {
        activeInstance.events[props.dayIndex] = '';
        updateProgramInstance(activeInstance);
      } else {
        router.push({
          name: 'event',
          query: {
            instance: activeInstance.id,
            block: props.workoutCoords[0],
            cycle: props.workoutCoords[1],
            workout: props.workoutCoords[2],
            day: props.dayIndex,
          },
        });
      }
    }
  };

  const goToProgram = () => {
    router.push({
      name: 'programs',
      query: {
        activity: props.activityID,
        instance: activeInstance.id,
      },
    });
  };

  const skipWorkout = (outstandingIndex) => {
    activeInstance.events[
      props.dayIndex - (slides.value.length - 1 - outstandingIndex)
    ] = '';

    updateProgramInstance(activeInstance);
  };

  const setEvent = (carouselIndex) => {
    const currentTableRowEl = document.getElementById(
      activeInstance.events[currentIndex.value]
    );
    if (
      currentTableRowEl &&
      currentTableRowEl.classList.contains('evtHighlight')
    ) {
      currentTableRowEl.classList.remove('evtHighlight');
    }
    currentIndex.value = carouselIndex;
    const tableRowEl = document.getElementById(
      activeInstance.events[carouselIndex]
    );
    if (tableRowEl) {
      tableRowEl.classList.add('evtHighlight');
    }
  };

  const buttonSize = '10px';
</script>
<template>
  <div>
    <div :class="[styles.carouselWrap]">
      <q-carousel
        v-model="slide"
        arrows
        dark
        swipeable
        animated
        transition-prev="slide-right"
        transition-next="slide-left"
        padding
        height="100%"
        @update:model-value="
          (value) => {
            setEvent(value);
          }
        "
      >
        <q-carousel-slide v-for="(w, index) of slides" :name="index">
          <WorkoutStatus
            :class="[styles.workoutStatus]"
            :eventID="activeInstance.events[index]"
            :todayIndex="dayIndex"
            :workoutIndex="index"
            :workout="w"
          />
        </q-carousel-slide>
      </q-carousel>
      <CarouselNav :items="statusColourStyles" :homeItem="props.dayIndex" />
    </div>
    <div :class="[styles.carouselNavButtonArray]">
      <q-btn
        round
        :size="buttonSize"
        icon="play_circle_filled"
        color="primary"
        @Click="
          () => {
            startWorkout(currentIndex);
          }
        "
        :disable="
          currentIndex != props.dayIndex ||
          activeInstance.events[currentIndex] != undefined ||
          slides[currentIndex].restDay
        "
      />
      <q-btn
        round
        :size="buttonSize"
        color="primary"
        icon="do_not_disturb"
        :disable="
          currentIndex >= props.dayIndex ||
          activeInstance.events[currentIndex] != undefined ||
          slides[currentIndex].restDay
        "
        @Click="
          () => {
            skipWorkout(currentIndex);
          }
        "
      />
      <q-btn
        round
        :size="buttonSize"
        color="primary"
        icon="visibility"
        @Click="goToProgram"
      />
    </div>
  </div>
</template>
