<script setup lang="ts">
import {useI18n} from "vue-i18n";
import {computed, onMounted, ref} from "vue";
import {useGamesStore} from "@/stores/games";
import GameCard from "@components/GameCard.vue";

const {t} = useI18n();

const gamesStore = useGamesStore();
const recentGamesPage = ref(0);
const allGamesPage = ref(0);
const gamesPerPage = 4;
const recentGamesToShow = computed(() => {
    const start = recentGamesPage.value * gamesPerPage;
    return gamesStore.recentGames.slice(start, start + gamesPerPage);
});
const allGamesToShow = computed(() => {
    const start = allGamesPage.value * gamesPerPage;
    return gamesStore.games.slice(start, start + gamesPerPage);
});
const nextRecent = () => {
    if ((recentGamesPage.value + 1) * gamesPerPage < gamesStore.recentGames.length) {
        recentGamesPage.value++;
    }
};
const prevRecent = () => {
    if (recentGamesPage.value > 0) {
        recentGamesPage.value--;
    }
};
const nextAll = () => {
    if ((allGamesPage.value + 1) * gamesPerPage < gamesStore.games.length) {
        allGamesPage.value++;
    }
};
const prevAll = () => {
    if (allGamesPage.value > 0) {
        allGamesPage.value--;
    }
};

const allGamesLoading = ref(true);
const recentGamesLoading = ref(true);
onMounted(() => {
    gamesStore.fetchRecentGames().then(() => {
        recentGamesLoading.value = false;
    }).catch(() => {
        recentGamesLoading.value = false;
    });
    gamesStore.fetchGames().then(() => {
        allGamesLoading.value = false;
    }).catch(() => {
        allGamesLoading.value = false;
    });
});
</script>

<template>
    <section class="mb-12">
        <div class="flex items-center justify-between mb-4">
            <h2 class="text-2xl font-bold">Recent Games</h2>
            <div v-if="!recentGamesLoading && gamesStore.recentGames.length > gamesPerPage">
                <button class="btn btn-circle btn-ghost"
                        :disabled="recentGamesPage === 0"
                        @click="prevRecent">
                    <span>&lt;</span>
                </button>
                <button class="btn btn-circle btn-ghost"
                        :disabled="recentGamesPage * gamesPerPage + gamesPerPage >= gamesStore.recentGames.length"
                        @click="nextRecent">
                    <span>&gt;</span>
                </button>
            </div>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6"
                v-if="!recentGamesLoading">
            <GameCard
                :enable-link="true"
                v-for="game in recentGamesToShow" :key="game.id" :game="game"/>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6" v-else>
            <GameCard v-for="i in gamesPerPage" :key="i" />
        </div>
    </section>

    <section>
        <div class="flex items-center justify-between mb-4">
            <h2 class="text-2xl font-bold">All Games</h2>
            <div v-if="!allGamesLoading && gamesStore.games.length > gamesPerPage">
                <button class="btn btn-primary btn-circle btn-ghost"
                        :disabled="allGamesPage === 0"
                        @click="prevAll">
                    <span>&lt;</span>
                </button>
                <button class="btn btn-primary btn-circle btn-ghost"
                        :disabled="allGamesPage * gamesPerPage + gamesPerPage >= gamesStore.games.length"
                        @click="nextAll">
                    <span>&gt;</span>
                </button>
            </div>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6" v-if="!allGamesLoading">
            <GameCard
                :enable-link="true"
                v-for="game in allGamesToShow"
                :key="game.id"
                :game="game"/>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6" v-else>
            <GameCard
                v-for="i in gamesPerPage"
                :key="i" />
        </div>
    </section>
</template>

<style scoped>
</style>
