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
    class="mt-auto h-11 w-11"
>
    <Sun
        class="h-5 w-5 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0"
    />
    <Moon
        class="absolute h-5 w-5 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100"
    />
    <span class="sr-only">Toggle theme</span>
</Button>
