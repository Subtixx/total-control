<script setup lang="ts">
import {usePluginsStore} from "@/stores/plugins";
import RepositoryEntry from "@components/RepositoryEntry.vue";
import {plugins} from "@wails/go/models.ts";
import {onMounted, ref} from 'vue';
import {toast} from "vue3-toastify";

const pluginsStore = usePluginsStore();

const autoUpdatePlugins = ref(false);
const pluginUpdateCheckEnabled = ref(false);
const updateCheckInterval = ref(1);
const updateNotificationEnabled = ref(false);
const newRepositoryUrl = ref('');

const pluginRepositories = ref<Array<plugins.PluginRepository>>([]);

onMounted(async () => {
    pluginRepositories.value = await pluginsStore.getPluginRepositories();
});

const addRepository = (url: string) => {
    if (!url) {
        toast('Please enter a valid repository URL!', {
            type: 'error',
            position: 'top-right',
            autoClose: 3000,
        });
        return;
    }

    if (!isRepositoryUrlValid(url)) {
        toast('Please enter a valid repository URL!', {
            type: 'error',
            position: 'top-right',
            autoClose: 3000,
        });
    }

    if (pluginRepositories.value.some(repo => repo.url === url)) {
        toast('Repository already exists!', {
            type: 'error',
            position: 'top-right',
            autoClose: 3000,
        });
        return;
    }

    const uuid = crypto.randomUUID();
    const newRepo: plugins.PluginRepository = plugins.PluginRepository.createFrom({
        id: uuid,
        url: url,
        plugins: [],
    });
    pluginRepositories.value.push(newRepo);
};
const removeRepository = (id: string) => {
    if (!id) {
        return;
    }

    const index = pluginRepositories.value.findIndex(repo => repo.id === id);
    if (index === -1) {
        return;
    }

    pluginRepositories.value.splice(index, 1);
};

const isRepositoryUrlValid = (url: string) => {
    try {
        new URL(url);
        return true;
    } catch (error) {
        console.error('Error parsing URL:', error);
        return false;
    }
}
</script>

<template>
    <div class="max-w-2xl mx-auto p-4">
        <div class="card bg-base-100 shadow-md mb-4">
            <div class="card-body">
                <h2 class="card-title">Plugin Manager Settings</h2>
                <p class="text-sm text-gray-500">Manage your plugins and their settings.</p>
                <div class="flex flex-col space-y-4 mt-4">
                    <div>
                        <label class="label cursor-pointer">
                            <span class="label-text">Auto-Update Plugins</span>
                            <input type="checkbox" class="toggle toggle-primary" v-model="autoUpdatePlugins"/>
                        </label>
                        <p class="text-sm text-gray-500 mt-2">Automatically update plugins when new versions are
                            available.</p>
                    </div>
                    <div>
                        <label class="label cursor-pointer">
                            <span class="label-text">Plugin Update Check</span>
                            <input type="checkbox" class="toggle toggle-primary" v-model="pluginUpdateCheckEnabled"/>
                        </label>
                        <p class="text-sm text-gray-500 mt-2">Enable or disable automatic checks for plugin updates.
                            When
                            enabled, the system will periodically check for updates to your installed plugins.</p>
                    </div>
                    <div
                        v-if="pluginUpdateCheckEnabled"
                    >
                        <label class="label cursor-pointer">
                            <span class="label-text">Plugin Update Check Interval</span>
                            <input type="number" class="input input-bordered w-full max-w-xs"
                                   v-model="updateCheckInterval"
                                   placeholder="Enter interval in hours" min="1"/>
                        </label>
                        <p class="text-sm text-gray-500 mt-2">Interval in hours for checking for plugin updates.</p>
                    </div>
                    <div v-if="pluginUpdateCheckEnabled && !autoUpdatePlugins">
                        <label class="label cursor-pointer">
                            <span class="label-text">Plugin Update Notification</span>
                            <input type="checkbox" class="toggle toggle-primary" v-model="updateNotificationEnabled"/>
                        </label>
                        <p class="text-sm text-gray-500 mt-2">Enable or disable notifications for plugin updates. When
                            enabled, you will receive notifications when updates are available for your installed
                            plugins.</p>
                    </div>
                    <div>
                        <label class="label cursor-pointer">
                            <span class="label-text">Plugin Repositories</span>
                        </label>
                        <p class="text-sm text-gray-500 mt-2">
                            Manage your plugin repositories. Add or remove
                            repositories to customize where plugins are sourced from.
                        </p>
                        <div class="flex flex-col space-y-2 mt-2">
                            <div class="flex items-center space-x-2">
                                <input type="text"
                                       class="focus:outline-none focus:ring-0 focus-within:outline-none focus-within:ring-0 flex-grow input"
                                       v-model="newRepositoryUrl"
                                       placeholder="Add new repository URL"/>
                                <button
                                    @click="addRepository(newRepositoryUrl)"
                                    :disabled="!isRepositoryUrlValid(newRepositoryUrl)"
                                    class="btn btn-info">Add Repository
                                </button>
                            </div>
                            <div class="h-36 bg-base-300 overflow-y-auto">
                                <RepositoryEntry
                                    v-for="repo in pluginRepositories"
                                    :key="repo.id"
                                    :repo="repo"
                                />
                            </div>
                        </div>
                    </div>
                    <div class="flex items-center justify-between">
                        <router-link :to="{ name: 'Plugins' }" class="btn">
                            Back to Plugins
                        </router-link>
                        <button class="btn btn-primary" @click="() => alert('Settings saved!')">Save Settings</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">

</style>
