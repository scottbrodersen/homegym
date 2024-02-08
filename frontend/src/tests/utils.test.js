import './patch-fetch.js';
import { expect, test, beforeAll, afterAll, afterEach } from 'vitest';
import * as utils from '../modules/utils.js';
import server from './../mocks/server.js';

beforeAll(() => {
  server.listen({ onUnhandledRequest: 'error' });
});
afterAll(() => server.close());
afterEach(() => server.resetHandlers());

test('hrZone validation', () => {
  const goodValues = [1, 2, 3, 4, 5];
  const badValues = [0, 6, 'r', '', '$'];

  for (const good of goodValues) {
    expect(utils.intensityProps('hrZone').validate(good)).toBeTruthy;
  }

  for (const bad of badValues) {
    expect(utils.intensityProps('hrZone').validate(bad)).toBeFalsy;
  }
});

test('when we log in using good credentials', async () => {
  let caught = false;
  try {
    const resp = await utils.login('id', 'password');
  } catch (error) {
    caught = true;
  }
  expect(caught).toBeFalsy;
});

test('when we log in using bad credentials', async () => {
  let caught = false;
  try {
    const resp = await utils.login('id', 'badpassword');
  } catch (error) {
    caught = true;
  }
  expect(caught).toBeTruthy;
});
