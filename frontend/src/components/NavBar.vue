<script setup lang="ts">
import SearchBar from "@components/SearchBar.vue";
import {useRouter} from "vue-router";

const router = useRouter();
</script>

<template>
    <div class="flex items-center gap-4 flex-shrink-0 h-16 px-4 border-b border-base-200 bg-base-300">
        <SearchBar/>
        <button
            v-if="false"
            class="flex items-center justify-center h-10 px-4 text-sm font-medium rounded btn hover:btn-primary">
            Action 1
        </button>
        <button
            v-if="false"
            class="flex items-center justify-center h-10 px-4 ml-2 text-sm font-medium rounded btn btn-primary">
            Action 2
        </button>
        <button
            v-if="false"
            class="relative ml-2 text-sm focus:outline-none group">
            <div class="flex items-center justify-between rounded btn hover:btn-primary">
                <svg class="w-5 h-5 mx-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                     stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"/>
                </svg>
            </div>
            <div
                class="absolute right-0 pt-2 flex-col items-start hidden w-40 pb-1 border shadow-lg group-focus:flex bg-base-300 border-base-100">
                <a class="w-full px-4 py-2 text-left hover:bg-primary hover:text-primary-content" href="#">Menu Item
                    1</a>
                <a class="w-full px-4 py-2 text-left hover:bg-primary hover:text-primary-content" href="#">Menu Item
                    2</a>
                <a class="w-full px-4 py-2 text-left hover:bg-primary hover:text-primary-content" href="#">Menu Item
                    3</a>
            </div>
        </button>
    </div>
</template>

