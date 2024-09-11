/**
 * @see https://prettier.io/docs/en/configuration.html
 * @type {import("prettier").Config}
 */
const config = {
    trailingComma: 'es5',
    tabWidth: 4,
    semi: true,
    singleQuote: true,
    plugins: ['prettier-plugin-svelte', 'prettier-plugin-tailwindcss'],
    tailwindConfig: './tailwind.config.js',
    overrides: [
        {
            files: '*.svelte',
            options: {
                parser: 'svelte',
            },
        },
    ],
};

export default config;
