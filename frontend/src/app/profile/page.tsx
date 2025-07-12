'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Avatar } from '@/components/ui/avatar'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  Star, 
  MapPin, 
  Calendar, 
  Edit, 
  MessageSquare, 
  CheckCircle,
  XCircle,
  Clock
} from 'lucide-react'
import { 
  getUserById,
  getSwapRequestsByUserId,
  getRatingsByUserId
} from '@/lib/dummy-data'

export default function ProfilePage() {
  const [activeTab, setActiveTab] = useState('overview')
  
  // For demo purposes, we'll show Sarah Johnson's profile (userId: '1')
  const currentUserId = '1'
  const user = getUserById(currentUserId)
  const swapRequests = getSwapRequestsByUserId(currentUserId)
  const ratings = getRatingsByUserId(currentUserId)

  if (!user) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-foreground mb-2">User not found</h1>
          <p className="text-muted-foreground">The user you&apos;re looking for doesn&apos;t exist.</p>
        </div>
      </div>
    )
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'accepted':
        return <CheckCircle className="w-4 h-4 text-green-500" />
      case 'rejected':
        return <XCircle className="w-4 h-4 text-red-500" />
      case 'pending':
        return <Clock className="w-4 h-4 text-yellow-500" />
      default:
        return null
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'accepted':
        return 'text-green-600'
      case 'rejected':
        return 'text-red-600'
      case 'pending':
        return 'text-yellow-600'
      default:
        return 'text-muted-foreground'
    }
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Profile Header */}
        <div className="bg-card rounded-lg p-6 mb-8">
          <div className="flex flex-col md:flex-row items-start md:items-center space-y-4 md:space-y-0 md:space-x-6">
            <Avatar
              src={user.photoUrl || undefined}
              alt={user.name}
              fallback={user.name.split(' ').map(n => n[0]).join('')}
              className="w-24 h-24"
            />
            
            <div className="flex-1">
              <div className="flex flex-col md:flex-row md:items-center md:justify-between">
                <div>
                  <h1 className="text-3xl font-bold text-foreground mb-2">{user.name}</h1>
                  <div className="flex items-center space-x-4 text-muted-foreground mb-3">
                    {user.location && (
                      <div className="flex items-center space-x-1">
                        <MapPin className="w-4 h-4" />
                        <span>{user.location}</span>
                      </div>
                    )}
                    <div className="flex items-center space-x-1">
                      <Calendar className="w-4 h-4" />
                      <span>Joined {new Date(user.createdAt).toLocaleDateString()}</span>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Star className="w-5 h-5 fill-yellow-400 text-yellow-400" />
                    <span className="font-medium">
                      {user.rating ? user.rating.toFixed(1) : 'No rating'}
                    </span>
                    <span className="text-muted-foreground">
                      ({ratings.length} {ratings.length === 1 ? 'review' : 'reviews'})
                    </span>
                  </div>
                </div>
                
                <div className="flex space-x-3 mt-4 md:mt-0">
                  <Button variant="outline">
                    <Edit className="w-4 h-4 mr-2" />
                    Edit Profile
                  </Button>
                  <Button>
                    <MessageSquare className="w-4 h-4 mr-2" />
                    Message
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Tabs */}
        <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="skills">Skills</TabsTrigger>
            <TabsTrigger value="swaps">Swap History</TabsTrigger>
            <TabsTrigger value="reviews">Reviews</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Skills Offered</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex flex-wrap gap-2">
                    {user.skillsOffered.filter(skill => skill).map((skill) => (
                      <Badge key={skill.skillId} variant="secondary">
                        {skill.name}
                      </Badge>
                    ))}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Skills Wanted</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex flex-wrap gap-2">
                    {user.skillsWanted.filter(skill => skill).map((skill) => (
                      <Badge key={skill.skillId} variant="outline">
                        {skill.name}
                      </Badge>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>

            <Card>
              <CardHeader>
                <CardTitle>Recent Activity</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {swapRequests.slice(0, 3).map((request) => {
                    const otherUser = request.requesterId === currentUserId 
                      ? getUserById(request.responderId)
                      : getUserById(request.requesterId)
                    
                    return (
                      <div key={request.swapId} className="flex items-center justify-between p-3 border border-border rounded-lg">
                        <div className="flex items-center space-x-3">
                          <Avatar
                            src={otherUser?.photoUrl || undefined}
                            alt={otherUser?.name || 'User'}
                            fallback={otherUser?.name?.split(' ').map(n => n[0]).join('') || 'U'}
                            className="w-10 h-10"
                          />
                          <div>
                            <p className="font-medium">
                              {request.requesterId === currentUserId ? 'You offered' : `${otherUser?.name} offered`} {request.offeredSkill.name}
                            </p>
                            <p className="text-sm text-muted-foreground">
                              for {request.wantedSkill.name}
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center space-x-2">
                          {getStatusIcon(request.status)}
                          <span className={`text-sm font-medium ${getStatusColor(request.status)}`}>
                            {request.status.charAt(0).toUpperCase() + request.status.slice(1)}
                          </span>
                        </div>
                      </div>
                    )
                  })}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="skills" className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Skills I Can Teach</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {user.skillsOffered.map((skill) => (
                      <div key={skill.skillId} className="p-4 border border-border rounded-lg">
                        <h3 className="font-medium mb-2">{skill.name}</h3>
                        {skill.description && (
                          <p className="text-sm text-muted-foreground mb-3">{skill.description}</p>
                        )}
                        <Button variant="outline" size="sm">
                          <Edit className="w-4 h-4 mr-2" />
                          Edit
                        </Button>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Skills I Want to Learn</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {user.skillsWanted.map((skill) => (
                      <div key={skill.skillId} className="p-4 border border-border rounded-lg">
                        <h3 className="font-medium mb-2">{skill.name}</h3>
                        {skill.description && (
                          <p className="text-sm text-muted-foreground mb-3">{skill.description}</p>
                        )}
                        <Button variant="outline" size="sm">
                          <Edit className="w-4 h-4 mr-2" />
                          Edit
                        </Button>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="swaps" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Swap History</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {swapRequests.map((request) => {
                    const otherUser = request.requesterId === currentUserId 
                      ? getUserById(request.responderId)
                      : getUserById(request.requesterId)
                    
                    return (
                      <div key={request.swapId} className="p-4 border border-border rounded-lg">
                        <div className="flex items-center justify-between mb-3">
                          <div className="flex items-center space-x-3">
                            <Avatar
                              src={otherUser?.photoUrl || undefined}
                              alt={otherUser?.name || 'User'}
                              fallback={otherUser?.name?.split(' ').map(n => n[0]).join('') || 'U'}
                              className="w-12 h-12"
                            />
                            <div>
                              <h3 className="font-medium">{otherUser?.name}</h3>
                              <p className="text-sm text-muted-foreground">
                                {new Date(request.createdAt).toLocaleDateString()}
                              </p>
                            </div>
                          </div>
                          <div className="flex items-center space-x-2">
                            {getStatusIcon(request.status)}
                            <span className={`text-sm font-medium ${getStatusColor(request.status)}`}>
                              {request.status.charAt(0).toUpperCase() + request.status.slice(1)}
                            </span>
                          </div>
                        </div>
                        
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                          <div className="p-3 bg-muted rounded-lg">
                            <p className="text-sm font-medium text-muted-foreground mb-1">Offered</p>
                            <p className="font-medium">{request.offeredSkill.name}</p>
                          </div>
                          <div className="p-3 bg-muted rounded-lg">
                            <p className="text-sm font-medium text-muted-foreground mb-1">Wanted</p>
                            <p className="font-medium">{request.wantedSkill.name}</p>
                          </div>
                        </div>
                      </div>
                    )
                  })}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="reviews" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Reviews & Ratings</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {ratings.map((rating) => {
                    const reviewer = getUserById(rating.ratedById)
                    
                    return (
                      <div key={rating.ratingId} className="p-4 border border-border rounded-lg">
                        <div className="flex items-start justify-between mb-3">
                          <div className="flex items-center space-x-3">
                            <Avatar
                              src={reviewer?.photoUrl || undefined}
                              alt={reviewer?.name || 'User'}
                              fallback={reviewer?.name?.split(' ').map(n => n[0]).join('') || 'U'}
                              className="w-10 h-10"
                            />
                            <div>
                              <h3 className="font-medium">{reviewer?.name}</h3>
                              <div className="flex items-center space-x-1">
                                {[...Array(5)].map((_, i) => (
                                  <Star
                                    key={i}
                                    className={`w-4 h-4 ${
                                      i < rating.score ? 'fill-yellow-400 text-yellow-400' : 'text-gray-300'
                                    }`}
                                  />
                                ))}
                              </div>
                            </div>
                          </div>
                          <span className="text-sm text-muted-foreground">
                            {new Date(rating.createdAt).toLocaleDateString()}
                          </span>
                        </div>
                        {rating.comment && (
                          <p className="text-sm text-muted-foreground">{rating.comment}</p>
                        )}
                      </div>
                    )
                  })}
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  )
} 