import ActivityPage from '../src/components/ActivityPage.vue';
import { config, enableAutoUnmount, mount } from '@vue/test-utils';
import { focus } from '../src/modules/directives';
import { states } from '../src/modules/utils';
import { installQuasarPlugin } from '@quasar/quasar-app-extension-testing-unit-vitest';

config.global.errorHandler = (err) => {
  throw err;
};

installQuasarPlugin();

config.global.directives = {
  focus: focus,
};

enableAutoUnmount(afterEach);

describe('ActivityPage component', () => {
  it('renders correctly read only', () => {
    const wrapper = mount(ActivityPage, {
      provide: {
        state() {
          return states.READ_ONLY;
        },
      },
    });
  });
});
