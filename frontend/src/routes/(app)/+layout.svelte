<script lang="ts">
    import Sidebar from '@/components/Sidebar.svelte';
</script>

<div class="flex min-h-screen">
    <Sidebar />

    <main class="mx-4 my-2 flex flex-1 items-center justify-center pl-20">
        <slot />
    </main>
</div>
