<script setup lang="ts">
import {type Ref, ref} from "vue";

import SettingsIcon from '~icons/mdi/cog';
import PlusIcon from '~icons/mdi/plus';
import PuzzlePieceIcon from '~icons/heroicons/puzzle-piece';

import {useGamesStore} from "@stores/games.ts";
import {games} from "@wails/go/models.ts";

const gamesStore = useGamesStore();

//curl 'https://www.steamgriddb.com/api/public/game/10052' --compressed -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:138.0) Gecko/20100101 Firefox/138.0' -H 'Accept: application/json, text/plain, */*' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br, zstd' -H 'Connection: keep-alive' -H 'Referer: https://www.steamgriddb.com/game/10052/icons' -H 'Cookie: cookieAgree=true' -H 'Sec-Fetch-Dest: empty' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Site: same-origin' -H 'TE: trailers'
const installedGamesLoading = ref(true);
const installedGames: Ref<games.Game[]> = ref([]);
gamesStore.fetchInstalledGames().then((games) => {
    installedGames.value = games;
    installedGamesLoading.value = false;
}).catch(() => {
    installedGamesLoading.value = false;
});
</script>

<template>
    <div class="flex flex-col flex-shrink-0 w-56 bg-base-100 text-base-content">
        <div class="flex items-center justify-center h-16 px-4 border-b border-base-200 bg-base-300">
            <router-link to="/" class="flex items-center">
                <img src="@/assets/logo_side.svg" alt="Logo" class="h-10 mr-2">
            </router-link>
        </div>
        <div class="flex flex-col flex-grow p-4 overflow-auto">
            <router-link
                v-for="game in installedGames"
                :key="game.name"
                class="flex items-center flex-shrink-0 h-10 px-2 text-sm font-medium rounded hover:bg-primary hover:text-primary-content"
                :to="{ name: 'GameDetails', params: { id: game.id } }"
            >
                <img :src="game.icon" alt="Game Icon" class="w-6 h-6 mr-2 inline-block">
                {{ game.name }}
            </router-link>
        </div>
        <div class="p-4">
            <router-link
                class="flex items-center flex-shrink-0 h-10 px-4 text-sm font-medium rounded btn btn-success"
                :to="{ name: 'AddGame' }">
                <PlusIcon class="w-5 h-5 ml-1 inline-block"/>
                Add Game
            </router-link>
        </div>

        <div class="flex items-center h-16 border-t border-base-200 bg-base-300 px-4">
            <!--Version number-->
            <span class="flex-grow text-sm font-medium">Version 1.0.0</span>
            <router-link
                :to="{ name: 'Plugins' }"
                class="flex items-center justify-center w-10 h-10 text-sm font-medium rounded hover:bg-primary hover:text-primary-content">
                <PuzzlePieceIcon/>
            </router-link>
            <router-link
                class="flex items-center justify-center w-10 h-10 text-sm font-medium rounded hover:bg-primary hover:text-primary-content"
                :to="{name:'Settings'}">
                <SettingsIcon class="w-5 h-5"/>
            </router-link>
        </div>
    </div>
</template>

<style scoped>

</style>
