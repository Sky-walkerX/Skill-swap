"use client"

import { useState, useEffect } from "react"
import LoadingAnimation from "@/components/LoadingAnimation"
import HeroSection from "@/components/Hero"
import FeaturesSection from "@/components/Features"
import SkillsMarketplace from "@/components/Skills"
import HowItWorksSection from "@/components/HowItWorks"
import CallToActionSection from "@/components/CTA"

export default function Home() {
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsLoading(false)
    }, 3000)

    return () => clearTimeout(timer)
  }, [])

  if (isLoading) {
    return <LoadingAnimation />
  }

  return (
    <main className="min-h-screen">
      <HeroSection />
      <FeaturesSection />
      <SkillsMarketplace />
      <HowItWorksSection />
      <CallToActionSection />
    </main>
  )
}
