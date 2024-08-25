<script lang="ts">
    import Rule from './Rule.svelte';

    let searchQuery = '';

    function handleFilter() {
        console.log('Filter button clicked');
    }

    function handleSort() {
        console.log('Sort by button clicked');
    }

    function handleAddRule() {
        console.log('Add new rule button clicked');
    }
</script>

<div class="card flex justify-between items-center p-4 rounded-lg">
    <!-- Search Bar -->
    <input
        type="text"
        bind:value={searchQuery}
        placeholder="search bar"
        class="p-2 rounded-lg bg-input text-text w-full max-w-md"
    />

    <!-- Buttons -->
    <div class="flex gap-4 ml-4">
        <button class="btn-primary" on:click={handleFilter}>filter</button>
        <button class="btn-primary" on:click={handleSort}>sort by</button>
        <button class="btn-primary" on:click={handleAddRule}>add new rule</button>
    </div>
</div>

<!-- Rules -->
<div class="flex flex-wrap overflow-y-auto p-8">
    <Rule />
    <Rule />
    <Rule />
</div>
