<script lang="ts">
    import type { Snippet } from 'svelte';
    import '../app.css';
    import { ModeWatcher } from 'mode-watcher';
    type Props = {
        children?: Snippet;
    };

    let { children }: Props = $props();
</script>

<ModeWatcher />
{@render children?.()}
