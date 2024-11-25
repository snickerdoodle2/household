import { z } from 'zod';

export enum NotificationType {
    Info = 'info',
    Error = 'error',
    Warning = 'warning',
    Success = 'success',
}

export const notificationSchema = z.object({
    type: z.nativeEnum(NotificationType),
    name: z.string().min(1).max(64),
    description: z.string().max(256),
    date: z
        .string()
        .or(z.date())
        .transform((d) => new Date(d)),
});

export type Notification = z.infer<typeof notificationSchema>;

// TODO: remove this after backend integration
export const dummyNotifications: Notification[] = [
    {
        type: NotificationType.Info,
        name: 'System Update',
        description:
            'A new system update is available. Please restart to apply the latest features and improvements.',
        date: new Date('2024-11-01T10:00:00Z'),
    },
    {
        type: NotificationType.Error,
        name: 'Connection Error',
        description:
            'Failed to connect to the server. Please check your internet connection and try again.',
        date: new Date('2024-11-02T14:30:00Z'),
    },
    {
        type: NotificationType.Success,
        name: 'Settings Saved',
        description: 'Your changes have been saved successfully.',
        date: new Date('2024-11-02T15:00:00Z'),
    },
    {
        type: NotificationType.Warning,
        name: 'Low Battery',
        description:
            'Battery level is below 20%. Please charge your device soon to avoid shutdown.',
        date: new Date('2024-11-03T08:00:00Z'),
    },
    {
        type: NotificationType.Info,
        name: 'Welcome Back!',
        description:
            'We’re glad to see you again. Check out what’s new in your dashboard.',
        date: new Date('2024-11-03T09:15:00Z'),
    },
    {
        type: NotificationType.Success,
        name: 'Subscription Activated',
        description:
            'Your subscription has been activated. Enjoy all the premium features!',
        date: new Date('2024-11-04T10:45:00Z'),
    },
    {
        type: NotificationType.Warning,
        name: 'Storage Almost Full',
        description:
            'Your storage is almost full. Consider deleting some files or upgrading your plan.',
        date: new Date('2024-11-04T12:00:00Z'),
    },
    {
        type: NotificationType.Error,
        name: 'File Upload Failed',
        description:
            'An error occurred while uploading your file. Please try again later.',
        date: new Date('2024-11-05T16:20:00Z'),
    },
    {
        type: NotificationType.Success,
        name: 'Profile Updated',
        description: 'Your profile information has been successfully updated.',
        date: new Date('2024-11-06T11:30:00Z'),
    },
    {
        type: NotificationType.Success,
        name: 'Payment Successful',
        description: 'Your payment has been processed successfully. Thank you!',
        date: new Date('2024-11-07T09:00:00Z'),
    },
];
