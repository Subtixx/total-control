import {fileURLToPath, URL} from "node:url";
import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import VueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from '@tailwindcss/vite'
import {visualizer} from 'rollup-plugin-visualizer'
import Icons from 'unplugin-icons/vite'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        tailwindcss(),
        Icons({autoInstall: true}),
        VueDevTools(),
        visualizer({open: false, gzipSize: true}),
    ],
    resolve: {
        alias: {
            "@": fileURLToPath(new URL("./src", import.meta.url)),
            "@stores": fileURLToPath(new URL("./src/stores", import.meta.url)),
            "@components": fileURLToPath(new URL("./src/components", import.meta.url)),
            "@assets": fileURLToPath(new URL("./src/assets", import.meta.url)),
            "@views": fileURLToPath(new URL("./src/views", import.meta.url)),
            "@router": fileURLToPath(new URL("./src/router", import.meta.url)),
            "@composables": fileURLToPath(new URL("./src/composables", import.meta.url)),
        },
    },
});
