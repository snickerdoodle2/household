<script lang="ts">
    import { Label } from '$lib/components/ui/label';
    import { Input } from '$lib/components/ui/input';
    import type { HTMLInputTypeAttribute } from 'svelte/elements';

    export let name: string;
    export let label: string;
    export let type: HTMLInputTypeAttribute;
    export let value: string;
    export let errors: Record<string, string>;
    export let disabled = false;
</script>

<Label
    for={name}
    class="flex items-center justify-between text-base font-semibold"
    >{label}
    {#if errors[name]}
        <span class="text-sm font-normal italic text-red-400"
            >{errors[name]}</span
        >
    {/if}
</Label>
<Input
    {type}
    {name}
    bind:value
    required
    {disabled}
    class={errors[name] ? 'border-2 border-red-600' : ''}
/>
