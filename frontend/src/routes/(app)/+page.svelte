<script lang="ts">
    // TODO: calculate grid size
    import Sensor from '@/components/sensor/Sensor.svelte';
    import type { LayoutData } from './$types';
    import { goto, preloadData, pushState } from '$app/navigation';
    import * as Dialog from '$lib/components/ui/dialog';
    import { page } from '$app/stores';
    import DetailsPage from './details/[id]/+page.svelte';
    import { Plus } from 'svelte-radix';
    import { Button } from '@/components/ui/button';

    let modalOpen = false;

    const handleDetails = async (e: Event) => {
        const { href } = e.currentTarget as HTMLAnchorElement;
        const res = await preloadData(href);
        if (res.type === 'loaded' && res.status == 200) {
            const data = await res.data.sensor;
            modalOpen = true;
            pushState(href, {
                selected: data,
            });
        } else {
            goto(href);
        }
    };

    export let data: LayoutData;
</script>

<div class="flex h-full items-start gap-4 md:py-20">
    {#await data.sensors then sensors}
        <div
            class="grid flex-1 grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5"
        >
            {#each sensors as sensor}
                <Sensor {sensor} on:click={handleDetails} />
            {/each}
        </div>
    {/await}
    <Button variant="outline" size="icon" href="/create">
        <Plus />
    </Button>
</div>

<Dialog.Root
    bind:open={modalOpen}
    onOpenChange={(opened) => {
        if (!opened) {
            history.back();
        }
    }}
>
    <Dialog.Portal>
        <Dialog.Overlay />
        <Dialog.Content
            class="flex max-w-[600px] items-center justify-center px-8 py-4"
        >
            {#if $page.state.selected}
                <DetailsPage data={{ sensor: $page.state.selected }} />
            {/if}
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>
