import { setupServer } from 'msw/node';
import { handlers } from './handlers.js';

// Configures a Service Worker with the given request handlers.
export default setupServer(...handlers);
