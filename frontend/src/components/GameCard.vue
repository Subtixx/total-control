<script setup lang="ts">
import type {Game} from "@stores/games.ts";

type Props = {
    game?: Game;
    enableLink?: boolean;
    showBadge?: boolean;
};

defineProps<Props>();
</script>

<template>
    <div v-if="game" :key="game.id" class="items-center">
        <figure>
            <div class="w-64 h-96 relative group">
                <router-link
                    :to="{ name: 'GameDetails', params: { id: game.id } }"
                    class="block" v-if="enableLink">
                    <img :src="game.capsule"
                         alt="" class="rounded-xl object-cover max-w-64">
                </router-link>
                <img :src="game.capsule"
                     alt="" class="rounded-xl object-cover max-w-64" v-else>
                <div class="absolute top-2 right-2 badge badge-error" v-if="showBadge">
                </div>
                <div
                    class="absolute bottom-0 left-0 w-full text-center bg-black bg-opacity-50 text-white p-2 rounded-b-xl opacity-0 group-hover:opacity-100 transition-opacity">
                    {{ game.name }}
                </div>
            </div>
        </figure>
    </div>
    <!-- Skeleton Loader -->
    <div v-else class="items-center animate-pulse rounded-xl w-64 h-96">
        <div class="skeleton rounded-xl w-64 h-96"></div>
    </div>
</template>

<style scoped lang="scss">

</style>
