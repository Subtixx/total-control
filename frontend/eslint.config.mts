import js from "@eslint/js";

import pluginVue from "eslint-plugin-vue";
import {
    defineConfigWithVueTs,
    vueTsConfigs,
} from "@vue/eslint-config-typescript";
//import prettierConfig from "@vue/eslint-config-prettier";

export default defineConfigWithVueTs(
    pluginVue.configs["flat/essential"],
    js.configs.recommended,
    js.configs.all,
    vueTsConfigs.recommended,
    //prettierConfig,
    {
        rules: {
            "capitalized-comments": "off",
            "sort-keys": "off",
            "sort-imports": "off",
            "id-length": "off",
            "no-useless-assignment": "off",
            "no-magic-numbers": "off",
            "no-empty-function": "off",
            "no-ternary": "off",
            "no-negated-condition": "off",
            "no-plusplus": "off",
            "func-style": "off",
            "@typescript-eslint/no-unused-expressions": "off",
        },
    },
);
