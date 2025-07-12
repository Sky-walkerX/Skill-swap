export interface MessageType {
	messageId: string;
	senderId: string;
	receiverId: string;
	text: string | null;
	image: string | null;
	createdAt: string;
	updatedAt: string;
}
