import {createApp} from "vue";
import {createPinia} from "pinia";
import {createThemePlugin} from "@/themePlugin.ts";

import Vue3Toastify, {type ToastContainerOptions} from 'vue3-toastify';

import 'vue3-toastify/dist/index.css';
import App from "./App.vue";
import i18n from "./i18n";
import router from "./router";
import "./style.css";

const availableThemes = [
    "light",
    "dark",
    "cupcake",
    "bumblebee",
    "emerald",
    "corporate",
    "synthwave",
    "retro",
    "cyberpunk",
    "valentine",
    "halloween",
    "garden",
    "forest",
    "aqua",
    "lofi",
    "pastel",
    "fantasy",
    "wireframe",
    "black",
    "luxury",
    "dracula",
    "cmyk",
    "autumn",
    "business",
    "acid",
    "lemonade",
    "night",
    "coffee",
    "winter",
    "dim",
    "nord",
    "sunset",
    "caramellatte",
    "abyss",
    "silk"
];

const app = createApp(App);

app.use(Vue3Toastify, {
    autoClose: 3000,
} as ToastContainerOptions);
app.use(createPinia());
app.use(createThemePlugin({
    defaultTheme: 'dark',
    themes: availableThemes,
    storageKey: 'app-theme'
}));
app.use(router);
app.use(i18n);

app.mount("#app");
