import { http, HttpResponse } from 'msw';

export const handlers = [
  http.post('/homegym/login', async ({ request }) => {
    const options = { status: 200 };
    const creds = await request.json();
    if (creds.password == 'badpassword') {
      options.status = 401;
    }
    return HttpResponse.json(null, options);
  }),
];
