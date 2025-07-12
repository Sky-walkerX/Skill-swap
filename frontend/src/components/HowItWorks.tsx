import { useEffect, useRef, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { gsap } from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { UserPlus, Search, MessageCircle, TrendingUp, Play } from 'lucide-react';

gsap.registerPlugin(ScrollTrigger);

const HowItWorksSection = () => {
  const sectionRef = useRef(null);
  const [activeStep, setActiveStep] = useState(0);
  const [isPlaying, setIsPlaying] = useState(false);

  const steps = [
    {
      icon: UserPlus,
      title: "Create Your Profile",
      description: "Sign up and showcase your skills. Tell us what you're good at and what you want to learn.",
      details: "Upload your portfolio, set your availability, and define your skill level to attract the right matches."
    },
    {
      icon: Search,
      title: "Find Perfect Matches",
      description: "Our AI algorithm connects you with people who have the skills you need and want what you offer.",
      details: "Browse through verified profiles, check ratings, and use advanced filters to find your ideal skill exchange partner."
    },
    {
      icon: MessageCircle,
      title: "Start Exchanging",
      description: "Connect with your matches, schedule sessions, and begin your skill exchange journey.",
      details: "Use our built-in video chat, collaborative tools, and progress tracking to make the most of your exchanges."
    },
    {
      icon: TrendingUp,
      title: "Track Progress",
      description: "Monitor your learning goals, earn badges, and build your professional network as you grow.",
      details: "Set milestones, receive feedback, and showcase your newly acquired skills to advance your career."
    }
  ];

  useEffect(() => {
    const ctx = gsap.context(() => {
      // Animate step cards on scroll
      steps.forEach((_, index) => {
        gsap.fromTo(`.step-card-${index}`,
          { x: index % 2 === 0 ? -100 : 100, opacity: 0 },
          {
            x: 0,
            opacity: 1,
            duration: 0.8,
            scrollTrigger: {
              trigger: `.step-card-${index}`,
              start: "top 80%",
              end: "bottom 20%",
              toggleActions: "play none none reverse"
            }
          }
        );
      });

      // Animate connecting lines
      gsap.fromTo('.connecting-line',
        { scaleY: 0 },
        {
          scaleY: 1,
          duration: 1,
          stagger: 0.3,
          scrollTrigger: {
            trigger: '.steps-container',
            start: "top 60%",
            end: "bottom 40%",
            toggleActions: "play none none reverse"
          }
        }
      );
    }, sectionRef);

    return () => ctx.revert();
  }, []);

  const handlePlayDemo = () => {
    setIsPlaying(true);
    // Simulate auto-playing through steps
    let currentStep = 0;
    const interval = setInterval(() => {
      currentStep = (currentStep + 1) % steps.length;
      setActiveStep(currentStep);
      
      if (currentStep === 0) {
        clearInterval(interval);
        setIsPlaying(false);
      }
    }, 2000);
  };

  return (
    <section ref={sectionRef} className="py-20 bg-gradient-to-br from-gray-50 to-blue-50 dark:from-gray-900 dark:to-blue-900">
      <div className="container mx-auto px-4">
        {/* Section Header */}
        <motion.div
          initial={{ opacity: 0, y: 50 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
          className="text-center mb-16"
        >
          <h2 className="text-3xl md:text-5xl font-bold mb-6 text-gray-900 dark:text-white">
            How It Works
          </h2>
          <p className="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto mb-8">
            Getting started with skill exchange is simple. Follow these four easy steps to begin your learning journey.
          </p>
          
          <Button
            onClick={handlePlayDemo}
            disabled={isPlaying}
            className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white px-6 py-3"
          >
            <Play className="mr-2 h-4 w-4" />
            {isPlaying ? 'Playing Demo...' : 'Watch Interactive Demo'}
          </Button>
        </motion.div>

        {/* Steps Container */}
        <div className="steps-container relative max-w-4xl mx-auto">
          {steps.map((step, index) => (
            <div key={index} className="relative mb-16 last:mb-0">
              {/* Connecting Line */}
              {index < steps.length - 1 && (
                <div 
                  className="connecting-line absolute left-1/2 top-full w-0.5 h-16 bg-gradient-to-b from-blue-600 to-purple-600 transform -translate-x-1/2 z-0"
                  style={{ transformOrigin: 'top' }}
                />
              )}
              
              {/* Step Card */}
              <motion.div
                className={`step-card-${index} relative z-10`}
                onMouseEnter={() => setActiveStep(index)}
                whileHover={{ scale: 1.02 }}
                transition={{ duration: 0.3 }}
              >
                <Card className={`transition-all duration-300 ${
                  activeStep === index 
                    ? 'shadow-2xl border-blue-500 dark:border-blue-400' 
                    : 'shadow-lg hover:shadow-xl'
                }`}>
                  <CardContent className="p-8">
                    <div className="flex flex-col md:flex-row items-center text-center md:text-left">
                      {/* Step Icon */}
                      <div className={`flex-shrink-0 mb-6 md:mb-0 md:mr-8 p-4 rounded-full transition-all duration-300 ${
                        activeStep === index
                          ? 'bg-gradient-to-r from-blue-600 to-purple-600 scale-110'
                          : 'bg-gray-200 dark:bg-gray-700'
                      }`}>
                        <step.icon className={`h-8 w-8 ${
                          activeStep === index ? 'text-white' : 'text-gray-600 dark:text-gray-300'
                        }`} />
                      </div>
                      
                      {/* Step Content */}
                      <div className="flex-1">
                        <div className="flex items-center justify-center md:justify-start mb-4">
                          <span className="text-sm font-semibold text-blue-600 dark:text-blue-400 mr-4">
                            Step {index + 1}
                          </span>
                          <h3 className="text-2xl font-bold text-gray-900 dark:text-white">
                            {step.title}
                          </h3>
                        </div>
                        
                        <p className="text-lg text-gray-600 dark:text-gray-300 mb-4">
                          {step.description}
                        </p>
                        
                        <AnimatePresence>
                          {activeStep === index && (
                            <motion.p
                              initial={{ opacity: 0, height: 0 }}
                              animate={{ opacity: 1, height: 'auto' }}
                              exit={{ opacity: 0, height: 0 }}
                              transition={{ duration: 0.3 }}
                              className="text-gray-500 dark:text-gray-400 leading-relaxed"
                            >
                              {step.details}
                            </motion.p>
                          )}
                        </AnimatePresence>
                      </div>
                      
                      {/* Step Number */}
                      <div className={`flex-shrink-0 mt-6 md:mt-0 md:ml-8 w-12 h-12 rounded-full flex items-center justify-center text-2xl font-bold transition-all duration-300 ${
                        activeStep === index
                          ? 'bg-gradient-to-r from-blue-600 to-purple-600 text-white scale-110'
                          : 'bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300'
                      }`}>
                        {index + 1}
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </motion.div>
            </div>
          ))}
        </div>

        {/* Interactive Demo Visualization */}
        <motion.div
          initial={{ opacity: 0, y: 50 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.3 }}
          viewport={{ once: true }}
          className="mt-20 text-center"
        >
          <div className="max-w-2xl mx-auto p-8 bg-white dark:bg-gray-800 rounded-2xl shadow-xl">
            <h3 className="text-xl font-bold mb-6 text-gray-900 dark:text-white">
              Ready to get started?
            </h3>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white px-8 py-3">
                Join Now - It's Free
              </Button>
              <Button variant="outline" className="px-8 py-3 border-2">
                Learn More
              </Button>
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  );
};

export default HowItWorksSection;