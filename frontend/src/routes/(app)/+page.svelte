<script lang="ts">
    // TODO: calculate grid size
    import Sensor from '@/components/sensor/Sensor.svelte';
    import type { LayoutData } from './$types';
    import { goto } from '$app/navigation';
    import { Plus } from 'svelte-radix';
    import { Button } from '@/components/ui/button';

    const handleDetails = async (e: Event) => {
        const { href } = e.currentTarget as HTMLAnchorElement;
        goto(href);
    };

    const handleCreate = () => {
        goto(window.location.href + 'create');
    };

    type Props = {
        data: LayoutData;
    };

    let { data }: Props = $props();
</script>

<div class="flex h-full items-start gap-4 md:py-20">
    {#await data.sensors then sensors}
        <div
            class="grid flex-1 grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5"
        >
            {#each sensors as sensor}
                <Sensor {sensor} on:click={handleDetails} />
            {/each}
            <div class="flex items-center justify-center">
                <Button variant="outline" size="icon" on:click={handleCreate}>
                    <Plus />
                </Button>
            </div>
        </div>
    {/await}
</div>
