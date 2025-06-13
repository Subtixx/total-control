<script setup lang="ts">
import LoadingIcon from '~icons/svg-spinners/blocks-shuffle-3'

import {Game, useGamesStore} from "@stores/games.ts";
import {computed, ref, type Ref} from "vue";

const gameStore = useGamesStore();

const isLoading: Ref<boolean> = ref(false);
const searchQuery: Ref<string> = ref('');
const searchResults: Ref<Game[]> = ref([]);
const searchGames = (query: string) => {
    if (query.length < 2) {
        searchResults.value = [];
        return;
    }
    if (isLoading.value) return; // Prevent multiple requests while loading
    searchResults.value = [];
    isLoading.value = true;
    gameStore.searchGames(query).then((results) => {
        searchResults.value = results.sort((a, b) => {
            const nameA = a.name.toLowerCase();
            const nameB = b.name.toLowerCase();
            if (nameA < nameB) return -1;
            if (nameA > nameB) return 1;
            return 0;
        });
        isLoading.value = false;
    }).catch((error) => {
        console.error("Error searching games:", error);
        searchResults.value = [];
    });
};

const searchResultsMax = computed(() => {
    return searchResults.value.slice(0, 5);
});

const showResults = computed(() => {
    return isLoading.value || searchResults.value.length > 0 || searchQuery.value.length >= 2;
});

const clearSearch = () => {
    searchQuery.value = '';
    searchResults.value = [];
};

const highlightMatch = (name: string, query: Ref<string> | string) => {
    const q = typeof query === 'string' ? query : query.value;
    if (!q) return name;
    const regex = new RegExp(`(${q.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'ig');
    return name.replace(regex, '<b class="text-accent">$1</b>');
};
</script>

<template>
    <div class="flex flex-grow relative group">
        <label
            class="input flex-grow rounded-none border-none focus:outline-none focus:ring-0 focus-within:outline-none focus-within:ring-0 focus-within:border-none">
            <svg class="h-[1em] opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <g
                    stroke-linejoin="round"
                    stroke-linecap="round"
                    stroke-width="2.5"
                    fill="none"
                    stroke="currentColor"
                >
                    <circle cx="11" cy="11" r="8"></circle>
                    <path d="m21 21-4.3-4.3"></path>
                </g>
            </svg>
            <input
                v-model="searchQuery"
                @input="searchGames(searchQuery)"
                @keydown.enter="searchGames(searchQuery)"
                @keydown.esc="clearSearch"
                type="search"
                class="grow"
                placeholder="Search"/>
            <kbd class="kbd kbd-sm">ctrl</kbd>
            <kbd class="kbd kbd-sm">K</kbd>
        </label>
        <ul
            v-if="showResults"
            class="bg-base-100 text-base-content z-9999 absolute top-8 left-0 w-full mt-2 rounded-b shadow-lg max-h-60 overflow-y-auto opacity-0 group-focus-within:opacity-100 transition-opacity duration-300">
            <li
                v-for="result in searchResultsMax"
                :key="result.id"
                class="">
                <router-link
                    :to="{ name: 'GameDetails', params: { id: result.id } }"
                    class="flex gap-2 border-t-2 border-base-200 hover:bg-primary hover:text-primary-content px-2 py-2 relative cursor-pointer"
                    @click="clearSearch">
                    <img :src="result.icon"
                         :alt="result.name"
                         class="w-6 h-6 rounded-full">
                    <span
                        class="flex-grow"
                        v-html="highlightMatch(result.name, searchQuery)"/>
                    <span class="badge badge-secondary">
                        {{ result.gamePath ? result.gamePath : 'N/A' }}
                    </span>
                </router-link>
            </li>
            <li
                v-if="isLoading"
                class="px-2 py-2 text-gray-500 flex items-center">
                <LoadingIcon class="w-4 h-4 mr-2 inline-block"/>
                Loading...
            </li>
            <li
                v-else-if="searchResults.length > 5 || searchQuery.length < 2 || searchResults.length === 0"
                class="px-2 py-2 text-gray-500">
                <span v-if="searchQuery.length > 2 && searchResults.length === 0">
                    No results found
                </span>
                <span v-else-if="searchQuery.length < 2">
                    Start typing to search
                </span>
                <span v-else-if="searchResults.length > 5">
                    {{ searchResults.length - 5 }} more..
                </span>
            </li>
        </ul>
    </div>
</template>
