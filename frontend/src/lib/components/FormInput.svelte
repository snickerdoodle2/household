<script lang="ts">
import { Label } from '$lib/components/ui/label';
import { Input } from '$lib/components/ui/input';
import type { HTMLInputTypeAttribute } from 'svelte/elements';

type Props = {
    name: string;
    label: string;
    type: HTMLInputTypeAttribute;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    value: any;
    errors: Record<string, string>;
    disabled?: boolean;
};

let {
    name,
    label,
    type,
    value = $bindable(),
    errors,
    disabled = false,
}: Props = $props();
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
    type={type}
    name={name}
    bind:value={value}
    required
    disabled={disabled}
    class={errors[name] ? 'border-2 border-red-600' : ''}
/>
