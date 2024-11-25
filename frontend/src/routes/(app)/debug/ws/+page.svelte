<script lang="ts">
    import * as Select from '@/components/ui/select';
    import { SensorWebsocket } from '@/helpers/socket.svelte';
    import type { PageData } from './$types';
    import { Button } from '@/components/ui/button';
    import { untrack } from 'svelte';
    import { X, History } from 'lucide-svelte';
    let ws = new SensorWebsocket();

    let selected = $state('');

    let updates = $state(0);

    $effect(() => {
        console.log(ws.data.forEach((e) => e.forEach((x) => x + 1)));
        untrack(() => (updates += 1));
    });

    type Props = {
        data: PageData;
    };

    let { data }: Props = $props();
</script>

<div class="flex w-[1024px] flex-col items-center justify-center gap-2">
    {#await data.sensors then sensors}
        <div class="flex w-96 flex-col gap-2">
            <Select.Root
                items={sensors.map((e) => ({ value: e.id, label: e.name }))}
                onSelectedChange={(e) => {
                    selected = e?.value ?? '';
                }}
            >
                <Select.Trigger class="max-w-96">
                    <Select.Value placeholder="Select a sensor..." />
                </Select.Trigger>
                <Select.Content>
                    {#each sensors as sensor}
                        <Select.Item value={sensor.id} label={sensor.name}
                            >{sensor.name}</Select.Item
                        >
                    {/each}
                </Select.Content>
            </Select.Root>
            <Button
                disabled={selected.length <= 0}
                onclick={() => {
                    ws.subscribe(selected);
                }}>Subscribe</Button
            >
        </div>
        <div>
            <ol>
                {#each ws.data.entries() as [key, value]}
                    <li class="flex gap-3">
                        <div>
                            <code
                                >{sensors.find((e) => e.id === key)?.name}</code
                            >
                            - {(
                                value
                                    .values()
                                    .reduce((acc, cur) => acc + cur, 0) /
                                value.values().reduce((acc) => acc + 1, 0)
                            ).toFixed(2)}
                            ({value
                                .values()
                                .reduce((acc, cur) => acc + cur, 0)
                                .toFixed()})
                        </div>
                        <button onclick={() => ws.requestSince(key, '24h')}
                            ><History /></button
                        >
                        <button onclick={() => ws.unsubscribe(key)}
                            ><X /></button
                        >
                    </li>
                {/each}
            </ol>
            <p class="text-sm text-stone-400">Updated {updates} times!</p>
        </div>
    {/await}
</div>
