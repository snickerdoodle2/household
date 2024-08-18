<script lang="ts">
    import { Button } from '$lib/components/ui/button/index.js';
    import { toggleMode } from 'mode-watcher';
    import { Sun } from 'radix-icons-svelte';
    import { Moon } from 'radix-icons-svelte';
</script>

<Button
    on:click={toggleMode}
    variant="outline"
    size="icon"
    class="w-11 h-11 mt-auto"
>
    <Sun
        class="rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0 w-5 h-5"
    />
    <Moon
        class="absolute rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100 w-5 h-5"
    />
    <span class="sr-only">Toggle theme</span>
</Button>
