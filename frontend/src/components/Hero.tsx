import { useEffect, useRef } from 'react';
import { motion } from 'framer-motion';
import { gsap } from 'gsap';
import { Button } from '@/components/ui/button';
import { ArrowRight, Users, BookOpen, TrendingUp, Zap } from 'lucide-react';

const HeroSection = () => {
  const heroRef = useRef(null);
  const skillIconsRef = useRef<(HTMLDivElement | null)[]>([]);

  useEffect(() => {
    const ctx = gsap.context(() => {
      // Animate floating skill icons
      skillIconsRef.current.forEach((icon, index) => {
        if (icon) {
          gsap.set(icon, {
            x: Math.random() * 200 - 100,
            y: Math.random() * 200 - 100,
          });
          
          gsap.to(icon, {
            y: "+=30",
            duration: 2 + index * 0.5,
            repeat: -1,
            yoyo: true,
            ease: "power2.inOut",
          });
        }
      });
    }, heroRef);

    return () => ctx.revert();
  }, []);

  const skillIcons = [Users, BookOpen, TrendingUp, Zap];

  return (
    <section 
      ref={heroRef}
      className="relative min-h-screen flex items-center justify-center overflow-hidden bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 dark:from-gray-900 dark:via-blue-900 dark:to-purple-900"
    >
      {/* Floating Skill Icons */}
      <div className="absolute inset-0 pointer-events-none">
        {skillIcons.map((Icon, index) => (
          <div
            key={index}
            ref={(el) => { skillIconsRef.current[index] = el; }}
            className="absolute opacity-10 dark:opacity-20"
            style={{
              left: `${20 + index * 20}%`,
              top: `${20 + index * 15}%`,
            }}
          >
            <Icon size={48} className="text-blue-600" />
          </div>
        ))}
      </div>

      <div className="container mx-auto px-4 z-10">
        <div className="text-center max-w-4xl mx-auto">
          {/* Animated Headline */}
          <motion.div
            initial={{ opacity: 0, y: 50 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.2 }}
          >
            <h1 className="text-4xl md:text-6xl lg:text-7xl font-bold mb-6 leading-tight">
              <motion.span
                initial={{ opacity: 0, x: -50 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.8, delay: 0.5 }}
                className="text-gray-900 dark:text-white"
              >
                Exchange Skills,
              </motion.span>
              <br />
              <motion.span
                initial={{ opacity: 0, x: 50 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.8, delay: 0.7 }}
                className="bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent"
              >
                Grow Together
              </motion.span>
            </h1>
          </motion.div>

          <motion.p
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.9 }}
            className="text-xl md:text-2xl text-gray-600 dark:text-gray-300 mb-8 max-w-2xl mx-auto"
          >
            Connect with professionals worldwide. Trade your expertise for new skills. 
            Build meaningful relationships while advancing your career.
          </motion.p>

          {/* CTA Buttons */}
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 1.1 }}
            className="flex flex-col sm:flex-row gap-4 justify-center items-center"
          >
            <motion.div
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              <Button 
                size="lg" 
                className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white px-8 py-4 text-lg group"
              >
                Start Trading Skills
                <ArrowRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
              </Button>
            </motion.div>
            
            <motion.div
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              <Button 
                variant="outline" 
                size="lg"
                className="px-8 py-4 text-lg border-2 hover:bg-gray-50 dark:hover:bg-gray-800"
              >
                Watch Demo
              </Button>
            </motion.div>
          </motion.div>

          {/* Stats */}
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 1.3 }}
            className="mt-16 grid grid-cols-3 gap-8 max-w-md mx-auto"
          >
            {[
              { number: "10K+", label: "Active Users" },
              { number: "500+", label: "Skills Available" },
              { number: "95%", label: "Success Rate" }
            ].map((stat, index) => (
              <div key={index} className="text-center">
                <motion.div
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  transition={{ duration: 0.5, delay: 1.5 + index * 0.1 }}
                  className="text-2xl md:text-3xl font-bold text-blue-600 dark:text-blue-400"
                >
                  {stat.number}
                </motion.div>
                <div className="text-sm text-gray-600 dark:text-gray-400">
                  {stat.label}
                </div>
              </div>
            ))}
          </motion.div>
        </div>
      </div>

      {/* Scroll Indicator */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 1, delay: 2 }}
        className="absolute bottom-8 left-1/2 transform -translate-x-1/2"
      >
        <motion.div
          animate={{ y: [0, 10, 0] }}
          transition={{ duration: 2, repeat: Infinity }}
          className="w-6 h-10 border-2 border-gray-400 dark:border-gray-600 rounded-full p-1"
        >
          <div className="w-1 h-3 bg-gray-400 dark:bg-gray-600 rounded-full mx-auto"></div>
        </motion.div>
      </motion.div>
    </section>
  );
};

export default HeroSection;