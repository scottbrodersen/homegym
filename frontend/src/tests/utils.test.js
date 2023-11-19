import './patch-fetch.js';
import { expect, test, beforeAll, afterAll, afterEach } from 'vitest';
import * as u from '../modules/utils.js';
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
    expect(u.intensityTypeProps.hrZone.validate(good)).toBeTruthy;
  }

  for (const bad of badValues) {
    expect(u.intensityTypeProps.hrZone.validate(bad)).toBeFalsy;
  }
});

test('when we log in using good credentials', async () => {
  let caught = false;
  try {
    const resp = await u.login('id', 'password');
  } catch (error) {
    caught = true;
  }
  expect(caught).toBeFalsy;
});

test('when we log in using bad credentials', async () => {
  let caught = false;
  try {
    const resp = await u.login('id', 'badpassword');
  } catch (error) {
    caught = true;
  }
  expect(caught).toBeTruthy;
});
