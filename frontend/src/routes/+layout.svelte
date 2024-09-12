<script lang="ts">
    import '../app.css';
    import { ModeWatcher } from 'mode-watcher';

    import Sidebar from '@/components/Sidebar.svelte';
    import { onMount } from 'svelte';
    import { authToken } from '@/auth/token';
    onMount(() => {
        const token = localStorage.getItem('authToken');
        if (token) {
            authToken.set(JSON.parse(token));
        }
    });
</script>

<ModeWatcher />
<div class="flex min-h-screen">
    <Sidebar />
    <main class="mx-4 my-2 flex flex-1 items-center justify-center pl-20">
        <slot />
    </main>
</div>
