<script lang="ts">
    import type { Snippet } from 'svelte';
    import '../app.css';
    import { ModeWatcher } from 'mode-watcher';
    import dayjs from 'dayjs';
    import relativeTime from 'dayjs/plugin/relativeTime';
    type Props = {
        children?: Snippet;
    };

    dayjs.extend(relativeTime);

    let { children }: Props = $props();
</script>

<ModeWatcher />
{@render children?.()}
