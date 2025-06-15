<script setup lang="ts">

import {onMounted, ref, type Ref} from "vue";

let availableThemes = [
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

const activeTheme: Ref<string | null> = ref(null);
onMounted(() => {
    activeTheme.value = localStorage.getItem('theme') || 'system';
    setTheme(activeTheme.value);
});

const setTheme = (theme: string) => {
    if (theme === 'system') {
        theme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    }
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
};
</script>

<template>
    <select
        class="select w-full"
        @change="setTheme($event?.target?.value)">
        <option
            v-for="theme in availableThemes"
            :selected="activeTheme === theme"
            :key="theme"
            :value="theme">
            {{ theme }}
        </option>
        <option value="system">System Default</option>
    </select>
</template>
