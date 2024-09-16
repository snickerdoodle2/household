<script lang="ts">
    // TODO: calculate grid size
    import Sensor from '@/components/sensor/Sensor.svelte';
    import type { LayoutData } from './$types';
    import { goto, preloadData, pushState } from '$app/navigation';
    import * as Dialog from '$lib/components/ui/dialog';
    import { page } from '$app/stores';
    import DetailsPage from './details/[id]/+page.svelte';

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

{#await data.sensors then sensors}
    <div class="grid h-full w-full grid-cols-5 grid-rows-6 gap-8 px-32 py-20">
        {#each sensors as sensor}
            <Sensor {sensor} on:click={handleDetails} />
        {/each}
    </div>
{/await}

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
