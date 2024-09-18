<script lang="ts">
    import type { PageData } from './$types';
    export let data: PageData;
</script>

{#await data.rules then rules}
    {#each rules as rule}
        <p>{rule.name}</p>
    {/each}
{/await}
