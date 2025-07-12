import { SkillType } from './skill';

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
