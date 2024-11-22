<script lang="ts">
import * as Card from '$lib/components/ui/card';
import { dummyNotifications, NotificationType } from '@/types/notification';
import {
    CheckCircled,
    CrossCircled,
    InfoCircled,
    ExclamationTriangle,
} from 'radix-icons-svelte';

const notifications = dummyNotifications.sort(
    (n1, n2) => n2.date.getTime() - n1.date.getTime()
);

const notificationMap = {
    [NotificationType.Success]: {
        icon: CheckCircled,
        iconClass: 'text-green-600',
        borderClass: 'border-green-600',
        textClass: 'text-green-600',
    },
    [NotificationType.Info]: {
        icon: InfoCircled,
        iconClass: 'text-blue-600',
        borderClass: 'border-blue-600',
        textClass: 'text-blue-600',
    },
    [NotificationType.Error]: {
        icon: CrossCircled,
        iconClass: 'text-red-600',
        borderClass: 'border-red-600',
        textClass: 'text-red-600',
    },
    [NotificationType.Warning]: {
        icon: ExclamationTriangle,
        iconClass: 'text-yellow-500 scale-90',
        borderClass: 'border-yellow-500',
        textClass: 'text-yellow-500',
    },
};
</script>

<Card.Root class="w-[600px] border-none shadow-none">
    <Card.Header class="text-3xl">
        <Card.Title>Notifications</Card.Title>
    </Card.Header>
    <Card.Content
        class="flex max-h-[600px] flex-col items-start gap-3 overflow-y-auto"
    >
        {#each notifications as notification}
            {#if notificationMap[notification.type]}
                <div
                    class="relative flex w-full rounded-md border border-l-8 p-2 {notificationMap[
                        notification.type
                    ].borderClass}"
                >
                    <!-- Icon  -->
                    <div class="items-start pl-2 pr-4 pt-2">
                        <svelte:component
                            this={notificationMap[notification.type].icon}
                            class="h-6 w-6 {notificationMap[notification.type]
                                .iconClass}"
                        />
                    </div>

                    <!-- Date -->
                    <p class="absolute right-2 top-2 text-xs text-gray-500">
                        {new Date(notification.date).toLocaleDateString(
                            'en-US',
                            {
                                weekday: 'short',
                                year: 'numeric',
                                month: 'short',
                                day: 'numeric',
                            }
                        )}
                        {new Date(notification.date).toLocaleTimeString(
                            'en-US',
                            {
                                hour: '2-digit',
                                minute: '2-digit',
                            }
                        )}
                    </p>

                    <!-- Content -->
                    <div>
                        <h3
                            class="text-xl font-bold {notificationMap[
                                notification.type
                            ].textClass}"
                        >
                            {notification.name}
                        </h3>
                        <p class="text-sm">{notification.description}</p>
                    </div>
                </div>
            {/if}
        {/each}
    </Card.Content>
    <Card.Footer class="flex justify-end gap-3"></Card.Footer>
</Card.Root>
