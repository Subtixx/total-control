<script setup lang="ts">
import {onMounted, ref, type Ref, watch} from "vue";
import {useTheme} from "@/themePlugin.ts";

const activeTheme: Ref<string | null> = ref(null);

const theme = useTheme();
onMounted(() => {
    activeTheme.value = theme.theme.value;
});

watch(activeTheme, (newTheme) => {
    if (newTheme === 'system') {
        theme.setSystemTheme();
        return
    }
    if (!newTheme) {
        return;
    }

    theme.setTheme(newTheme);
});
</script>

<template>
    <select
        class="select w-full"
        v-model="activeTheme">
        <option
            v-for="theme in theme.availableThemes"
            :selected="activeTheme === theme"
            :key="theme"
            :value="theme">
            {{ theme }}
        </option>
        <option value="system">System Default</option>
    </select>
</template>
