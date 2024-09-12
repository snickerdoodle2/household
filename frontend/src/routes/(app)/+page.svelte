<script lang="ts">
    import type { LayoutData } from './$types';

    export let data: LayoutData;
</script>

{#await data.sensors then sensors}
    {sensors.length}
{/await}
