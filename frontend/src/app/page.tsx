"use client";
import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import HeroSection from '@/components/Hero';
import FeaturesSection from '@/components/Features';
import HowItWorksSection from '@/components/HowItWorks';
import SkillsMarketplace from '@/components/Skills';
import CallToActionSection from '@/components/CTA';
import LoadingAnimation from '@/components/LoadingAnimation';

function Home() {
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Simulate loading time
    const timer = setTimeout(() => {
      setIsLoading(false);
    }, 3000);

    return () => clearTimeout(timer);
  }, []);

  return (
    <div className="min-h-screen">
      <AnimatePresence mode="wait">
        {isLoading ? (
          <LoadingAnimation key="loading" />
        ) : (
          <motion.div
            key="main"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.5 }}
            className="overflow-x-hidden"
          >
            <HeroSection />
            <FeaturesSection />
            <HowItWorksSection />
            <SkillsMarketplace />
            <CallToActionSection />
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}

export default Home;