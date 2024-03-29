import { installQuasarPlugin } from '@quasar/quasar-app-extension-testing-unit-vitest';
import { config, mount, enableAutoUnmount } from '@vue/test-utils';
import BinaryRep from './../src/components/BinaryRep.vue';

installQuasarPlugin();

config.global.errorHandler = (err) => {
  throw err;
};

enableAutoUnmount(afterEach);

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
