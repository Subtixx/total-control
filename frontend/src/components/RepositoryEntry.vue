<script setup lang="ts">
import {plugins} from "@wails/go/models.ts";
import ExclamationTriangleIcon from '~icons/heroicons-solid/exclamation-triangle.svg';
import TrashIcon from '~icons/heroicons-solid/trash.svg';

type Props = {
    repo: plugins.PluginRepository;
    canRemove?: boolean;
    error?: string;
}

const props = defineProps<Props>();

defineEmits<{
    remove: [id: string],
}>();
</script>

<template>
    <div :key="props.repo.id"
         class="flex gap-4 items-center p-2 group hover:bg-base-200">
        <div class="tooltip text-error" v-if="props.error">
            <div class="tooltip-content" v-text="props.error">
            </div>
            <div>
                <ExclamationTriangleIcon/>
            </div>
        </div>
        <span
            class="flex-grow text-sm text-gray-500 max-w-xl overflow-hidden text-ellipsis">
                {{ props.repo.url }}
            </span>
        <button
            v-if="props.canRemove"
            @click="$emit('remove', props.repo.id)"
            class="btn btn-error btn-sm opacity-0 group-hover:opacity-100 transition-opacity">
            <TrashIcon/>
        </button>
    </div>
</template>
