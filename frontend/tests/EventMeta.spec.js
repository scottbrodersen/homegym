import EventMeta, { labels } from '../src/components/EventMeta.vue';
import { config, mount } from '@vue/test-utils';
import { Quasar, QInput, QRating } from 'quasar';
import { focus, select } from '../src/modules/directives';

config.global.plugins.push(Quasar);
config.global.directives = {
  focus: focus,
  select: select,
};
config.global.errorHandler = (err) => {
  throw err;
};

describe('EventMeta component', () => {
  it('renders correctly with all meta in read only', () => {
    const testMeta = {
      mood: 1,
      energy: 2,
      motivation: 3,
      overall: 4,
      notes: 'test notes',
    };

    const wrapper = mount(EventMeta, {
      components: { QInput, QRating },
      props: {
        mood: testMeta.mood,
        energy: testMeta.energy,
        motivation: testMeta.motivation,
        overall: testMeta.overall,
        notes: testMeta.notes,
        readonly: true,
      },
    });
    let found = [];
    const elements = wrapper.findAll('div');
    for (const meta in testMeta) {
      for (let i = 0; i < elements.length; i++) {
        if (elements[i].text() == `${labels[meta]}: ${testMeta[meta]}`) {
          found.push(meta);
          break;
        }
      }
    }

    for (const meta in testMeta) {
      expect(found.includes(meta)).toEqual(true);
    }
  });

  it('renders correctly with a subset of meta in read only', () => {
    const testMeta = {
      mood: 1,
      energy: 2,
      overall: 4,
    };

    const wrapper = mount(EventMeta, {
      components: { QInput, QRating },
      props: {
        mood: testMeta.mood,
        energy: testMeta.energy,
        overall: testMeta.overall,
        readonly: true,
      },
    });
    let found = [];
    const elements = wrapper.findAll('div');
    for (const meta in testMeta) {
      for (let i = 0; i < elements.length; i++) {
        if (elements[i].text() == `${labels[meta]}: ${testMeta[meta]}`) {
          found.push(meta);
          break;
        }
      }
    }

    for (const meta in testMeta) {
      expect(found.includes(meta)).toEqual(true);
    }
  });

  it('renders correctly not in read only', () => {
    const testMeta = {
      mood: 1,
      energy: 2,
      overall: 4,
    };

    const wrapper = mount(EventMeta, {
      components: { QInput, QRating },
      props: {
        readonly: false,
      },
    });
    let found = [];
    const ratings = wrapper.findAllComponents(QRating);

    expect(ratings).toHaveLength(4);
  });
});
