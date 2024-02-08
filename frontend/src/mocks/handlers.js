import { http, HttpResponse } from 'msw';
import { fetchEventExercises } from '../modules/utils';
import * as data from 'data.js';

const loginHandler = http.post('/homegym/login', async ({ request }) => {
  const options = { status: 200 };
  const creds = await request.json();
  if (creds.password == 'badpassword') {
    options.status = 401;
  }
  return HttpResponse.json(null, options);
});

const eventPageHandler = http.get('/homegym/api/events/', ({ params }) => {
  const { previousID, date } = params;

  return HttpResponse.json(data.fetchedEvents(10), { status: 200 });
});
export const handlers = [loginHandler, eventPageHandler];
