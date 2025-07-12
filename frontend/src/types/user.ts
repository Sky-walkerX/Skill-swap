import { SkillType } from './skill';

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
