// import ProgramInstanceModal from '../src/components/ProgramInstanceModal.vue';
// import { config, mount, DOMWrapper } from '@vue/test-utils';
// import {
//   useQuasar,
//   Dialog,
//   Notify,
//   Quasar,
//   useDialogPluginComponent,
//   QCard,
//   QCardActions,
//   QDate,
//   QDialog,
//   QInput,
// } from 'quasar';
// import { installQuasarPlugin } from '@quasar/quasar-app-extension-testing-unit-vitest';
// import { focus } from '../src/modules/directives';
// import { expect } from 'vitest';

// config.global.plugins.push(Quasar);
// config.global.directives = {
//   focus: focus,
// };

// installQuasarPlugin({
//   components: { QCard, QCardActions, QDate, QDialog, QInput },
//   plugins: { Dialog },
// });

// describe('ProgramInstanceModal component', () => {
//   // it('should mount the document body and expose for testing', () => {
//   // const d = Dialog.create({
//   //   component: ProgramInstanceModal,
//   //   componentProps: {
//   //     activityID: 'test-activity-id',
//   //     programID: 'test-program-id',
//   //   },
//   // });
//   //     const dialogWrapper = mount(ProgramInstanceModal, {
//   //       components: { Dialog, QCard, QCardActions, QDate, QDialog, QInput },
//   //       props: {
//   //         activityID: 'test-activity-id',
//   //         programID: 'test-program-id',
//   //       },
//   //       data: () => ({
//   //         isDialogOpen: true,
//   //       }),
//   //     });
//   //     const wrapper = new DOMWrapper(document.body);

//   //     expect(wrapper.find('.q-dialog').exists()).toBeTruthy();
//   //   });

//   it('renders correctly', () => {
//     document.body.innerHTML = `
//   <div>
//     <h1>Non Vue app</h1>
//     <div id="app"></div>
//   </div>`;
//     const $q = useQuasar;

//     const d = $q.dialog({
//       component: ProgramInstanceModal,
//       componentProps: {
//         activityID: 'test-activity-id',
//         programID: 'test-program-id',
//       },
//     });

//     const wrapper = mount(ProgramInstanceModal, {
//       attachTo: document.getElementById('app'),

//       components: { Dialog, QCard, QCardActions, QDate, QDialog, QInput },
//       props: {
//         activityID: 'test-activity-id',
//         programID: 'test-program-id',
//       },
//     });
//     const input = wrapper.findAllComponents(QInput);
//     expect(input).toHaveLength(1);

//     const date = wrapper.findAllComponents(QDate);
//     expect(date).toHaveLength(1);
//   });
// });
