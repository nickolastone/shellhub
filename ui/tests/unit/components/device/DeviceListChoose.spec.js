import Vuex from 'vuex';
import { mount, createLocalVue } from '@vue/test-utils';
import Vuetify from 'vuetify';
import DeviceListChoose from '@/components/device/DeviceListChooser';

describe('DeviceListChoose', () => {
  const localVue = createLocalVue();
  const vuetify = new Vuetify();
  localVue.use(Vuex);

  let wrapper;

  const numberDevices = 2;
  const devicesSelected = [];
  const action = 'suggestedDevices';

  const pagination = {
    groupBy: [],
    groupDesc: [],
    itemsPerPage: 10,
    multiSort: false,
    mustSort: false,
    page: 1,
    sortBy: [],
    sortDesc: [false],
  };

  const headers = [
    {
      text: 'Hostname',
      value: 'hostname',
      align: 'center',
    },
    {
      text: 'Operating System',
      value: 'info.pretty_name',
      align: 'center',
      sortable: false,
    },
    {
      text: 'SSHID',
      value: 'namespace',
      align: 'center',
      sortable: false,
    },
  ];

  const devices = [
    {
      uid: 'a582b47a42d',
      name: '39-5e-2a',
      identity: {
        mac: '00:00:00:00:00:00',
      },
      info: {
        id: 'linuxmint',
        pretty_name: 'Linux Mint 19.3',
        version: '',
      },
      public_key: '----- PUBLIC KEY -----',
      tenant_id: '00000000',
      last_seen: '2020-05-20T18:58:53.276Z',
      online: false,
      namespace: 'user',
      status: 'accepted',
    },
    {
      uid: 'a582b47a42e',
      name: '39-5e-2b',
      identity: {
        mac: '00:00:00:00:00:00',
      },
      info: {
        id: 'linuxmint',
        pretty_name: 'Linux Mint 19.3',
        version: '',
      },
      public_key: '----- PUBLIC KEY -----',
      tenant_id: '00000001',
      last_seen: '2020-05-20T19:58:53.276Z',
      online: true,
      namespace: 'user',
      status: 'accepted',
    },
  ];

  const store = new Vuex.Store({
    namespaced: true,
    state: {
      devices,
      numberDevices,
      devicesSelected,
    },
    getters: {
      'devices/getDevicesForUserToChoose': (state) => state.devices,
      'devices/getNumberForUserToChoose': (state) => state.numberDevices,
      'devices/getDevicesSelected': (state) => state.devicesSelected,
    },
    actions: {
      'devices/fetch': () => {},
      'devices/setDevicesSelected': () => {},
      'devices/getDevicesMostUsed': () => {},
      'devices/setDevicesForUserToChoose': () => {},
      'devices/getNumberForUserToChoose': () => {},
      'devices/resetListDevices': () => {},
      'snackbar/showSnackbarCopy': () => {},
      'snackbar/showSnackbarErrorAssociation': () => {},
      'snackbar/showSnackbarErrorLoading': () => {},
    },
  });

  ///////
  // In this case, it is tested when device is online.
  ///////

  describe('Suggested devices', () => {
    beforeEach(() => {
      wrapper = mount(DeviceListChoose, {
        store,
        localVue,
        stubs: ['fragment', 'router-link'],
        propsData: {
          action,
        },
        vuetify,
      });
    });

    ///////
    // Component Rendering
    //////

    it('Is a Vue instance', () => {
      expect(wrapper).toBeTruthy();
    });
    it('Renders the component', () => {
      expect(wrapper.html()).toMatchSnapshot();
    });

    ///////
    // Data and Props checking
    //////

    it('Compare data with default value', () => {
      expect(wrapper.vm.hostname).toEqual('localhost');
      expect(wrapper.vm.pagination).toEqual(pagination);
      expect(wrapper.vm.headers).toEqual(headers);
    });
    it('Process data in the computed', () => {
      expect(wrapper.vm.getListDevices).toEqual(devices);
      expect(wrapper.vm.getNumberDevices).toEqual(numberDevices);
    });
    it('Process data in methods', () => {
      Object.keys(devices).forEach((device) => {
        const address = `${device.namespace}.${device.name}@localhost`;
        expect(wrapper.vm.address(device)).toEqual(address);
      });
    });

    //////
    // HTML validation
    //////

    it('Renders the template with components', async () => {
      expect(wrapper.find('[data-test="deviceIcon-component"]').exists()).toEqual(true);
    });
    it('Renders the template with data', () => {
      const dt = wrapper.find('[data-test="devices-dataTable"]');
      const dataTableProps = dt.vm.$options.propsData;

      expect(dataTableProps.items).toHaveLength(numberDevices);
    });
  });
});
