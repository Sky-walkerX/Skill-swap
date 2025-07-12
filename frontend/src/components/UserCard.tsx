import { UserType } from '@/types/user';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Avatar } from '@/components/ui/avatar';
import { Star, MapPin, MessageSquare, ArrowRight } from 'lucide-react';
import Link from 'next/link';

interface UserCardProps {
	user: UserType;
	currentUserId?: string;
}

export function UserCard({ user, currentUserId }: UserCardProps) {
	const isCurrentUser = currentUserId === user.userId;

	return (
		<Card className="hover:shadow-md transition-shadow duration-200">
			<CardHeader className="pb-3">
				<div className="flex items-start justify-between">
					<div className="flex items-center space-x-3">
						<Avatar
							src={user.photoUrl || undefined}
							alt={user.name}
							fallback={user.name
								.split(' ')
								.map((n) => n[0])
								.join('')}
							className="w-12 h-12"
						/>
						<div>
							<h3 className="font-semibold text-lg">{user.name}</h3>
							<div className="flex items-center space-x-2 text-sm text-muted-foreground">
								{user.location && (
									<>
										<MapPin className="w-4 h-4" />
										<span>{user.location}</span>
									</>
								)}
							</div>
						</div>
					</div>
					<div className="flex items-center space-x-1">
						<Star className="w-4 h-4 fill-yellow-400 text-yellow-400" />
						<span className="text-sm font-medium">
							{user.rating ? user.rating.toFixed(1) : 'No rating'}
						</span>
					</div>
				</div>
			</CardHeader>

			<CardContent className="space-y-4">
				{/* Skills Offered */}
				<div>
					<h4 className="text-sm font-medium text-muted-foreground mb-2">
						Skills Offered
					</h4>
					<div className="flex flex-wrap gap-1">
						{user.skillsOffered
							.filter((skill) => skill)
							.map((skill) => (
								<Badge key={skill.skillId} variant="secondary" className="text-xs">
									{skill.name}
								</Badge>
							))}
					</div>
				</div>

				{/* Skills Wanted */}
				<div>
					<h4 className="text-sm font-medium text-muted-foreground mb-2">
						Skills Wanted
					</h4>
					<div className="flex flex-wrap gap-1">
						{user.skillsWanted
							.filter((skill) => skill)
							.map((skill) => (
								<Badge key={skill.skillId} variant="outline" className="text-xs">
									{skill.name}
								</Badge>
							))}
					</div>
				</div>

				{/* Actions */}
				<div className="flex space-x-2 pt-2">
					{!isCurrentUser ? (
						<>
							<Button variant="outline" size="sm" className="flex-1">
								<MessageSquare className="w-4 h-4 mr-2" />
								Message
							</Button>
							<Link href={`/profile/${user.userId}`}>
								<Button size="sm" className="flex-1">
									<ArrowRight className="w-4 h-4 mr-2" />
									View Profile
								</Button>
							</Link>
						</>
					) : (
						<Button variant="outline" size="sm" className="w-full" disabled>
							This is you
						</Button>
					)}
				</div>
			</CardContent>
		</Card>
	);
}
