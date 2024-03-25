import './patch-fetch.js';
import { expect, test, beforeAll, afterAll, afterEach } from 'vitest';
import * as utils from '../src/modules/utils.js';
import server from '../mocks/server.js';

describe('utils', () => {
  beforeAll(() => {
    server.listen({ onUnhandledRequest: 'error' });
  });
  afterAll(() => server.close());
  afterEach(() => server.resetHandlers());

  it('hrZone validation', () => {
    const goodValues = [1, 2, 3, 4, 5];
    const badValues = [0, 6, 'r', '', '$'];

    for (const good of goodValues) {
      expect(utils.intensityProps('hrZone').validate(good)).toBeTruthy;
    }

    for (const bad of badValues) {
      expect(utils.intensityProps('hrZone').validate(bad)).toBeFalsy;
    }
  });

  it('hrZone formatting', () => {
    const hrZoneProps = utils.intensityProps('hrZone');
    expect(hrZoneProps.format(2)).toEqual('2');
    expect(hrZoneProps.format(2.1)).toEqual('2');
  });

  it('pace validation', () => {
    const goodValues = ['33:33', '03:33', '00:33'];
    const badValues = ['3:33', '33:3', '3:3', '3333'];
    for (const good of goodValues) {
      expect(utils.intensityProps('pace').validate(good)).toBeTruthy;
    }

    for (const bad of badValues) {
      expect(utils.intensityProps('pace').validate(bad)).toBeFalsy;
    }
  });

  it('pace formatting', () => {
    const formatted = {
      60: '01:00',
      61: '01:01',
      600: '10:00',
      70: '01:10',
    };
    const format = utils.intensityProps('pace').format;
    for (const value in formatted) {
      expect(format(value)).toEqual(formatted[value]);
    }
  });

  it('when we log in using good credentials', async () => {
    let caught = false;
    try {
      const resp = await utils.login('id', 'password');
    } catch (error) {
      caught = true;
    }
    expect(caught).toBeFalsy;
  });

  it('when we log in using bad credentials', async () => {
    let caught = false;
    try {
      const resp = await utils.login('id', 'badpassword');
    } catch (error) {
      caught = true;
    }
    expect(caught).toBeTruthy;
  });
});
