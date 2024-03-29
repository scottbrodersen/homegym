import { vi } from 'vitest';
import {
  config,
  enableAutoUnmount,
  flushPromises,
  mount,
} from '@vue/test-utils';
import ActivityProgram from './../src/components/ActivityProgram.vue';
import { states } from './../src/modules/utils';
import ProgramBlock from '../src/components/ProgramBlock.vue';
import { QBtn, QInput } from 'quasar';
import * as data from '../mocks/data';
import { focus } from '../src/modules/directives';
import { nextTick } from 'vue';
import { installQuasarPlugin } from '@quasar/quasar-app-extension-testing-unit-vitest';

config.global.directives = {
  focus: focus,
};

config.global.errorHandler = (err) => {
  throw err;
};

enableAutoUnmount(afterEach);

installQuasarPlugin({
  components: { QBtn, QInput },
});

vi.mock('../src/modules/state');

describe('ActivityProgram component', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('renders read only', () => {
    const wrapper = mount(ActivityProgram, {
      components: { ProgramBlock },
      global: {
        provide: {
          state: { value: states.READ_ONLY },

          activity: { value: data.fetchedTestActivities[0] },
        },
      },
      props: {
        programID: 'anything',
      },
    });
  });

  it('form components read only', async () => {
    const wrapper = mount(ActivityProgram, {
      components: { ProgramBlock },
      global: {
        provide: {
          state: { value: states.READ_ONLY },
          activity: { value: data.fetchedTestActivities[0] },
        },
      },
      props: {
        programID: data.testProgram().id,
      },
    });
    nextTick();
    await flushPromises();
    wrapper.findAllComponents(QBtn).forEach((w) => {
      expect(w.isVisible()).toBeFalsy;
    });
    wrapper.findAllComponents(QInput).forEach((w) => {
      expect(w.isVisible()).toBeFalsy;
    });
    expect(wrapper.findAllComponents(ProgramBlock)).toHaveLength(
      data.testProgram().blocks.length
    );
  });

  it('done and save buttons in EDIT state', async () => {
    const wrapper = mount(ActivityProgram, {
      components: { ProgramBlock },
      global: {
        provide: {
          state: { value: states.EDIT },
          activity: { value: data.fetchedTestActivities[0] },
        },
      },
      props: {
        programID: 'anything',
      },
    });
    const buttons = wrapper.findAllComponents(QBtn);

    let doneWrapper;
    let updateWrapper;
    for (let i = 0; i < buttons.length; i++) {
      if (buttons[i].text() == 'Done') {
        doneWrapper = buttons[i];
      } else if (buttons[i].text() == 'Update') {
        updateWrapper = buttons[i];
      }
      if (doneWrapper && updateWrapper) break;
    }
    expect(doneWrapper).not.toBeUndefined();
    expect(doneWrapper.isVisible()).toBeTruthy();
    expect(doneWrapper.attributes().disabled).toBeFalsy();
    expect(updateWrapper).not.toBeUndefined();
    expect(updateWrapper.isVisible()).toBeTruthy();
    expect(updateWrapper.attributes().disabled).not.toBeUndefined();

    let foundBlockTitle = false;
    let p;
    wrapper.findAll('input').forEach((w) => {
      console.log(w.text());
      if (
        !foundBlockTitle &&
        w.element.value == data.testProgram().blocks[0].title
      ) {
        foundBlockTitle = true;
        p = w.setValue('changed').then(() => {
          expect(doneWrapper.text()).toEqual('Cancel');
          expect(updateWrapper.attributes().disabled).toBeFalsy();
        });
      }
    });
    expect(foundBlockTitle).toBeTruthy();
  });
});
