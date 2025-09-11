<script setup>
  import WorkoutStatus from './WorkoutStatus.vue';
  import { inject, provide, ref, watch } from 'vue';
  import { programInstanceStore } from '../modules/state';
  import { QCarousel, QCarouselSlide } from 'quasar';
  import { updateProgramInstance } from './../modules/utils';
  import * as styles from '../style.module.css';
  import { useRouter } from 'vue-router';
  import * as utils from '../modules/programUtils';
  import CarouselNav from './CarouselNav.vue';

  const router = useRouter();

  const { focusedEvent, setFocusedEvent } = inject('focusedEvent');
  const { selectedEvent, setSelectedEvent } = inject('selectedEvent');

  const props = defineProps({
    activityID: String,
    workoutCoords: Array,
    dayIndex: Number,
  });

  const currentInstance = programInstanceStore.getCurrent(props.activityID);

  if (!currentInstance.events) {
    currentInstance['events'] = {};
  }

  // content for slides
  const slides = ref(utils.getWorkouts(props.activityID));

  // carousel model (slide index) is the workout day
  const slide = ref(props.dayIndex);

  provide('current', slide);

  setFocusedEvent(currentInstance.events[slide.value]);

  const statusColourStyles = [];
  for (let i = 0; i < slides.value.length; i++) {
    const status = utils.getWorkoutStatus(
      currentInstance.events[i],
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
        currentInstance.events[props.dayIndex] = '';
        updateProgramInstance(currentInstance);
      } else {
        router.push({
          name: 'event',
          query: {
            instance: currentInstance.id,
            block: props.workoutCoords[0],
            cycle: props.workoutCoords[1],
            workout: props.workoutCoords[2],
            day: props.dayIndex,
          },
        });
      }
    }
  };

  const skipWorkout = (outstandingIndex) => {
    currentInstance.events[
      props.dayIndex - (slides.value.length - 1 - outstandingIndex)
    ] = '';

    updateProgramInstance(currentInstance);
  };

  watch(
    () => selectedEvent.value,
    (newID) => {
      for (let i = 0; i < Object.entries(currentInstance.events).length; i++) {
        if (currentInstance.events[i] == newID) {
          slide.value = i;
          setFocusedEvent(newID);
          break;
        }
      }
    }
  );
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
            setFocusedEvent(currentInstance.events[value]);
          }
        "
      >
        <q-carousel-slide v-for="(w, index) of slides" :name="index">
          <div :class="[styles.carouselSlideContent]">
            <WorkoutStatus
              :class="[styles.workoutStatus]"
              :eventID="currentInstance.events[index]"
              :todayIndex="dayIndex"
              :workoutIndex="index"
              :workout="w"
            />
            <div
              :class="[styles.carouselNavButtonArray]"
              v-if="
                slide == props.dayIndex &&
                currentInstance.events[slide] == undefined &&
                !slides[slide].restDay
              "
            >
              <q-btn
                round
                icon="play_circle_filled"
                color="primary"
                @Click="
                  () => {
                    startWorkout(slide);
                  }
                "
                :disable="
                  slide != props.dayIndex ||
                  currentInstance.events[slide] != undefined ||
                  slides[slide].restDay
                "
              />
              <q-btn
                round
                color="primary"
                icon="do_not_disturb"
                :disable="
                  slide > props.dayIndex ||
                  currentInstance.events[slide] != undefined ||
                  slides[slide].restDay
                "
                @Click="
                  () => {
                    skipWorkout(slide);
                  }
                "
              />
            </div>
          </div>
        </q-carousel-slide>
      </q-carousel>
      <CarouselNav :items="statusColourStyles" :homeItem="props.dayIndex" />
    </div>
  </div>
</template>
