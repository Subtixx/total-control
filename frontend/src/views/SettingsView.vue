<script setup lang="ts">
import ThemeSelect from "@components/ThemeSelect.vue";
</script>

<template>
    <div class="flex flex-col items-center justify-center">
        <div class="card w-1/2 bg-base-300 shadow-sm">
            <div class="card-body">
                <h2 class="card-title">Settings</h2>

                <div class="py-2">
                    <ThemeSelect/>
                </div>

                <div class="justify-end card-actions">
                    <button class="btn btn-success">Save</button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">

</style>
