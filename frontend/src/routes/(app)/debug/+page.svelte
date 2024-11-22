<script lang="ts">
import { authToken } from '@/auth/token';
import type { LayoutData } from '../$types';
import { Button } from '@/components/ui/button';

type Props = {
    data: LayoutData;
};

let { data }: Props = $props();
</script>

<div class="flex flex-col">
    <div>
        <code><pre>{JSON.stringify($authToken, null, 4)}</pre></code>
        {#await data.user then user}
            <code><pre>{JSON.stringify(user, null, 4)}</pre></code>
        {/await}
    </div>
    <Button href="/debug/ws">Websocket</Button>
</div>
