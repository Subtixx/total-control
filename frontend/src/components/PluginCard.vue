<script setup lang="ts">
import {computed, ref} from 'vue';
import TrashIcon from '~icons/heroicons-solid/trash.svg';
import CloudDownloadIcon from '~icons/heroicons-solid/cloud-arrow-down.svg';
import {type Plugin, usePluginsStore} from '@stores/plugins';
import MarkdownRenderer from "@components/MarkdownRenderer.vue";

type Props = {
    plugin: Plugin;
};

defineProps<Props>();
const pluginsStore = usePluginsStore();

const showChangelog = ref(false);
</script>

<template>
    <dialog class="modal" :class="{'modal-open': showChangelog}" v-if="plugin.changelog">
        <div class="modal-box">
            <h3 class="text-lg font-bold">Changelog</h3>
            <p class="py-4">
                <MarkdownRenderer :content="plugin.changelog"/>
            </p>
            <div class="modal-action">
                <button class="btn" @click="showChangelog=false">Close</button>
            </div>
        </div>
    </dialog>
    <div class="card bg-base-100 shadow-md mb-4">
        <div class="card-body">
            <div class="card-title flex items-center space-x-2">
                <a href="#" class="flex-grow" target="_blank" rel="noopener noreferrer">
                    {{ plugin.name }}
                </a>
                <span v-if="plugin.isPacked" class="badge badge-info">Packed</span>
                <span v-if="plugin.enabled" class="badge badge-success">Enabled</span>
                <span v-else class="badge badge-error">Disabled</span>
                <button class="badge"
                        @click="showChangelog=true"
                        :class="{
                            'cursor-pointer': plugin.changelog,
                    'badge-info': !plugin.version_available,
                    'badge-success': plugin.version_available,
                    'badge-error': plugin.version_available && plugin.version_available !== plugin.version
                }"
                >
                    {{ plugin.version }}
                </button>
            </div>
            <p class="text-sm text-gray-500">
                By {{ plugin.author }}
            </p>
            <p class="mt-2">{{ plugin.description }}</p>
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
</template>
