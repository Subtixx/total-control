import {defineStore} from 'pinia';
import {ref} from "vue";
import {GetAvailablePlugins, GetInstalledPlugins, GetPluginRepositories} from "@wails/go/main/App";
import {plugins} from "@wails/go/models.ts";

export interface Plugin {
    id: string;
    name: string;
    description: string;
    author: string;

    version: string;
    version_available?: string;

    enabled?: boolean;
    dir?: string; // Directory where the plugin is located
    isPacked?: boolean; // Whether the plugin is packed or not

    changelog?: string;
}

interface PluginsState {
    plugins: Plugin[];
    repositories: plugins.PluginRepository[];
}

const demoInstalledPlugins = ref<Array<Plugin>>([
    {
        id: 'com.github.subtixx.factorio',
        name: 'Factorio',
        description: 'Adds support for Factorio mods',
        version: '1.0.0',
        version_available: '1.0.1',
        author: 'IT-Hock',
        enabled: true,
    },
]);

const demoAvailablePlugins = ref<Array<plugins.PluginRepositoryInfo>>([
    {
        id: 'com.github.subtixx.unity',
        name: 'Unity',
        description: 'Adds support for generic Unity games',
        version: '1.0.0',
        changelog: "# 1.0.0\nInitial Release",
        author: 'IT-Hock',
    },
    {
        id: 'com.github.subtixx.unreal',
        name: 'Unreal Engine',
        description: 'Adds support for generic Unreal Engine games',
        version: '1.0.0',
        changelog: "# 1.0.0\nInitial Release",
        author: 'IT-Hock'
    },
    {
        id: 'com.github.subtixx.satisfactory',
        name: 'Satisfactory',
        description: 'Adds support for Satisfactory mods',
        version: '1.0.0',
        changelog: "# 1.0.0\nInitial Release",
        author: 'IT-Hock'
    },
    {
        id: 'com.github.subtixx.terraria',
        name: 'Terraria',
        description: 'Adds support for Terraria mods',
        version: '1.0.0',
        changelog: "# 1.0.0\nInitial Release",
        author: 'IT-Hock'
    },
    {
        id: 'com.github.subtixx.starbound',
        name: 'Starbound',
        description: 'Adds support for Starbound mods',
        version: '1.0.0',
        changelog: "# 1.0.0\nInitial Release",
        author: 'IT-Hock'
    },
    {
        id: 'com.github.subtixx.spaceengineers',
        name: 'Space Engineers',
        description: 'Adds support for Space Engineers mods',
        version: '1.0.0',
        changelog: "# 1.0.0\nInitial Release",
        author: 'IT-Hock'
    },
]);

export const usePluginsStore = defineStore('plugins', {
    state: (): PluginsState => ({
        plugins: [],
        repositories: [],
    }),
    actions: {
        async fetchInstalledPlugins(): Promise<Array<Plugin>> {
            if (window.go === undefined) {
                const randomTimeout = Math.floor(Math.random() * 2000) + 500;
                return new Promise(resolve => {
                    setTimeout(() => {
                        resolve(demoInstalledPlugins.value);
                    }, randomTimeout);
                });
            }
            this.plugins = []; // Clear existing plugins

            const installedPlugins = await GetInstalledPlugins()
            for (const plugin of installedPlugins) {
                this.plugins.push({
                    // @ts-expect-error Auto generated type, might not match exactly
                    id: plugin.id as string,
                    name: plugin.name,
                    author: plugin.author,
                    description: "",
                    version: plugin.version,
                    dir: plugin.PluginDir,
                    enabled: plugin.Enabled,
                    isPacked: plugin.IsPacked,
                } as Plugin);
            }
            return this.plugins;
        },
        async fetchAvailablePlugins(): Promise<Array<plugins.PluginRepositoryInfo>> {
            if (window.go === undefined) {
                const randomTimeout = Math.floor(Math.random() * 2000) + 500;
                return new Promise(resolve => {
                    setTimeout(() => {
                        resolve(demoAvailablePlugins.value);
                    }, randomTimeout);
                });
            }

            return await GetAvailablePlugins();
        },
        async getPluginRepositories(): Promise<Array<plugins.PluginRepository>> {
            if (window.go === undefined) {
                const randomTimeout = Math.floor(Math.random() * 2000) + 500;
                return new Promise(resolve => {
                    setTimeout(() => {
                        resolve([
                            {
                                id: 'demo',
                                url: 'https://github.com/subtixx/total-control/tree/main/plugins/repository.json',
                                plugins: demoAvailablePlugins.value
                            } as plugins.PluginRepository
                        ]);
                    }, randomTimeout);
                });
            }

            if (this.repositories.length > 0) {
                return this.repositories;
            }

            const repositories = await GetPluginRepositories();
            this.repositories = repositories.map(repo => ({
                id: repo.id,
                url: repo.url,
                plugins: repo.plugins,
            }) as plugins.PluginRepository);
            return this.repositories;
        },
        addPlugin(plugin: Plugin) {
            this.plugins.push(plugin);
        },
        removePlugin(id: string) {
            this.plugins = this.plugins.filter(plugin => plugin.id !== id);
        },
        setPlugins(plugins: Plugin[]) {
            this.plugins = plugins;
        },
        isPluginInstalled(id: string): boolean {
            return this.plugins.find(plugin => plugin.id === id)?.enabled !== undefined;
        },
        isPluginEnabled(id: string): boolean {
            const plugin = this.plugins.find(plugin => plugin.id === id);
            return plugin ? plugin.enabled ?? false : false;
        },
        getPlugin(id: string): Plugin | undefined {
            return this.plugins.find(plugin => plugin.id === id);
        },
        installPlugin(plugin: Plugin) {
            if (!this.isPluginInstalled(plugin.id)) {
                const foundPlugin = this.plugins.find(p => p.id === plugin.id);
                if (foundPlugin) {
                    foundPlugin.enabled = true;
                } else {
                    this.plugins.push({...plugin, enabled: true});
                }
            }
        },
        hasPluginUpdate(plugin: Plugin): boolean {
            const existingPlugin = this.plugins.find(p => p.id === plugin.id);
            return existingPlugin ? existingPlugin.version_available !== undefined : false;
        },
        canUninstallPlugin(plugin: Plugin): boolean {
            // TODO: Check if the user has any games that depend on this plugin
            console.log(`Checking if plugin ${plugin.id} can be uninstalled...`);
            return plugin.name.toLowerCase() !== 'factorio'; // Example condition, adjust as needed
        },
        uninstallPlugin(plugin: Plugin) {
            const index = this.plugins.findIndex(p => p.id === plugin.id);
            this.plugins[index].enabled = undefined;
        },
        updatePlugin(plugin: Plugin) {
            const index = this.plugins.findIndex(p => p.id === plugin.id);
            this.plugins[index].version = plugin.version_available || plugin.version;
            this.plugins[index].version_available = undefined;
        },
        togglePluginStatus(plugin: Plugin) {
            const index = this.plugins.findIndex(p => p.id === plugin.id);
            if (index !== -1) {
                this.plugins[index].enabled = !this.plugins[index].enabled;
            } else {
                console.warn(`Plugin with id ${plugin.id} not found.`);
            }
        }
    },
    getters: {
        enabledPlugins: (state) => state.plugins.filter(plugin => plugin.enabled),
        disabledPlugins: (state) => state.plugins.filter(plugin => !plugin.enabled),
        allPlugins: (state) => state.plugins,
        installedPlugins: (state) => state.plugins.filter(plugin => plugin.enabled !== undefined),
        availablePlugins: (state) => state.plugins.filter(plugin => plugin.enabled === undefined),
        getPluginById: (state) => (id: string) => state.plugins.find(plugin => plugin.id === id),
    },
});
