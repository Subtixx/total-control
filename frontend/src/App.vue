<script setup lang="ts">
import Sidebar from "./components/Sidebar.vue";
import NavBar from "@components/NavBar.vue";
</script>

<template>
    <div class="flex w-screen h-screen">
        <sidebar/>

        <div class="flex flex-col flex-grow bg-base-200">
            <NavBar/>
            <div class="flex-grow p-6 overflow-auto border-l border-base-300">
                <router-view/>
            </div>
        </div>
    </div>
</template>

<style lang="scss"></style>
