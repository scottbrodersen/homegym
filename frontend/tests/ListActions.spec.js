import { config, enableAutoUnmount, mount } from '@vue/test-utils';
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

  it('The menu items are rendered', () => {
    const wrapper = mount(ListActions);

    console.log('listActions:');
    console.log(wrapper.html());

    wrapper.get('button').trigger('click');
    const menuWrapper = wrapper.getComponent(QMenu);
    expect(menuWrapper).to.not.be.null;
    console.log('menuwrapper:');

    console.log(menuWrapper.html());

    expect(menuWrapper.findAll('div.q-list')).toHaveLength(6);
  });

  it('The correct events are emitted', () => {
    const wrapper = mount(ListActions);

    wrapper.findAll('div.q-item').forEach((btn) => {
      btn.trigger('click');
    });

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
