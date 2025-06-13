import {createRouter, createWebHashHistory} from "vue-router";
import HomeView from "@views/HomeView.vue";
import AboutView from "@views/AboutView.vue";
import AddGameView from "@views/AddGameView.vue";
import GameDetailsView from "@views/GameDetailsView.vue";

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
            path: "/about",
            name: "about",
            component: AboutView,
        },
    ],
});

export default router;
