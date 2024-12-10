<script lang="ts">
    import type { PageData } from './$types';
    import * as Table from '$lib/components/ui/table';
    import { Button } from '@/components/ui/button';
    import { UserPen } from 'lucide-svelte';

    type Props = {
        data: PageData;
    };

    const { data }: Props = $props();
    const { users } = data;
</script>

<Table.Root class="max-w-[1024px] mx-auto">
    <Table.Header>
        <Table.Row>
            <Table.Cell>Username</Table.Cell>
            <Table.Cell>Name</Table.Cell>
            <Table.Cell>Created At</Table.Cell>
            <Table.Cell>Action</Table.Cell>
        </Table.Row>
    </Table.Header>
    <Table.Body>
        {#each users as user}
            <Table.Row>
                <Table.Cell>{user.username}</Table.Cell>
                <Table.Cell>{user.name}</Table.Cell>
                <Table.Cell>{user.created_at.toLocaleString()}</Table.Cell>
                <Table.Cell
                    ><Button
                        size="icon"
                        variant="ghost"
                        href={`/settings/users/${user.id}`}><UserPen /></Button
                    ></Table.Cell
                >
            </Table.Row>
        {/each}
    </Table.Body>
</Table.Root>
