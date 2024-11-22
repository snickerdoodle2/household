<script lang="ts">
import '../app.css';
import { ModeWatcher } from 'mode-watcher';
type Props = {
    children?: import('svelte').Snippet;
};

let { children }: Props = $props();
</script>

<ModeWatcher />
{@render children?.()}
