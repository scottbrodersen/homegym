<script setup>
  import { inject, onMounted, watch } from 'vue';
  import { QBtn } from 'quasar';
  import * as styles from '../style.module.css';
  /**
   * Indexed icons that represent items in a carousel.
   * The current item (i.e. the selected item) is injected.
   * Distinguishes a 'home' item by its index.
   *
   * Props:
   *  items: an array of css class names. Each item represents a navigation target. The css class name controls the appearance.
   *  homeItem: the index of the item to show as the home item.
   */
  const props = defineProps({ items: Array, homeItem: Number });
  // items is an array of css class names

  const current = inject('current');

  const scroll = (currentSlide) => {
    const cur = document.getElementById('nav' + currentSlide);
    cur.scrollIntoView({
      behavior: 'smooth',
      block: 'center',
      inline: 'center',
    });
  };

  watch(
    () => {
      return current.value;
    },
    (newCurrent) => {
      scroll(newCurrent);
    }
  );

  onMounted(() => {
    scroll(current.value);
  });
</script>
<template>
  <div :class="[styles.carouselNav]">
    <div :class="[styles.carouselNavWindow]">
      <div :class="[styles.carouselNavBtnWrap]">
        <div
          v-for="(item, index) in props.items"
          :key="index"
          :id="'nav' + index"
        >
          <q-btn
            :round="index != props.homeItem"
            :icon="index == props.homeItem ? 'today' : ''"
            :class="[
              index == current
                ? styles.carouselNavCurrent
                : styles.carouselNavBtn,
              styles[item],
              index == props.homeItem ? styles.carouselNavToday : '',
            ]"
          />
        </div>
      </div>
    </div>
  </div>
</template>
