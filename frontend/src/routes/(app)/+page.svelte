<script lang="ts">
    // TODO: calculate grid size
    import Sensor from '@/components/Sensor.svelte';
    import type { LayoutData } from './$types';

    export let data: LayoutData;
</script>

{#await data.sensors then sensors}
    <div class="grid h-full w-full grid-cols-5 grid-rows-6 gap-8 px-32 py-20">
        {#each sensors as sensor}
            <Sensor {sensor} />
        {/each}
    </div>
{/await}
