<script lang="ts">
    import type { SvelteMap } from 'svelte/reactivity';

    type Props = {
        data: SvelteMap<Date, number>;
    };

    let { data }: Props = $props();
</script>

{data.values().reduce((acc, c) => acc + c, 0)}
