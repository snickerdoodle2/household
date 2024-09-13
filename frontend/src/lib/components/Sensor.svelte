<script lang="ts">
    import { type Sensor } from '@/types/sensor';
    import { DotsVertical } from 'svelte-radix';
    export let sensor: Sensor;
</script>

<div class="rounded-lg bg-accent px-4 py-2">
    <div class="flex items-center justify-between">
        <span class="text-xl">{sensor.name}</span>
        <a href={`/details/${sensor.id}`}><DotsVertical class="h-5 w-5" /></a>
    </div>
</div>
