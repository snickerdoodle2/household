<script lang="ts">
    import { onMount } from 'svelte';
    import type { NewRule } from '$lib/types/rule';
    import type { PageData } from './$types';

    export let data: PageData;

    let loading = true;
    let sensors: { label: string; value: string }[] = [];
    let selectedSensor: { label: string; value: string };
    let rule: NewRule = {
        name: '',
        description: '',
        on_valid: {
            to: '',
            payload: {},
        },
        // @ts-expect-error nah dont wanna do this
        internal: {},
    };
    let errors: Record<string, string> = {};
    let internal = '';
    let payload = '';

    $: if (!loading && selectedSensor) {
        rule.on_valid.to = selectedSensor.value;
    }

    $: if (!loading) {
        try {
            rule.on_valid.payload = JSON.parse(payload);
            delete errors['payload'];
            errors = errors;
        } catch {
            errors['payload'] = 'Not a valid JSON';
        }
    }

    onMount(async () => {
        sensors = (await data.sensors).map((e) => ({
            value: e.id,
            label: e.name,
        }));
        loading = false;
    });
</script>

{#if loading}
    <p>loading</p>
{:else}{/if}
