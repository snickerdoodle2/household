<script lang="ts">
    import { AppWebsocket } from '@/helpers/socket.svelte';

    const ws = new AppWebsocket();
</script>

<ol>
    {#each ws.notifications as notification}
        <li>
            ({notification.level}) {notification.title} - {notification.description}
        </li>
    {/each}
</ol>
