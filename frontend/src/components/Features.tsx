import { useEffect, useRef } from 'react';
import { motion } from 'framer-motion';
import { gsap } from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import { useInView } from 'react-intersection-observer';
import { Card, CardContent } from '@/components/ui/card';
import { Users, Shield, Zap, Target, Star, Award } from 'lucide-react';

gsap.registerPlugin(ScrollTrigger);

const FeaturesSection = () => {
  const sectionRef = useRef(null);
  const [ref, inView] = useInView({
    threshold: 0.1,
    triggerOnce: true,
  });

  useEffect(() => {
    const ctx = gsap.context(() => {
      const cards = gsap.utils.toArray('.feature-card');
      
      cards.forEach((card: any, index) => {
        gsap.fromTo(card, 
          { y: 100, opacity: 0 },
          {
            y: 0,
            opacity: 1,
            duration: 0.8,
            delay: index * 0.2,
            scrollTrigger: {
              trigger: card,
              start: "top 80%",
              end: "bottom 20%",
              toggleActions: "play none none reverse"
            }
          }
        );
      });

      // Parallax effect for the section
      gsap.to('.features-bg', {
        yPercent: -50,
        ease: "none",
        scrollTrigger: {
          trigger: sectionRef.current,
          start: "top bottom",
          end: "bottom top",
          scrub: true
        }
      });
    }, sectionRef);

    return () => ctx.revert();
  }, []);

  const features = [
    {
      icon: Users,
      title: "Global Community",
      description: "Connect with skilled professionals from around the world and expand your network.",
      stats: "50+ Countries"
    },
    {
      icon: Shield,
      title: "Verified Skills",
      description: "Our verification system ensures quality interactions and trustworthy skill exchanges.",
      stats: "95% Verified"
    },
    {
      icon: Zap,
      title: "Instant Matching",
      description: "Advanced AI algorithm matches you with the perfect skill exchange partners in seconds.",
      stats: "< 2 sec"
    },
    {
      icon: Target,
      title: "Goal Tracking",
      description: "Set learning objectives and track your progress with detailed analytics and insights.",
      stats: "100% Tracked"
    },
    {
      icon: Star,
      title: "Quality Assured",
      description: "Rating system and reviews ensure high-quality skill exchanges and learning experiences.",
      stats: "4.9/5 Rating"
    },
    {
      icon: Award,
      title: "Skill Badges",
      description: "Earn digital badges and certificates as you master new skills and help others learn.",
      stats: "500+ Badges"
    }
  ];

  const testimonials = [
    {
      name: "Sarah Chen",
      role: "UX Designer",
      text: "I exchanged my design skills for coding lessons. Amazing community!",
      avatar: "https://images.pexels.com/photos/415829/pexels-photo-415829.jpeg?w=150"
    },
    {
      name: "Marcus Johnson",
      role: "Data Scientist",
      text: "Found the perfect mentor for machine learning in just 2 days.",
      avatar: "https://images.pexels.com/photos/2379004/pexels-photo-2379004.jpeg?w=150"
    },
    {
      name: "Elena Rodriguez",
      role: "Marketing Expert",
      text: "The skill verification system gave me confidence in every exchange.",
      avatar: "https://images.pexels.com/photos/774909/pexels-photo-774909.jpeg?w=150"
    }
  ];

  return (
    <section ref={sectionRef} className="relative py-20 bg-white dark:bg-gray-900 overflow-hidden">
      {/* Background Pattern */}
      <div className="features-bg absolute inset-0 opacity-5 dark:opacity-10">
        <div className="absolute inset-0 bg-gradient-to-r from-blue-600 to-purple-600"></div>
      </div>

      <div ref={ref} className="container mx-auto px-4 relative z-10">
        {/* Section Header */}
        <motion.div
          initial={{ opacity: 0, y: 50 }}
          animate={inView ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.8 }}
          className="text-center mb-16"
        >
          <h2 className="text-3xl md:text-5xl font-bold mb-6 text-gray-900 dark:text-white">
            Why Choose SkillShare Connect?
          </h2>
          <p className="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
            Join thousands of professionals who are already growing their careers through skill exchange
          </p>
        </motion.div>

        {/* Features Grid */}
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8 mb-20">
          {features.map((feature, index) => (
            <motion.div
              key={index}
              className="feature-card"
              whileHover={{ y: -10, scale: 1.02 }}
              transition={{ duration: 0.3 }}
            >
              <Card className="h-full border-0 shadow-lg hover:shadow-xl transition-all duration-300 bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900">
                <CardContent className="p-8">
                  <div className="flex items-center mb-4">
                    <div className="p-3 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg mr-4">
                      <feature.icon className="h-6 w-6 text-white" />
                    </div>
                    <div className="text-sm font-semibold text-blue-600 dark:text-blue-400">
                      {feature.stats}
                    </div>
                  </div>
                  <h3 className="text-xl font-bold mb-3 text-gray-900 dark:text-white">
                    {feature.title}
                  </h3>
                  <p className="text-gray-600 dark:text-gray-300 leading-relaxed">
                    {feature.description}
                  </p>
                </CardContent>
              </Card>
            </motion.div>
          ))}
        </div>

        {/* Testimonials */}
        <motion.div
          initial={{ opacity: 0, y: 50 }}
          animate={inView ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.8, delay: 0.5 }}
          className="text-center"
        >
          <h3 className="text-2xl md:text-3xl font-bold mb-12 text-gray-900 dark:text-white">
            What Our Community Says
          </h3>
          
          <div className="grid md:grid-cols-3 gap-8">
            {testimonials.map((testimonial, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, scale: 0.9 }}
                animate={inView ? { opacity: 1, scale: 1 } : {}}
                transition={{ duration: 0.5, delay: 0.7 + index * 0.2 }}
                whileHover={{ scale: 1.05 }}
              >
                <Card className="p-6 border-0 shadow-lg hover:shadow-xl transition-all duration-300">
                  <CardContent className="p-0 text-center">
                    <img
                      src={testimonial.avatar}
                      alt={testimonial.name}
                      className="w-16 h-16 rounded-full mx-auto mb-4 object-cover"
                    />
                    <p className="text-gray-600 dark:text-gray-300 mb-4 italic">
                      "{testimonial.text}"
                    </p>
                    <div>
                      <div className="font-semibold text-gray-900 dark:text-white">
                        {testimonial.name}
                      </div>
                      <div className="text-sm text-gray-500 dark:text-gray-400">
                        {testimonial.role}
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </motion.div>
            ))}
          </div>
        </motion.div>
      </div>
    </section>
  );
};

export default FeaturesSection;