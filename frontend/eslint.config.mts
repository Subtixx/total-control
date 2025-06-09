import js from '@eslint/js';

import pluginVue from 'eslint-plugin-vue'
import {
  defineConfigWithVueTs,
  vueTsConfigs,
} from '@vue/eslint-config-typescript'
import prettierConfig from "@vue/eslint-config-prettier";

export default defineConfigWithVueTs(
  pluginVue.configs['flat/essential'],
  js.configs.recommended,
  js.configs.all,
  prettierConfig,
  vueTsConfigs.recommended,
)
