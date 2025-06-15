import {createRouter, createWebHashHistory} from "vue-router";
import HomeView from "@views/HomeView.vue";
import AddGameView from "@views/AddGameView.vue";
import GameDetailsView from "@views/GameDetailsView.vue";
import SettingsView from "@views/SettingsView.vue";
import PluginsView from "@views/PluginsView.vue";
import NotFoundView from "@views/NotFoundView.vue";

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: "/",
            name: "home",
            component: HomeView,
        },
        {
            path: "/add",
            name: "AddGame",
            component: AddGameView,
        },
        {
            path: '/games/:id',
            name: 'GameDetails',
            component: GameDetailsView,
        },
        {
            path: "/settings",
            name: "Settings",
            component: SettingsView,
        },
        {
            path: "/plugins",
            name: "Plugins",
            component: PluginsView,
        },
        {
            path: "/:pathMatch(.*)*",
            name: "NotFound",
            component: NotFoundView,
        }
    ],
});

export default router;
