import ProgramsPage from './../src/components/ProgramsPage.vue';
import ActivityProgram from './../src/components/ActivityProgram.vue';
import ProgramSelect from '../src/components/ProgramSelect.vue';
import { Dialog, QBtn, QDialog, QSelect } from 'quasar';
import { focus } from '../src/modules/directives';
import { vi, afterEach } from 'vitest';
import { installQuasarPlugin } from '@quasar/quasar-app-extension-testing-unit-vitest';
import { config, mount, enableAutoUnmount } from '@vue/test-utils';
import * as data from '../mocks/data';

config.global.directives = {
  focus: focus,
};

config.global.errorHandler = (err) => {
  throw err;
};

installQuasarPlugin({
  components: { QBtn, QDialog, QSelect },
  plugins: { Dialog },
});

vi.mock('../src/modules/state');

enableAutoUnmount(afterEach);

describe('ProgramsPage component', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });
  it('has expected initial state', async () => {
    const wrapper = mount(ProgramsPage, {
      components: { ActivityProgram, ProgramSelect },
      attachTo: document.body,
    });

    // activity selector
    expect(wrapper.get('#activity')).toBeDefined;
    expect(wrapper.get('#activity').isVisible()).toEqual(true);
    expect(wrapper.get('#activity').isDisabled()).toEqual(false);

    // program selector
    expect(wrapper.getComponent(ProgramSelect)).toBeDefined;
    expect(wrapper.getComponent(ProgramSelect).isVisible()).toEqual(false);

    // new program button
    expect(wrapper.get('#new')).toBeDefined;
    expect(wrapper.get('#new').isVisible()).toEqual(true);
    expect(wrapper.get('#new').isDisabled()).toEqual(true);

    // edit program button
    expect(wrapper.find('#edit')).toBeDefined;
    expect(wrapper.find('#edit').isDisabled()).toEqual(true);
    expect(wrapper.find('#edit').isVisible()).toEqual(false);

    // start program button
    expect(wrapper.get('#start')).toBeDefined;
    expect(wrapper.get('#start').isVisible()).toEqual(false);

    expect(wrapper.getComponent(ActivityProgram)).toBeDefined;
    expect(wrapper.getComponent(ActivityProgram).props()).toEqual({
      programID: '',
    });
  });

  it('shows the program selector when an activity is selected', async () => {
    const wrapper = mount(ProgramsPage, {
      components: { QSelect, QBtn, QDialog, ActivityProgram },
      attachTo: document.body,
      props: { activity: data.fetchedTestActivities[0] },
    });

    expect(wrapper.getComponent(ProgramSelect).isVisible()).toEqual(true);
    expect(wrapper.find('#edit').isVisible()).toEqual(true);
    expect(wrapper.find('#edit').isDisabled()).toEqual(true);
    expect(wrapper.find('#new').isVisible()).toEqual(true);
    expect(wrapper.find('#new').isDisabled()).toEqual(false);
    expect(wrapper.find('#start').isVisible()).toEqual(false);
  });

  it('shows the program when a program is selected', async () => {
    const wrapper = mount(ProgramsPage, {
      components: { QSelect, QBtn, QDialog, ActivityProgram },
      attachTo: document.body,
      props: {
        activity: data.fetchedTestActivities[0],
        programID: data.testProgramID,
      },
    });

    expect(wrapper.find('#edit').isDisabled()).toEqual(false);
    expect(wrapper.find('#start').isVisible()).toEqual(true);
    expect(wrapper.getComponent(ActivityProgram).props().programID).toEqual(
      data.testProgramID
    );
  });
});
