<script lang="ts">
    import type { PageData } from './$types';
    export let data: PageData;
</script>

{#await data.rule then rule}
    <p>{rule.name}</p>
{/await}
