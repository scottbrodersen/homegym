import ActivityPage from '../src/components/ActivityPage.vue';
import { config, mount } from '@vue/test-utils';
import { Quasar } from 'quasar';
import { focus } from '../src/modules/directives';
import { states } from '../src/modules/utils';
config.global.plugins.push(Quasar);
config.global.directives = {
  focus: focus,
};

console.log(config);
describe('Program module', () => {
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
