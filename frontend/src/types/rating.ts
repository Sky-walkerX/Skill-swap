export interface RatingType {
	ratingId: string;
	userId: string;
	ratedById: string;
	score: number; // Score from 1 to 5
	comment: string | null;
	createdAt: string;
	updatedAt: string;
}
