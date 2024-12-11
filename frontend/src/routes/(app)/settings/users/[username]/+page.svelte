<script lang="ts">
    import UserManagementForm from './UserManagementForm.svelte';
    import type { PageData } from './$types';
    type Props = {
        data: PageData;
    };
    const { data }: Props = $props();
    let { user } = data;
</script>

<UserManagementForm action="edit" {user} />
