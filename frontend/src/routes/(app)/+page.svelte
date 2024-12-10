<script lang="ts">
    // TODO: calculate grid size
    import Sensor from '@/components/sensor/Sensor.svelte';
    import type { LayoutData } from './$types';
    import { Button } from '@/components/ui/button';
    import { Eye, EyeOff, Plus } from 'lucide-svelte';

    type Props = {
        data: LayoutData;
    };

    let { data }: Props = $props();
    const { sensors } = data;
    let disableShowButton = $derived(
        sensors.filter((sensor) => sensor.hidden).length === 0
    );

    let showHidden = $state(false);
    let shownSensors = $derived.by(() => {
        if (showHidden) {
            return sensors;
        }
        return sensors.filter((sensor) => !sensor.hidden);
    });
</script>

<div class="flex flex-col h-full w-full px-4 py-4 gap-6">
    <div class="flex justify-end gap-2">
        <Button
            variant="outline"
            size="icon"
            disabled={disableShowButton}
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
    <div
        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-8"
    >
        {#each shownSensors as sensor (sensor.id)}
            <Sensor {sensor} />
        {/each}
    </div>
</div>
