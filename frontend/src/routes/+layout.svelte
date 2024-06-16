<script>
    import '../app.css';
    import { ModeWatcher } from 'mode-watcher';

    import Sidebar from '@/components/Sidebar.svelte';
</script>

<ModeWatcher />
<div class="flex min-h-screen">
    <Sidebar />
    <slot />
</div>
