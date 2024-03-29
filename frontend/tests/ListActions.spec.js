import { config, enableAutoUnmount, mount } from '@vue/test-utils';
import ListActions from './../src/components/ListActions.vue';
import { OrderedList } from '../src/modules/utils';
import { QBtn } from 'quasar';
import { installQuasarPlugin } from '@quasar/quasar-app-extension-testing-unit-vitest';

config.global.errorHandler = (err) => {
  throw err;
};

installQuasarPlugin({ components: { QBtn } });

enableAutoUnmount(afterEach);

describe('ListActions', () => {
  it('The correct number of buttons are rendered', () => {
    const wrapper = mount(ListActions);
    expect(wrapper.findAll('button')).toHaveLength(4);
  });
  it('The correct events are emitted', () => {
    const wrapper = mount(ListActions);

    wrapper.findAll('button').forEach((btn) => {
      btn.trigger('click');
    });

    const emittedEvents = wrapper.emitted('update');
    expect(emittedEvents).toHaveLength(4);
    const args = [];
    emittedEvents.forEach((evt) => {
      args.push(evt[0]);
    });

    const expectedArgs = [
      OrderedList.ADD,
      OrderedList.DELETE,
      OrderedList.MOVEBACK,
      OrderedList.MOVEFWD,
    ];
    expect(args.sort()).toEqual(expectedArgs.sort());
  });
});
