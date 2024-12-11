<script lang="ts">
    import UserManagementForm from './UserManagementForm.svelte';
    import type { PageData } from './$types';
    import type { LayoutData } from '../../$types';
    type Props = {
        data: PageData;
        layoutData: LayoutData;
    };
    const { data }: Props = $props();
    let { user, currentUser } = data;
</script>

<UserManagementForm action="edit" {user} isMe={currentUser.id === user.id} />
