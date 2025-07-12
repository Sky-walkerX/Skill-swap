// ---------------------------------------------------------------------------
// Skills
export interface SkillType {
	skillId: string;
	name: string;
	description: string | null;
}

// ---------------------------------------------------------------------------
// User
export type UserRoleType = 'user' | 'admin';

export interface UserType {
	userId: string;
	name: string;
	email: string;
	role: UserRoleType;
	rating: number | null;
	location: string | null;
	photoUrl: string | null;
	isPublic: boolean;
	createdAt: string;
	updatedAt: string;
	deletedAt: string | null;
	skillsOffered: SkillType[];
	skillsWanted: SkillType[];
}

// ---------------------------------------------------------------------------
// Rating
export interface RatingType {
	ratingId: string;
	userId: string;
	ratedById: string;
	score: number; // Score from 1 to 5
	comment: string | null;
	createdAt: string;
	updatedAt: string;
}

// ---------------------------------------------------------------------------
// Swap Request
export type SwapStatusType = 'pending' | 'accepted' | 'rejected' | 'cancelled';

export interface SwapRequest {
	swapId: string;
	requesterId: string;
	responderId: string;
	offeredSkill: SkillType;
	wantedSkill: SkillType;
	status: SwapStatusType;
	createdAt: string;
	updatedAt: string;
	deletedAt: string | null;
}

// ---------------------------------------------------------------------------
// Message
export interface MessageType {
	messageId: string;
	senderId: string;
	receiverId: string;
	text: string | null;
	image: string | null;
	createdAt: string;
	updatedAt: string;
}

// ---------------------------------------------------------------------------
// Notification
export type NotificationType = 'swapRequest' | 'swapAccepted' | 'swapRejected' | 'messageReceived';

export interface Notification {
	notificationId: string;
	userId: string;
	type: NotificationType;
	content: string;
	isRead: boolean;
	createdAt: string;
	updatedAt: string;
}
