<script setup lang="ts">
import {computed} from "vue";
import hljs from "highlight.js";
import {Marked} from "marked";
import {markedHighlight} from "marked-highlight";

const marked = new Marked(
    markedHighlight({
        emptyLangClass: 'hljs',
        langPrefix: 'hljs language-',
        highlight(code: string, lang: string, info: string) {
            console.log(info)
            const language = hljs.getLanguage(lang) ? lang : 'plaintext';
            return hljs.highlight(code, {language}).value;
        }
    })
);

type Props = {
    content: string;
};
const props = defineProps<Props>();

const htmlContent = computed(() => {
    return marked.parse(props.content);
})
</script>

<template>
    <div v-html="htmlContent"></div>
</template>
