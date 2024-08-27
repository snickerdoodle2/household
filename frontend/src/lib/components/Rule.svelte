<script lang="ts">
    // Hardcoded rule data
    const rule = {
        name: 'Podlewanie grządek',
        type: 'Static rule',
        lastActivation: '18:04 10-12-2024',
        description: '--- no description privded ----',
    };

    function handleSeeMore() {
        console.log('See more button clicked');
    }
</script>

<div class="btn-secondary flex justify-between items-center px-20 py-10 my-4 rounded-full w-full">
    <!-- Rule Information -->
    <div class="flex flex-col space-y-1">
        <p class="font-bold text-lg">{rule.name}</p>
        <p class="text-xs">last activation: {rule.lastActivation}</p>
    </div>

    <div class="flex-1">
        <p class="text-sm">{rule.type}</p>
    </div>

    <!-- Rule Description -->
    <div class="flex-1 mx-4 text-sm">
        {rule.description}
    </div>

    <!-- See More Button -->
    <button class="btn-primary" on:click={handleSeeMore}>see more</button>
</div>

<style>
</style>
