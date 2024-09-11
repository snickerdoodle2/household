// @ts-check

import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';

export default tseslint.config(
    eslint.configs.recommended,
    ...tseslint.configs.strict,
    ...tseslint.configs.stylistic,
    {
        ignores: [".svelte-kit/*"],
    },
    {
        rules: {
            "@typescript-eslint/consistent-type-definitions": [
                "error",
                "type"
            ],
        }
    }
);
