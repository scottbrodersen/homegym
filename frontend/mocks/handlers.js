import { http, HttpResponse } from 'msw';
import * as data from './data';

const loginHandler = http.post('/homegym/login', async ({ request }) => {
  const options = { status: 200 };
  const credentials = await request.json();
  if (credentials.password == 'badpassword') {
    options.status = 401;
  }
  return HttpResponse.json(null, options);
});

const eventPageHandler = http.get('/homegym/api/events/', ({ params }) => {
  const { previousID, date } = params;

  return HttpResponse.json(data.fetchedEvents(10), { status: 200 });
});

const programInstanceHandler = http.post(
  '/homegym/api/activities/:activityID/programs/:programID/instances/',
  ({ params }) => {
    let { instID } = params;
    return HttpResponse.json(data.testProgramInstance(data.testActivityID), {
      status: 200,
    });
  }
);

export const handlers = [
  loginHandler,
  eventPageHandler,
  programInstanceHandler,
];
