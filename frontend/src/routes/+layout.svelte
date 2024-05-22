<script>
import "../app.css";
import { ModeWatcher } from "mode-watcher";

import Sidebar from "@/components/Sidebar.svelte";
</script>

<ModeWatcher />
<div class="flex min-h-screen">
    <Sidebar />
    <main class="pl-20 mx-4 my-2 flex-1 flex justify-center items-center">
        <slot />
    </main>
</div>
