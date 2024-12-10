<script lang="ts">
    // TODO: calculate grid size
    import type { LayoutData } from './$types';
    import { Button } from '@/components/ui/button';
    import { Eye, EyeOff, Plus } from 'lucide-svelte';

    type Props = {
        data: LayoutData;
    };

    let { data }: Props = $props();
    const { sensors } = data;
    let showHidden = $state(false);
    let shownSensors = $derived.by(() => {
        if (showHidden) {
            return sensors;
        }
        return sensors.filter((sensor) => !sensor.hidden);
    });
</script>

<div class="flex flex-col h-full w-full px-4 py-4">
    <div class="flex justify-end gap-2">
        <Button
            variant="outline"
            size="icon"
            onclick={() => {
                showHidden = !showHidden;
            }}
        >
            {#if showHidden}
                <EyeOff class="w-6 h-6" />
            {:else}
                <Eye class="w-6 h-6" />
            {/if}
        </Button>
        <Button variant="outline" size="icon" href="/create"
            ><Plus class="w-6 h-6" /></Button
        >
    </div>
    {#each shownSensors as sensor}
        {sensor.name}
    {/each}
</div>
