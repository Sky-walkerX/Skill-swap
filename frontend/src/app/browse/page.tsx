'use client';

import { useState, useMemo } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { UserCard } from '@/components/UserCard';
import { dummyUsers, dummySkills } from '@/lib/dummy-data';
import { Search, Filter, MapPin } from 'lucide-react';

export default function BrowsePage() {
	const [searchTerm, setSearchTerm] = useState('');
	const [selectedSkills, setSelectedSkills] = useState<string[]>([]);
	const [locationFilter, setLocationFilter] = useState('');

	const filteredUsers = useMemo(() => {
		return dummyUsers.filter((user) => {
			// Search by name
			const matchesSearch = user.name.toLowerCase().includes(searchTerm.toLowerCase());

			// Filter by skills
			const matchesSkills =
				selectedSkills.length === 0 ||
				selectedSkills.some(
					(skillId) =>
						user.skillsOffered
							.filter((skill) => skill)
							.some((skill) => skill.skillId === skillId) ||
						user.skillsWanted
							.filter((skill) => skill)
							.some((skill) => skill.skillId === skillId)
				);

			// Filter by location
			const matchesLocation =
				!locationFilter ||
				user.location?.toLowerCase().includes(locationFilter.toLowerCase());

			return matchesSearch && matchesSkills && matchesLocation;
		});
	}, [searchTerm, selectedSkills, locationFilter]);

	const toggleSkillFilter = (skillId: string) => {
		setSelectedSkills((prev) =>
			prev.includes(skillId) ? prev.filter((id) => id !== skillId) : [...prev, skillId]
		);
	};

	const clearFilters = () => {
		setSearchTerm('');
		setSelectedSkills([]);
		setLocationFilter('');
	};

	const uniqueLocations = [
		...new Set(
			dummyUsers
				.map((user) => user.location)
				.filter((location): location is string => Boolean(location))
		)
	];

	return (
		<div className="min-h-screen bg-background">
			<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
				{/* Header */}
				<div className="mb-8">
					<h1 className="text-3xl font-bold text-foreground mb-2">Discover People</h1>
					<p className="text-muted-foreground">
						Find people to exchange skills with in your area
					</p>
				</div>

				{/* Search and Filters */}
				<div className="mb-8 space-y-4">
					{/* Search Bar */}
					<div className="relative">
						<Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
						<Input
							placeholder="Search by name..."
							value={searchTerm}
							onChange={(e) => setSearchTerm(e.target.value)}
							className="pl-10"
						/>
					</div>

					{/* Filters */}
					<div className="flex flex-wrap gap-4 items-center">
						<div className="flex items-center space-x-2">
							<Filter className="w-4 h-4 text-muted-foreground" />
							<span className="text-sm font-medium">Filters:</span>
						</div>

						{/* Location Filter */}
						<div className="flex items-center space-x-2">
							<MapPin className="w-4 h-4 text-muted-foreground" />
							<select
								value={locationFilter}
								onChange={(e) => setLocationFilter(e.target.value)}
								className="text-sm border border-border rounded-md px-3 py-1 bg-background"
							>
								<option value="">All Locations</option>
								{uniqueLocations.map((location) => (
									<option key={location} value={location}>
										{location}
									</option>
								))}
							</select>
						</div>

						{/* Skill Filters */}
						<div className="flex flex-wrap gap-2">
							{dummySkills.map((skill) => (
								<Badge
									key={skill.skillId}
									variant={
										selectedSkills.includes(skill.skillId)
											? 'default'
											: 'outline'
									}
									className="cursor-pointer hover:bg-primary/10"
									onClick={() => toggleSkillFilter(skill.skillId)}
								>
									{skill.name}
								</Badge>
							))}
						</div>

						{(searchTerm || selectedSkills.length > 0 || locationFilter) && (
							<Button
								variant="outline"
								size="sm"
								onClick={clearFilters}
								className="ml-auto"
							>
								Clear Filters
							</Button>
						)}
					</div>
				</div>

				{/* Results */}
				<div className="mb-6">
					<div className="flex items-center justify-between">
						<h2 className="text-xl font-semibold">
							{filteredUsers.length}{' '}
							{filteredUsers.length === 1 ? 'person' : 'people'} found
						</h2>
						{filteredUsers.length > 0 && (
							<p className="text-sm text-muted-foreground">
								Showing {filteredUsers.length} of {dummyUsers.length} users
							</p>
						)}
					</div>
				</div>

				{/* User Grid */}
				{filteredUsers.length > 0 ? (
					<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
						{filteredUsers.map((user) => (
							<UserCard
								key={user.userId}
								user={user}
								currentUserId="1" // Assuming current user is Sarah Johnson
							/>
						))}
					</div>
				) : (
					<div className="text-center py-12">
						<div className="w-16 h-16 mx-auto mb-4 bg-muted rounded-full flex items-center justify-center">
							<Search className="w-8 h-8 text-muted-foreground" />
						</div>
						<h3 className="text-lg font-medium text-foreground mb-2">No users found</h3>
						<p className="text-muted-foreground mb-4">
							Try adjusting your search terms or filters
						</p>
						<Button onClick={clearFilters}>Clear all filters</Button>
					</div>
				)}

				{/* Popular Skills Section */}
				<div className="mt-16">
					<h2 className="text-2xl font-bold text-foreground mb-6">Popular Skills</h2>
					<div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-4">
						{dummySkills.map((skill) => {
							const userCount = dummyUsers.filter(
								(user) =>
									user.skillsOffered
										.filter((s) => s)
										.some((s) => s.skillId === skill.skillId) ||
									user.skillsWanted
										.filter((s) => s)
										.some((s) => s.skillId === skill.skillId)
							).length;

							return (
								<div
									key={skill.skillId}
									className="p-4 border border-border rounded-lg hover:border-primary/50 transition-colors cursor-pointer"
									onClick={() => toggleSkillFilter(skill.skillId)}
								>
									<h3 className="font-medium text-sm mb-1">{skill.name}</h3>
									<p className="text-xs text-muted-foreground mb-2">
										{skill.description}
									</p>
									<p className="text-xs text-primary">
										{userCount} {userCount === 1 ? 'person' : 'people'}
									</p>
								</div>
							);
						})}
					</div>
				</div>
			</div>
		</div>
	);
}
