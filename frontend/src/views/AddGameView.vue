<script setup lang="ts">
import {ref} from 'vue';
import {useRouter} from 'vue-router';
import {useDebounceFn} from '@vueuse/core';
import {Game, useGamesStore} from '@stores/games';
import GameCard from "@components/GameCard.vue";

const router = useRouter();
const gameStore = useGamesStore();

const gamePath = ref('');
const isLoading = ref(false);
const foundGame = ref<Game | null>(null);
const addGame = async () => {
    if (!gamePath.value || !foundGame.value) {
        alert('Please fill in all fields.');
        return;
    }

    try {
        await gameStore.addGame(foundGame.value);
        alert('Game added successfully!');
        router.push({name: 'home'});
    } catch (error) {
        console.error('Error adding game:', error);
        alert('Failed to add game. Please try again.');
    }
};

const fetchGame = async () => {
    if (!gamePath.value) {
        foundGame.value = null;
        return;
    }
    if (isLoading.value) return; // Prevent multiple fetches
    if (foundGame.value && foundGame.value.gamePath === gamePath.value) {
        return; // No need to fetch if the path hasn't changed
    }

    isLoading.value = true;

    try {
        const game = await gameStore.detectGame(gamePath.value);
        if (game) {
            foundGame.value = game;
        } else {
            foundGame.value = null;
        }
    } catch (error) {
        console.error('Error fetching game:', error);
        alert('Unknown error occurred while fetching the game.');
    }
    isLoading.value = false;
};
</script>

<template>
    <div class="flex gap-4 items-start">
        <div class="card bg-base-300">
            <div class="w-md card-body">
                <h2 class="card-title justify-center mb-4">Add Game From Disk</h2>
                <div class="mb-4">
                    <fieldset class="fieldset">
                        <legend class="fieldset-legend">
                            Path to Game
                        </legend>
                        <input
                            :disabled="isLoading"
                            v-model="gamePath"
                            @keyup.enter="fetchGame"
                            @blur="fetchGame"
                            type="text"
                            placeholder="Enter path to game"
                            class="input input-bordered w-full"
                            :class="{'input-error': gamePath && foundGame == null}"
                            required
                        />
                        <p class="label text-error" v-if="gamePath && foundGame == null">
                            No supported game found at this path.
                        </p>
                    </fieldset>
                </div>
                <div class="card-actions justify-between">
                    <button
                        class="btn btn-error"
                    >
                        Cancel
                    </button>
                    <button
                        :disabled="foundGame == null || isLoading"
                        @click="addGame"
                        class="btn btn-primary">
                        Add Game
                    </button>
                </div>
            </div>
        </div>
        <GameCard v-if="!isLoading && foundGame" :game="foundGame"/>
        <GameCard v-else-if="isLoading"/>
    </div>
</template>

<style scoped lang="scss">

</style>
