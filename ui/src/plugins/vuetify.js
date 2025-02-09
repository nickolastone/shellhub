import Vue from 'vue';
import Vuetify from 'vuetify/lib/framework';
import '@mdi/font/css/materialdesignicons.css';
import 'font-logos/assets/font-logos.css';
import { library } from '@fortawesome/fontawesome-svg-core';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { fas } from '@fortawesome/free-solid-svg-icons';
import { fab } from '@fortawesome/free-brands-svg-icons';
import Clipboard from 'v-clipboard';
import Fragment from 'vue-fragment';
import { StripePlugin } from '@vue-stripe/vue-stripe';
import env from '@/env';

Vue.use(env);

const options = {
  pk: Vue.env.stripePublishableKey,
};

// import '../styles/variables.scss'

Vue.component('FontAwesomeIcon', FontAwesomeIcon); // Register component globally
library.add(fas, fab); // Include needed icons.

Vue.use(Vuetify);
Vue.use(Clipboard);
Vue.use(Fragment.Plugin);

if (Vue.env.billingEnable) {
  Vue.use(StripePlugin, options);
}

export default new Vuetify({
  iconfont: 'md',
  theme: {
    dark: true,

    themes: {
      dark: {
        primary: '#667acc',
        secondary: '#1E2127',
        background: '#18191B',
        tabs: '#1E1E1E',
        foreground: '#1E1E1E',
        paymentForm: '#E0E0E0',
      },
      light: {
        primary: '#667acc',
        secondary: '#c4c4c4',
        background: '#FFFFFF',
        tabs: '#FFFFFF',
        foreground: '#FFFFFF',
        paymentForm: '#FFFFFF',
      },
    },
  },
  options: {
    customProperties: true,
  },
  breakpoint: {
    mobileBreakpoint: 'xs',
  },
});
