export type NotificationType = 'swapRequest' | 'swapAccepted' | 'swapRejected' | 'messageReceived';

export interface Notification {
	notificationId: string;
	type: NotificationType;
	content: string;
	isRead: boolean;
	createdAt: string;
	updatedAt: string;
}
