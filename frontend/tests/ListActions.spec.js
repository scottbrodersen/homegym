import { config, mount } from '@vue/test-utils';
import ListActions from './../src/components/ListActions.vue';
import { OrderedList } from '../src/modules/utils';
import { Quasar } from 'quasar';

config.global.plugins.push(Quasar);
describe('ListActions', () => {
  it('The correct nubmer of buttons are rendered', () => {
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
