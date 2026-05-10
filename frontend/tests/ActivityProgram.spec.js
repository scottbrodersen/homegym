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
import { nextTick, ref } from 'vue';
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
          editProgramTitle: {
            editProgramTitle: ref(false),
            toggleProgramTitle: () => {},
          },
          newProgram: { newProgram: ref(false), toggleNewProgram: () => {} },
          cloneProgram: {
            cloneProgram: ref(false),
            toggleCloneProgram: () => {},
          },
          //     activity: { value: data.fetchedTestActivities[0] },
        },
      },
      props: {
        programID: 'anything',
        activityID: data.fetchedTestActivities[0].id,
      },
    });
  });

  it('form components read only', async () => {
    const wrapper = mount(ActivityProgram, {
      components: { ProgramBlock },
      global: {
        provide: {
          state: { value: states.READ_ONLY },
          editProgramTitle: {
            editProgramTitle: ref(false),
            toggleProgramTitle: () => {},
          },
          newProgram: { newProgram: ref(false), toggleNewProgram: () => {} },
          cloneProgram: {
            cloneProgram: ref(false),
            toggleCloneProgram: () => {},
          },
        },
      },
      props: {
        programID: data.testProgram().id,
        activityID: data.fetchedTestActivities[0].id,
      },
    });
    nextTick();
    await flushPromises();
    console.log(wrapper.html());
    wrapper.findAllComponents(QBtn).forEach((w) => {
      expect(w.isVisible()).toBeFalsy;
    });
    wrapper.findAllComponents(QInput).forEach((w) => {
      expect(w.isVisible()).toBeFalsy;
    });
    expect(wrapper.findAllComponents(ProgramBlock)).toHaveLength(1);
  });

  it('edit button in EDIT state', async () => {
    const wrapper = mount(ActivityProgram, {
      components: { ProgramBlock },
      global: {
        provide: {
          state: { value: states.EDIT },
          //          activity: { value: data.fetchedTestActivities[0] },
          editProgramTitle: {
            editProgramTitle: ref(false),
            toggleProgramTitle: () => {},
          },
          newProgram: { newProgram: ref(false), toggleNewProgram: () => {} },
          cloneProgram: {
            cloneProgram: ref(false),
            toggleCloneProgram: () => {},
          },
        },
      },
      props: {
        programID: 'anything',
        activityID: data.fetchedTestActivities[0].id,
      },
    });
    await nextTick();
    await flushPromises();
    const buttons = wrapper.findAllComponents(QBtn);
    const editButton = buttons.find((btn) => btn.props().icon === 'edit_note');
    expect(editButton).toBeDefined();
    expect(editButton.isVisible()).toBeTruthy();
  });
});
