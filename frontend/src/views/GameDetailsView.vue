<script setup lang="ts">
import {useRoute} from 'vue-router';
import {computed, ref, type Ref, watch} from 'vue';
import {Game, useGamesStore} from "@stores/games.ts";
import {Mod, useModsStore} from "@stores/mods.ts";
import {useInfiniteScroll} from "@vueuse/core";

const route = useRoute();
const gamesStore = useGamesStore();
const modsStore = useModsStore();

const el = ref<HTMLElement | null>(null);
const isLoading = ref(true);
const modsLoading = ref(false);
const gameDetails: Ref<Game | null> = ref(null);

const mods: Ref<Mod[]> = ref([]);

const reset = () => {
    isLoading.value = true;
    modsLoading.value = false;
    gameDetails.value = null;
    mods.value = [];
};

const loadGameDetails = async (gameId: number) => {
    try {
        isLoading.value = true;
        gameDetails.value = await gamesStore.fetchGameById(gameId);
        isLoading.value = false;
    } catch (error) {
        console.error("Failed to fetch game details:", error);
        isLoading.value = false;
    }
};

const canLoadMore = computed(() => {
    return !modsLoading.value && modsStore.getModCount(gameDetails.value?.id || 0) > mods.value.length;
});

const loadMoreMods = async (gameId: number) => {
    if (modsLoading.value) {
        console.error("Already loading mods, please wait.");
        return;
    }
    modsLoading.value = true;

    try {
        const newMods = await modsStore.fetchMods(
            gameId,
            mods.value.length,
            20
        );
        if (newMods.length > 0) {
            mods.value.push(...newMods);
        }
    } catch (error) {
        console.error("Failed to fetch mods:", error);
    } finally {
        modsLoading.value = false;
    }
};

watch(
    () => route.params.id,
    (newId) => {
        const gameId = parseInt(newId as string);
        reset();
        loadGameDetails(gameId);
        loadMoreMods(gameId);
    },
    {immediate: true}
)

useInfiniteScroll(
    el,
    () => loadMoreMods(gameDetails.value?.id || 0),
    {
        distance: 10,
        canLoadMore: () => {
            return canLoadMore.value;
        },
    }
)
</script>

<template>
    <div class="flex gap-8 items-start" v-if="!isLoading && gameDetails">
        <!-- Game Info (Left) -->
        <div class="w-full max-w-sm bg-base-300 rounded-xl shadow p-6 flex flex-col items-center">
            <img class="w-48 h-64 object-cover rounded-lg mb-4"
                 :src="gameDetails?.capsule"
                 :alt="gameDetails?.name"
            />
            <h1 class="text-3xl font-bold mb-2">
                {{ gameDetails?.name }}
            </h1>
            <p class="text-base-content/70 text-center">
                {{ gameDetails?.description }}
            </p>
        </div>
        <!-- Mods Grid (Right) -->
        <div class="flex-1 overflow-y-auto p-4" ref="el">
            <h2 class="text-2xl font-semibold mb-4">Mods</h2>
            <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
                <div
                    v-for="mod in mods"
                    :key="mod.id"
                    class="card bg-base-100 shadow rounded-lg p-4"
                >
                    <div class="font-bold text-lg mb-1">{{ mod.name }}</div>
                    <div class="text-base-content/60">
                        {{ mod.description }}
                    </div>
                    <div class="mt-2">
                        <a
                            class="btn btn-primary btn-sm"
                        >
                            Install
                        </a>
                    </div>
                </div>
            </div>
            <div v-if="canLoadMore && !modsLoading" class="text-center mt-4">
                <button
                    class="btn btn-primary"
                    @click="loadMoreMods(gameDetails?.id || 0)"
                >
                    Load More Mods
                </button>
            </div>
            <div v-if="modsLoading" class="flex justify-center mt-4">
                <span class="loading loading-spinner loading-md"></span>
                <span class="ml-2">Loading more mods...</span>
            </div>
            <div v-if="!canLoadMore && !modsLoading && mods.length === 0" class="text-center mt-4">
                <p class="text-base-content/70">No mods available for this game.</p>
            </div>
        </div>
    </div>
    <div v-else-if="!gameDetails && !isLoading" class="flex items-center justify-center h-screen">
        <div class="text-center">
            <h2 class="text-2xl font-semibold mb-4">Game Not Found</h2>
            <p class="text-base-content/70">The game you are looking for does not exist or has been removed.</p>
        </div>
    </div>
    <div v-else class="flex items-center justify-center h-screen">
        <span class="loading loading-spinner loading-lg"></span>
        <span class="ml-4 text-lg">Loading game details...</span>
    </div>
</template>

<style scoped lang="scss">

</style>
