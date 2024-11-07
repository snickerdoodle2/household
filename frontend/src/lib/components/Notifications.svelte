<script lang="ts">
    import * as Card from '$lib/components/ui/card';
    import { dummyNotifications, NotificationType } from '@/types/notification';
    import { CheckCircled, CrossCircled, InfoCircled, ExclamationTriangle } from 'radix-icons-svelte';

    const notifications = dummyNotifications.sort((n1, n2) => n2.date.getTime() - n1.date.getTime());

    const notificationMap = {
        [NotificationType.Success]: {
            icon: CheckCircled,
            iconClass: 'text-green-600',
            borderClass: 'border-green-600',
            textClass: 'text-green-600'
        },
        [NotificationType.Info]: {
            icon: InfoCircled,
            iconClass: 'text-blue-600',
            borderClass: 'border-blue-600',
            textClass: 'text-blue-600'
        },
        [NotificationType.Error]: {
            icon: CrossCircled,
            iconClass: 'text-red-600',
            borderClass: 'border-red-600',
            textClass: 'text-red-600'
        },
        [NotificationType.Warning]: {
            icon: ExclamationTriangle,
            iconClass: 'text-yellow-500 scale-90',
            borderClass: 'border-yellow-500',
            textClass: 'text-yellow-500'
        }
    };
</script>

<Card.Root class="w-[600px] border-none shadow-none">
    <Card.Header class="text-3xl">
        <Card.Title>Notifications</Card.Title>
    </Card.Header>
    <Card.Content
        class="flex flex-col items-start gap-3 overflow-y-auto max-h-[600px]"
    >
        {#each notifications as notification}
            {#if notificationMap[notification.type]}
                <div
                    class="flex relative w-full border rounded-md p-2 border-l-8 {notificationMap[
                        notification.type
                    ].borderClass}"
                >
                    <!-- Icon  -->
                    <div class="pr-4 pl-2 pt-2 items-start">
                        <svelte:component
                            this={notificationMap[notification.type].icon}
                            class="w-6 h-6 {notificationMap[notification.type]
                                .iconClass}"
                        />
                    </div>

                    <!-- Date -->
                    <p class="text-xs text-gray-500 absolute top-2 right-2">
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
