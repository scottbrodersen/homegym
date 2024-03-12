import { mount } from '@vue/test-utils';
import BinaryRep from './../src/components/BinaryRep.vue';

describe('BinaryRep', () => {
  test('renders a good rep', () => {
    const wrapper = mount(BinaryRep, { props: true });

    expect(wrapper.find('span').exists()).toBe(true);
    expect(wrapper.findAll('span')).toHaveLength(1);
    const span = wrapper.get('span');
    expect(span.text()).toBe('circle');
  });
  test('renders a bad rep', () => {
    const wrapper = mount(BinaryRep, { props: false });

    expect(wrapper.find('span').exists()).toBe(true);
    expect(wrapper.findAll('span')).toHaveLength(1);
    const span = wrapper.get('span');
    expect(span.text()).toBe('circle');
  });
});
