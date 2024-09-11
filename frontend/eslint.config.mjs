// @ts-check

import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';
import svelteParser from 'svelte-eslint-parser';
import eslintPluginSvelte from 'eslint-plugin-svelte';
import eslintConfigPrettier from 'eslint-config-prettier';
import globals from 'globals';

export default tseslint.config(
    eslint.configs.recommended,
    ...tseslint.configs.strict,
    ...tseslint.configs.stylistic,
    {
        rules: {
            '@typescript-eslint/consistent-type-definitions': ['error', 'type'],
        },
    },
    ...eslintPluginSvelte.configs['flat/recommended'],
    eslintConfigPrettier,
    ...eslintPluginSvelte.configs['flat/prettier'],
    {
        files: ['*.svelte', '**/*.svelte'],
        languageOptions: {
            parser: svelteParser,
            parserOptions: {
                parser: tseslint.parser,
            },
            globals: {
                ...globals.browser,
            },
        },
    },
    {
        ignores: ['.svelte-kit/*', 'src/lib/components/ui/*'],
    }
);
