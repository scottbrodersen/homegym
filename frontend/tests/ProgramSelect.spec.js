import ProgramSelect from '../src/components/ProgramSelect.vue';
import { QDialog, QItem, QSelect } from 'quasar';
import { focus } from '../src/modules/directives';
import { vi, afterEach } from 'vitest';
import { installQuasarPlugin } from '@quasar/quasar-app-extension-testing-unit-vitest';
import {
  config,
  mount,
  enableAutoUnmount,
  flushPromises,
} from '@vue/test-utils';
import * as data from '../mocks/data';

config.global.directives = {
  focus: focus,
};

config.global.errorHandler = (err) => {
  throw err;
};

installQuasarPlugin({
  components: { QDialog, QItem, QSelect },
});

vi.mock('../src/modules/state');

describe('ProgramSelect component', () => {
  enableAutoUnmount(afterEach);

  afterEach(() => {
    vi.restoreAllMocks();
  });
  it('has expected initial state', async () => {
    const wrapper = mount(ProgramSelect, {
      components: {},
      attachTo: document.body,
    });
    expect(wrapper.getComponent(QSelect)).toBeDefined;
    expect(wrapper.getComponent(QSelect).isVisible()).toEqual(true);
  });
});
