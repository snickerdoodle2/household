import { fontFamily } from 'tailwindcss/defaultTheme';

/** @type {import('tailwindcss').Config} */
const config = {
    darkMode: ['class'],
    content: [
        './src/**/*.{html,js,svelte,ts}',
        './node_modules/flowbite/**/*.js',
        './node_modules/flowbite-svelte/**/*.js',
    ],
    plugins: [
        require('flowbite/plugin'),
        function ({ addComponents }) {
            addComponents({
                // Primary Button
                '.btn-primary': {
                    backgroundColor: 'hsl(var(--primary))',
                    color: 'hsl(var(--primary-foreground))',
                    borderRadius: 'var(--radius)',
                    padding: '0.5rem 1rem',
                    textAlign: 'center',
                    display: 'inline-block',
                    fontWeight: '600',
                    transition: 'background-color 0.3s ease',
                    '&:hover': {
                        backgroundColor: 'hsl(var(--primary-hover))',
                    },
                },

                // Secondary Button
                '.btn-secondary': {
                    backgroundColor: 'hsl(var(--secondary))',
                    color: 'hsl(var(--secondary-foreground))',
                    borderRadius: 'var(--radius)',
                    padding: '0.5rem 1rem',
                    textAlign: 'center',
                    display: 'inline-block',
                    fontWeight: '600',
                    transition: 'background-color 0.3s ease',
                    '&:hover': {
                        backgroundColor: 'hsl(var(--secondary-hover))',
                    },
                },

                // Card
                '.card': {
                    backgroundColor: 'hsl(var(--card))',
                    color: 'hsl(var(--card-foreground))',
                    borderRadius: 'var(--radius)',
                    padding: '1rem',
                    boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
                },

                // Muted Card
                '.card-muted': {
                    backgroundColor: 'hsl(var(--muted))',
                    color: 'hsl(var(--muted-foreground))',
                    borderRadius: 'var(--radius)',
                    padding: '1rem',
                    boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
                },

                // Accent Card
                '.card-accent': {
                    backgroundColor: 'hsl(var(--accent))',
                    color: 'hsl(var(--accent-foreground))',
                    borderRadius: 'var(--radius)',
                    padding: '1rem',
                    boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
                },

                // Destructive Button
                '.btn-destructive': {
                    backgroundColor: 'hsl(var(--destructive))',
                    color: 'hsl(var(--destructive-foreground))',
                    borderRadius: 'var(--radius)',
                    padding: '0.5rem 1rem',
                    textAlign: 'center',
                    display: 'inline-block',
                    fontWeight: '600',
                    transition: 'background-color 0.3s ease',
                    '&:hover': {
                        backgroundColor: 'hsl(var(--destructive-hover))',
                    },
                },
                // Input Field
                '.input-field': {
                    backgroundColor: 'hsl(var(--input))',
                    color: 'hsl(var(--foreground))',
                    borderRadius: 'var(--radius)',
                    padding: '0.5rem 1rem',
                    border: '1px solid hsl(var(--border))',
                    transition: 'border-color 0.3s ease, box-shadow 0.3s ease',
                    width: '100%',
                    fontSize: '1rem',
                    lineHeight: '1.5',
                    '&:focus': {
                        borderColor: 'hsl(var(--primary))',
                        boxShadow: '0 0 0 3px hsl(var(--primary) / 0.5)',
                        outline: 'none',
                    },
                    '&:hover': {
                        borderColor: 'hsl(var(--primary-hover))',
                    },
                    '&::placeholder': {
                        color: 'hsl(var(--muted-foreground))',
                        opacity: '0.7',
                    },
                    '&[disabled]': {
                        backgroundColor: 'hsl(var(--muted))',
                        color: 'hsl(var(--muted-foreground))',
                        cursor: 'not-allowed',
                    },
                },
            });
        },
    ],
    safelist: ['dark'],
    theme: {
        container: {
            center: true,
            padding: '2rem',
            screens: {
                '2xl': '1400px',
            },
        },
        extend: {
            colors: {
                border: 'hsl(var(--border) / <alpha-value>)',
                input: 'hsl(var(--input) / <alpha-value>)',
                ring: 'hsl(var(--ring) / <alpha-value>)',
                error: 'hsl(var(--error) / <alpha-value>)',
                background: 'hsl(var(--background) / <alpha-value>)',
                foreground: 'hsl(var(--foreground) / <alpha-value>)',
                primary: {
                    DEFAULT: 'hsl(var(--primary) / <alpha-value>)',
                    hover: 'hsl(var(--primary-hover))',
                    foreground: 'hsl(var(--primary-foreground) / <alpha-value>)',
                },
                secondary: {
                    DEFAULT: 'hsl(var(--secondary) / <alpha-value>)',
                    foreground: 'hsl(var(--secondary-foreground) / <alpha-value>)',
                },
                destructive: {
                    DEFAULT: 'hsl(var(--destructive) / <alpha-value>)',
                    foreground: 'hsl(var(--destructive-foreground) / <alpha-value>)',
                },
                muted: {
                    DEFAULT: 'hsl(var(--muted) / <alpha-value>)',
                    foreground: 'hsl(var(--muted-foreground) / <alpha-value>)',
                },
                accent: {
                    DEFAULT: 'hsl(var(--accent) / <alpha-value>)',
                    foreground: 'hsl(var(--accent-foreground) / <alpha-value>)',
                },
                popover: {
                    DEFAULT: 'hsl(var(--popover) / <alpha-value>)',
                    foreground: 'hsl(var(--popover-foreground) / <alpha-value>)',
                },
                card: {
                    DEFAULT: 'hsl(var(--card) / <alpha-value>)',
                    foreground: 'hsl(var(--card-foreground) / <alpha-value>)',
                },
            },
            borderRadius: {
                lg: 'var(--radius)',
                md: 'calc(var(--radius) - 2px)',
                sm: 'calc(var(--radius) - 4px)',
            },
            fontFamily: {
                sans: [...fontFamily.sans],
            },
            height: {
                p20: '10%',
                p20: '20%',
                p30: '30%',
                p40: '40%',
                p50: '50%',
                p60: '60%',
                p70: '70%',
                p80: '80%',
                p90: '90%',
                p100: '100%',
            },
        },
    },
};

export default config;
