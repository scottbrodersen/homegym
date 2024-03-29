import ExerciseIntensity from '../src/components/ExerciseIntensity.vue';
import { config, enableAutoUnmount, mount } from '@vue/test-utils';
import { QInput } from 'quasar';
import { focus, select } from '../src/modules/directives';
import { intensityProps } from '../src/modules/utils';
import { installQuasarPlugin } from '@quasar/quasar-app-extension-testing-unit-vitest';

installQuasarPlugin({
  components: { QInput },
});

config.global.directives = {
  focus: focus,
  select: select,
};
config.global.errorHandler = (err) => {
  throw err;
};

enableAutoUnmount(afterEach);

describe('ExerciseIntensity component', () => {
  it('weight value renders correctly read only', () => {
    const testWeight = 100;
    const testType = 'weight';
    const wrapper = mount(ExerciseIntensity, {
      props: {
        intensity: testWeight,
        writable: false,
        type: testType,
      },
    });
    let found = false;
    const elements = wrapper.findAll('div');
    for (let i = 0; i < elements.length; i++) {
      if (elements[i].text == intensityProps(testType).format(testWeight)) {
        found = true;
        break;
      }
    }
    expect(found).toBeTruthy;
  });

  it('bodyweight type renders correctly read only', () => {
    const testWeight = 1;
    const testType = 'bodyweight';
    const wrapper = mount(ExerciseIntensity, {
      props: {
        intensity: testWeight,
        writable: false,
        type: testType,
      },
    });
    let found = false;
    const elements = wrapper.findAll('div');
    for (let i = 0; i < elements.length; i++) {
      if (elements[i].text == intensityProps(testType).format(testWeight)) {
        found = true;
        break;
      }
    }
    expect(found).toBeTruthy;
  });
});

it('weight value renders correctly when editing', async () => {
  const testWeight = 100;
  const testType = 'weight';
  const wrapper = mount(ExerciseIntensity, {
    attachTo: document.body,
    props: {
      intensity: testWeight,
      writable: true,
      type: testType,
    },
  });
  const input = wrapper.find('input');
  expect(input.exists()).toBe(true);

  await input.setValue(testWeight + 1);
  expect(wrapper.emitted('update')).toBeDefined;
  expect(wrapper.emitted('update')).toHaveLength(1);
  expect(wrapper.emitted('update')[0]).toEqual([
    intensityProps(testType).value(testWeight + 1),
  ]);
});
