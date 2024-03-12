// workaround for https://github.com/mswjs/msw/issues/1625
const { fetch: globalFetch } = global;
global.fetch = async (...args) => {
  let [resource, options] = args;
  if (resource.startsWith('/')) {
    resource = 'http://localhost:3000' + resource;
  }
  const response = await globalFetch(resource, options);
  console.log(response);
  return response;
};
