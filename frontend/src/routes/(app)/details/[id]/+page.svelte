<script lang="ts">
    import type { PageData } from './$types';

    export let data: PageData;
</script>

{#await data.sensor then sensor}
    <pre>{JSON.stringify(sensor, null, '    ')}</pre>
{/await}
