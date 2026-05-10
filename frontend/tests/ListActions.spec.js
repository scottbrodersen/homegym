import {
  config,
  enableAutoUnmount,
  mount,
  flushPromises,
} from '@vue/test-utils';
import ListActions from './../src/components/ListActions.vue';
import { OrderedList } from '../src/modules/utils';
import { QBtn, QMenu } from 'quasar';
import { installQuasarPlugin } from '@quasar/quasar-app-extension-testing-unit-vitest';

config.global.errorHandler = (err) => {
  throw err;
};

installQuasarPlugin({ components: { QBtn } });

enableAutoUnmount(afterEach);

describe('ListActions', () => {
  it('The correct number of buttons are rendered', () => {
    const wrapper = mount(ListActions);
    expect(wrapper.findAll('button')).toHaveLength(1);
  });

  it('The menu items are rendered', async () => {
    const wrapper = mount(ListActions);

    wrapper.get('button').trigger('click');
    await flushPromises();

    // QMenu is teleported to the body, so search the document
    const menuItems = document.querySelectorAll('div.q-item');
    expect(menuItems).toHaveLength(6);
  });

  it('The correct events are emitted', async () => {
    const wrapper = mount(ListActions);

    wrapper.get('button').trigger('click');
    await flushPromises();

    // Click the menu items in the teleported menu
    const menuItems = document.querySelectorAll('div.q-item');
    menuItems.forEach((item) => {
      item.click();
    });
    await flushPromises();

    const emittedEvents = wrapper.emitted('update');
    expect(emittedEvents).toHaveLength(6);
    const args = [];
    emittedEvents.forEach((evt) => {
      args.push(evt[0]);
    });

    const expectedArgs = [
      OrderedList.ADD,
      OrderedList.DELETE,
      OrderedList.MOVEBACK,
      OrderedList.MOVEFWD,
      OrderedList.COPY,
      OrderedList.PASTE,
    ];
    expect(args.sort()).toEqual(expectedArgs.sort());
  });
});
