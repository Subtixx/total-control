<script setup lang="ts">
import ArrowPathIcon from '~icons/heroicons/arrow-path.svg';
import MagnifyingGlassIcon from '~icons/heroicons-solid/magnifying-glass.svg';
import CogIcon from '~icons/heroicons-solid/cog.svg';

import {onMounted, type Ref, ref} from 'vue';
import {type Plugin, usePluginsStore} from '@stores/plugins';
import PluginCard from "@components/PluginCard.vue";
import {plugins} from "@wails/go/models.ts";
import PluginRepositoryCard from "@components/PluginRepositoryCard.vue";
import {toast} from "vue3-toastify";
import PluginRepositoryInfo = plugins.PluginRepositoryInfo;

const pluginsStore = usePluginsStore();

const searchQuery = ref('');
const isLoadingInstalledPlugins = ref(false);
const isLoadingAvailablePlugins = ref(false);
const installedPlugins: Ref<Array<Plugin>> = ref([]);
const availablePlugins: Ref<Array<PluginRepositoryInfo>> = ref([]);

const refreshAvailablePlugins = async () => {
    if (isLoadingAvailablePlugins.value) {
        console.warn('Refresh already in progress, skipping...');
        return;
    }
    availablePlugins.value = [];
    isLoadingAvailablePlugins.value = true;
    try {
        availablePlugins.value = await pluginsStore.fetchAvailablePlugins();
    } catch (error) {
        console.error('Error refreshing available plugins:', error);
        toast('Failed to refresh available plugins!', {
            type: 'error',
            position: 'top-right',
            autoClose: 3000,
        });
    }
    searchQuery.value = '';
    isLoadingAvailablePlugins.value = false;
};

const refreshInstalledPlugins = async () => {
    if (isLoadingInstalledPlugins.value) {
        console.warn('Refresh already in progress, skipping...');
        return;
    }
    installedPlugins.value = [];
    isLoadingInstalledPlugins.value = true;
    try {
        installedPlugins.value = await pluginsStore.fetchInstalledPlugins();
    } catch (error) {
        console.error('Error refreshing installed plugins:', error);
        toast('Failed to refresh installed plugins!', {
            type: 'error',
            position: 'top-right',
            autoClose: 3000,
        });
    }
    isLoadingInstalledPlugins.value = false;
};

onMounted(() => {
    refreshInstalledPlugins().catch((error) => {
        console.error('Error fetching installed plugins:', error);
        toast('Failed to fetch installed plugins!', {
            type: 'error',
            position: 'top-right',
            autoClose: 3000,
        });
    });
    refreshAvailablePlugins().catch((error) => {
        console.error('Error fetching available plugins:', error);
        toast('Failed to fetch available plugins!', {
            type: 'error',
            position: 'top-right',
            autoClose: 3000,
        });
    });
});
</script>

<template>
    <div class="container mx-auto p-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div>
                <h2 class="text-2xl font-bold mb-4">Installed Plugins</h2>
                <div class="mt-18" :class="{'opacity-50': isLoadingInstalledPlugins}">
                    <div v-if="isLoadingInstalledPlugins" class="flex items-center justify-center h-48">
                        <span class="loading loading-spinner loading-lg"></span>
                    </div>
                    <div v-if="!isLoadingInstalledPlugins && installedPlugins.length === 0"
                         class="flex items-center justify-center h-48">
                        <span class="text-gray-500">No plugins installed.</span>
                    </div>
                    <div v-if="!isLoadingInstalledPlugins && installedPlugins.length > 0">
                        <PluginCard :plugin="plugin" v-for="plugin in installedPlugins" :key="plugin.name"/>
                    </div>
                </div>
            </div>
            <div>
                <div class="flex items-center justify-between">
                    <h2 class="text-2xl font-bold mb-4">
                        Available Plugins
                    </h2>
                    <div class="flex space-x-2">
                        <div class="tooltip"
                             data-tip="Refresh available plugins">
                            <span class="sr-only">Refresh</span>
                            <button class="btn btn-xs btn-success"
                                    :disabled="isLoadingAvailablePlugins"
                                    @click="refreshAvailablePlugins"
                            >
                                <ArrowPathIcon/>
                            </button>
                        </div>
                        <div class="tooltip"
                             data-tip="Plugin Manager Settings">
                            <span class="sr-only">Settings</span>
                            <router-link
                                class="btn btn-xs btn-primary"
                                :to="{ name: 'PluginManagerSettings' }"
                            >
                                <CogIcon/>
                            </router-link>
                        </div>
                    </div>
                </div>
                <div class="mb-4">
                    <label
                        v-if="!isLoadingAvailablePlugins"
                        class="input flex-grow rounded-none border-none focus:outline-none focus:ring-0 focus-within:outline-none focus-within:ring-0 focus-within:border-none w-full">
                        <MagnifyingGlassIcon class="h-[1em] opacity-50"/>
                        <input
                            v-model="searchQuery"
                            type="search"
                            placeholder="Search"/>
                    </label>
                </div>
                <div :class="{'opacity-50': isLoadingAvailablePlugins}">
                    <div v-if="isLoadingAvailablePlugins" class="flex items-center justify-center h-48">
                        <span class="loading loading-spinner loading-lg"></span>
                    </div>
                    <div v-if="!isLoadingAvailablePlugins && availablePlugins.length === 0"
                         class="flex items-center justify-center h-48">
                        <span class="text-gray-500">No plugins available.</span>
                    </div>
                    <div v-if="!isLoadingAvailablePlugins && availablePlugins.length > 0">
                        <PluginRepositoryCard v-for="plugin in availablePlugins" :key="plugin.id" :plugin="plugin"/>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">

</style>
