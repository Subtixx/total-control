<script setup lang="ts">
import {ref} from 'vue';

import {plugins} from "@wails/go/models.ts";
import CloudDownloadIcon from "~icons/heroicons-solid/cloud-arrow-down.svg";
import MarkdownRenderer from "@components/MarkdownRenderer.vue";

type Props = {
    plugin: plugins.PluginRepositoryInfo;
};

defineProps<Props>();

const showChangelog = ref(false);
</script>

<template>
    <dialog class="modal" :class="{'modal-open': showChangelog}">
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
            <div class="card-title">
                <div class="flex-grow">
                    <a href="#" target="_blank" rel="noopener noreferrer">
                        {{ plugin.name }}
                    </a>
                </div>
                <button class="badge badge-info cursor-pointer" @click="showChangelog = true;">
                    {{ plugin.version }}
                </button>
            </div>
            <p class="text-sm text-gray-500">
                By {{ plugin.author }}
            </p>
            <p class="mt-2">{{ plugin.description.slice(0, 50) }}</p>
            <div class="flex space-x-2 mt-5">
                <div class="flex flex-grow space-x-2">
                </div>
                <div class="tooltip" data-tip="Install">
                    <button class="btn btn-xs btn-primary">
                        <CloudDownloadIcon/>
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
