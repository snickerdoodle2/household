<script>
import "../app.css";
import { ModeWatcher } from "mode-watcher";

import Sidebar from "@/components/Sidebar.svelte";
</script>

<ModeWatcher />
<div class="flex">
    <Sidebar />
    <main class="pl-20 mx-4 my-2">
        <slot />
    </main>
</div>
