'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Check } from 'lucide-react'
import { dummySkills } from '@/lib/dummy-data'

export default function RegisterPage() {
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: '',
    location: '',
    skillsOffered: [] as string[],
    skillsWanted: [] as string[]
  })
  const [step, setStep] = useState(1)

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    })
  }

  const toggleSkill = (skillId: string, type: 'offered' | 'wanted') => {
    if (type === 'offered') {
      setFormData(prev => ({
        ...prev,
        skillsOffered: prev.skillsOffered.includes(skillId)
          ? prev.skillsOffered.filter(id => id !== skillId)
          : [...prev.skillsOffered, skillId]
      }))
    } else {
      setFormData(prev => ({
        ...prev,
        skillsWanted: prev.skillsWanted.includes(skillId)
          ? prev.skillsWanted.filter(id => id !== skillId)
          : [...prev.skillsWanted, skillId]
      }))
    }
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (step < 3) {
      setStep(step + 1)
    } else {
      // Handle registration
      console.log('Registering user:', formData)
    }
  }

  const isStepValid = () => {
    switch (step) {
      case 1:
        return formData.name && formData.email && formData.password && formData.confirmPassword && formData.password === formData.confirmPassword
      case 2:
        return formData.skillsOffered.length > 0 && formData.skillsWanted.length > 0
      case 3:
        return true
      default:
        return false
    }
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-foreground">Join SkillShare</h1>
          <p className="mt-2 text-muted-foreground">
            Connect with others and exchange your skills
          </p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Create your account</CardTitle>
            <CardDescription>
              Step {step} of 3: {step === 1 ? 'Basic Information' : step === 2 ? 'Skills' : 'Review'}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-6">
              {step === 1 && (
                <div className="space-y-4">
                  <div>
                    <Label htmlFor="name">Full Name</Label>
                    <Input
                      id="name"
                      name="name"
                      type="text"
                      value={formData.name}
                      onChange={handleInputChange}
                      placeholder="Enter your full name"
                      required
                    />
                  </div>
                  <div>
                    <Label htmlFor="email">Email</Label>
                    <Input
                      id="email"
                      name="email"
                      type="email"
                      value={formData.email}
                      onChange={handleInputChange}
                      placeholder="Enter your email"
                      required
                    />
                  </div>
                  <div>
                    <Label htmlFor="password">Password</Label>
                    <Input
                      id="password"
                      name="password"
                      type="password"
                      value={formData.password}
                      onChange={handleInputChange}
                      placeholder="Create a password"
                      required
                    />
                  </div>
                  <div>
                    <Label htmlFor="confirmPassword">Confirm Password</Label>
                    <Input
                      id="confirmPassword"
                      name="confirmPassword"
                      type="password"
                      value={formData.confirmPassword}
                      onChange={handleInputChange}
                      placeholder="Confirm your password"
                      required
                    />
                    {formData.confirmPassword && formData.password !== formData.confirmPassword && (
                      <p className="text-sm text-destructive mt-1">Passwords don&apos;t match</p>
                    )}
                  </div>
                  <div>
                    <Label htmlFor="location">Location</Label>
                    <Input
                      id="location"
                      name="location"
                      type="text"
                      value={formData.location}
                      onChange={handleInputChange}
                      placeholder="City, State"
                    />
                  </div>
                </div>
              )}

              {step === 2 && (
                <div className="space-y-6">
                  <div>
                    <Label>Skills You Can Offer</Label>
                    <p className="text-sm text-muted-foreground mb-3">
                      Select the skills you&apos;re confident teaching to others
                    </p>
                    <div className="grid grid-cols-2 gap-2">
                      {dummySkills.map((skill) => (
                        <button
                          key={skill.skillId}
                          type="button"
                          onClick={() => toggleSkill(skill.skillId, 'offered')}
                          className={`p-3 text-left rounded-lg border transition-colors ${
                            formData.skillsOffered.includes(skill.skillId)
                              ? 'border-primary bg-primary/10'
                              : 'border-border hover:border-primary/50'
                          }`}
                        >
                          <div className="flex items-center justify-between">
                            <span className="text-sm font-medium">{skill.name}</span>
                            {formData.skillsOffered.includes(skill.skillId) && (
                              <Check className="w-4 h-4 text-primary" />
                            )}
                          </div>
                        </button>
                      ))}
                    </div>
                  </div>

                  <div>
                    <Label>Skills You Want to Learn</Label>
                    <p className="text-sm text-muted-foreground mb-3">
                      Select the skills you&apos;d like to learn from others
                    </p>
                    <div className="grid grid-cols-2 gap-2">
                      {dummySkills.map((skill) => (
                        <button
                          key={skill.skillId}
                          type="button"
                          onClick={() => toggleSkill(skill.skillId, 'wanted')}
                          className={`p-3 text-left rounded-lg border transition-colors ${
                            formData.skillsWanted.includes(skill.skillId)
                              ? 'border-primary bg-primary/10'
                              : 'border-border hover:border-primary/50'
                          }`}
                        >
                          <div className="flex items-center justify-between">
                            <span className="text-sm font-medium">{skill.name}</span>
                            {formData.skillsWanted.includes(skill.skillId) && (
                              <Check className="w-4 h-4 text-primary" />
                            )}
                          </div>
                        </button>
                      ))}
                    </div>
                  </div>
                </div>
              )}

              {step === 3 && (
                <div className="space-y-4">
                  <div className="p-4 bg-muted rounded-lg">
                    <h3 className="font-medium mb-2">Account Summary</h3>
                    <div className="space-y-2 text-sm">
                      <p><strong>Name:</strong> {formData.name}</p>
                      <p><strong>Email:</strong> {formData.email}</p>
                      <p><strong>Location:</strong> {formData.location || 'Not specified'}</p>
                    </div>
                  </div>
                  
                  <div>
                    <h4 className="font-medium mb-2">Skills Offered ({formData.skillsOffered.length})</h4>
                    <div className="flex flex-wrap gap-1">
                      {formData.skillsOffered.map(skillId => {
                        const skill = dummySkills.find(s => s.skillId === skillId)
                        return skill ? (
                          <Badge key={skillId} variant="secondary">
                            {skill.name}
                          </Badge>
                        ) : null
                      })}
                    </div>
                  </div>

                  <div>
                    <h4 className="font-medium mb-2">Skills Wanted ({formData.skillsWanted.length})</h4>
                    <div className="flex flex-wrap gap-1">
                      {formData.skillsWanted.map(skillId => {
                        const skill = dummySkills.find(s => s.skillId === skillId)
                        return skill ? (
                          <Badge key={skillId} variant="outline">
                            {skill.name}
                          </Badge>
                        ) : null
                      })}
                    </div>
                  </div>
                </div>
              )}

              <div className="flex space-x-3">
                {step > 1 && (
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => setStep(step - 1)}
                    className="flex-1"
                  >
                    Back
                  </Button>
                )}
                <Button
                  type="submit"
                  disabled={!isStepValid()}
                  className="flex-1"
                >
                  {step === 3 ? 'Create Account' : 'Next'}
                </Button>
              </div>
            </form>

            <div className="mt-6 text-center">
              <p className="text-sm text-muted-foreground">
                Already have an account?{' '}
                <Link href="/auth/signin" className="text-primary hover:underline">
                  Sign in
                </Link>
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
} 