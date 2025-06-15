<script setup lang="ts">
import TrashIcon from '~icons/heroicons-solid/trash.svg';
import ArrowPathIcon from '~icons/heroicons/arrow-path.svg';
import CloudDownloadIcon from '~icons/heroicons-solid/cloud-arrow-down.svg';
import MagnifyingGlassIcon from '~icons/heroicons-solid/magnifying-glass.svg';

import {computed, onMounted, type Ref, ref} from 'vue';
import {type Plugin, usePluginsStore} from '@stores/plugins';

const pluginsStore = usePluginsStore();

const searchQuery = ref('');
const isLoading = ref(true);
const installedPlugins:Ref<Array<Plugin>> = ref([]);
const availablePlugins: Ref<Array<Plugin>> = ref([]);

onMounted(async () => {
    isLoading.value = true;
    installedPlugins.value = await pluginsStore.fetchInstalledPlugins();
    //availablePlugins.value = await pluginsStore.fetchAvailablePlugins();
    isLoading.value = false;
});
</script>

<template>
    <div class="container mx-auto p-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-8" :class="{'opacity-50': isLoading}">
            <div>
                <h2 class="text-2xl font-bold mb-4">Installed Plugins</h2>
                <div v-for="plugin in installedPlugins" :key="plugin.name" class="card bg-base-100 shadow-md mb-4">
                    <div class="card-body">
                        <h3 class="card-title">
                            {{ plugin.name }}
                            <span v-if="plugin.isPacked" class="badge badge-info">Packed</span>
                            <span v-if="plugin.enabled" class="badge badge-success">Enabled</span>
                            <span v-else class="badge badge-error">Disabled</span>
                        </h3>
                        <p class="text-sm text-gray-500">
                            {{ plugin.version }} - By {{ plugin.author }}
                        </p>
                        <p class="mt-2">{{ plugin.description.slice(0, 50) }}</p>
                        <div class="flex space-x-2 mt-5">
                            <div class="flex flex-grow space-x-2">
                                <div class="tooltip"
                                     :data-tip="pluginsStore.canUninstallPlugin(plugin) ? 'Uninstall' : 'Games depend on this plugin'">
                                    <button class="btn btn-xs btn-error"
                                            :disabled="!pluginsStore.canUninstallPlugin(plugin)"
                                            @click="pluginsStore.uninstallPlugin(plugin)">
                                        <TrashIcon/>
                                    </button>
                                </div>
                                <div class="tooltip"
                                     data-tip="Disable">
                                    <button class="btn btn-xs"
                                            :class="plugin.enabled ? 'btn-error' : 'btn-success'"
                                            @click="pluginsStore.togglePluginStatus(plugin)">
                                        {{ plugin.enabled ? 'Disable' : 'Enable' }}
                                    </button>
                                </div>
                            </div>
                            <div class="tooltip"
                                 :data-tip="plugin.version_available ? `Update to ${plugin.version_available}` : 'No updates available'">
                                <button class="btn btn-xs btn-info"
                                        :disabled="!pluginsStore.hasPluginUpdate(plugin)"
                                        @click="pluginsStore.updatePlugin(plugin)">
                                    <CloudDownloadIcon/>
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div>
                <h2 class="text-2xl font-bold mb-4">
                    Available Plugins
                    <button class="btn btn-xs btn-success">
                        <ArrowPathIcon/>
                    </button>
                </h2>
                <div class="mb-4">
                    <label
                        class="input flex-grow rounded-none border-none focus:outline-none focus:ring-0 focus-within:outline-none focus-within:ring-0 focus-within:border-none w-full">
                        <MagnifyingGlassIcon class="h-[1em] opacity-50"/>
                        <input
                            v-model="searchQuery"
                            type="search"
                            placeholder="Search"/>
                    </label>
                </div>
                <div v-for="plugin in availablePlugins" :key="plugin.name" class="card bg-base-100 shadow-md mb-4">
                    <div class="card-body">
                        <h3 class="card-title">{{ plugin.name }}</h3>
                        <p class="text-sm text-gray-500">By {{ plugin.author }}</p>
                        <p class="mt-2">{{ plugin.description.slice(0, 50) }}</p>

                        <div class="flex space-x-2 justify-end">
                            <button class="btn btn-xs btn-primary"
                                    @click="pluginsStore.installPlugin(plugin)">
                                <CloudDownloadIcon/>
                                Install
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">

</style>
