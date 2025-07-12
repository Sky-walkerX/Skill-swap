import { useState, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { gsap } from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { 
  ArrowRight, 
  CheckCircle, 
  Users, 
  Star, 
  Globe,
  Sparkles,
  Mail,
  Lock
} from 'lucide-react';

gsap.registerPlugin(ScrollTrigger);

const CallToActionSection = () => {
  const sectionRef = useRef(null);
  const [email, setEmail] = useState('');
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const socialProofStats = [
    { icon: Users, value: '10,000+', label: 'Active Members' },
    { icon: Star, value: '4.9/5', label: 'User Rating' },
    { icon: Globe, value: '50+', label: 'Countries' }
  ];

  const benefits = [
    'Free to join and start exchanging',
    'Verified skill assessments',
    '24/7 community support',
    'Earn certificates and badges'
  ];

  useEffect(() => {
    const ctx = gsap.context(() => {
      // Floating animation for the form
      gsap.to('.floating-form', {
        y: -20,
        duration: 3,
        repeat: -1,
        yoyo: true,
        ease: "power2.inOut"
      });

      // Parallax effect for background elements
      gsap.to('.cta-bg-element', {
        y: -100,
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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    setIsSubmitted(true);
    setIsLoading(false);
  };

  return (
    <section 
      ref={sectionRef}
      className="relative py-20 bg-gradient-to-br from-blue-900 via-purple-900 to-indigo-900 overflow-hidden"
    >
      {/* Background Elements */}
      <div className="absolute inset-0 pointer-events-none">
        <div className="cta-bg-element absolute top-20 left-10 w-32 h-32 bg-blue-500 rounded-full opacity-10 blur-xl"></div>
        <div className="cta-bg-element absolute top-40 right-20 w-24 h-24 bg-purple-500 rounded-full opacity-10 blur-xl"></div>
        <div className="cta-bg-element absolute bottom-20 left-1/3 w-28 h-28 bg-pink-500 rounded-full opacity-10 blur-xl"></div>
      </div>

      <div className="container mx-auto px-4 relative z-10">
        <div className="max-w-6xl mx-auto">
          <div className="grid lg:grid-cols-2 gap-12 items-center">
            
            {/* Left Content */}
            <motion.div
              initial={{ opacity: 0, x: -50 }}
              whileInView={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.8 }}
              viewport={{ once: true }}
              className="text-white"
            >
              <div className="mb-6">
                <Badge className="bg-white/10 text-white border-white/20 mb-4">
                  <Sparkles className="mr-1 h-3 w-3" />
                  Join Today
                </Badge>
                <h2 className="text-3xl md:text-5xl font-bold mb-6 leading-tight">
                  Ready to Transform
                  <br />
                  Your Career?
                </h2>
                <p className="text-xl text-blue-100 mb-8 leading-relaxed">
                  Join thousands of professionals who are already growing their skills and 
                  expanding their networks through meaningful exchanges.
                </p>
              </div>

              {/* Benefits List */}
              <div className="space-y-4 mb-8">
                {benefits.map((benefit, index) => (
                  <motion.div
                    key={index}
                    initial={{ opacity: 0, x: -20 }}
                    whileInView={{ opacity: 1, x: 0 }}
                    transition={{ duration: 0.5, delay: index * 0.1 }}
                    viewport={{ once: true }}
                    className="flex items-center"
                  >
                    <CheckCircle className="h-5 w-5 text-green-400 mr-3 flex-shrink-0" />
                    <span className="text-blue-100">{benefit}</span>
                  </motion.div>
                ))}
              </div>

              {/* Social Proof */}
              <div className="grid grid-cols-3 gap-6">
                {socialProofStats.map((stat, index) => (
                  <motion.div
                    key={index}
                    initial={{ opacity: 0, y: 20 }}
                    whileInView={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.5, delay: 0.3 + index * 0.1 }}
                    viewport={{ once: true }}
                    className="text-center"
                  >
                    <div className="inline-flex p-2 bg-white/10 rounded-lg mb-2">
                      <stat.icon className="h-5 w-5 text-blue-300" />
                    </div>
                    <div className="text-2xl font-bold text-white">{stat.value}</div>
                    <div className="text-sm text-blue-200">{stat.label}</div>
                  </motion.div>
                ))}
              </div>
            </motion.div>

            {/* Right Content - Floating Form */}
            <motion.div
              initial={{ opacity: 0, x: 50 }}
              whileInView={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.8, delay: 0.2 }}
              viewport={{ once: true }}
              className="floating-form"
            >
              <Card className="bg-white/95 backdrop-blur-lg border-0 shadow-2xl">
                <CardContent className="p-8">
                  <AnimatePresence mode="wait">
                    {!isSubmitted ? (
                      <motion.div
                        key="form"
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        transition={{ duration: 0.3 }}
                      >
                        <div className="text-center mb-6">
                          <h3 className="text-2xl font-bold text-gray-900 mb-2">
                            Start Your Journey
                          </h3>
                          <p className="text-gray-600">
                            Create your free account and get matched with your first skill exchange partner
                          </p>
                        </div>

                        <form onSubmit={handleSubmit} className="space-y-4">
                          <div className="relative">
                            <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                            <Input
                              type="email"
                              placeholder="Enter your email address"
                              value={email}
                              onChange={(e) => setEmail(e.target.value)}
                              required
                              className="pl-10 h-12"
                            />
                          </div>
                          
                          <Button
                            type="submit"
                            disabled={isLoading}
                            className="w-full h-12 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white font-semibold"
                          >
                            {isLoading ? (
                              <div className="flex items-center">
                                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                                Creating Account...
                              </div>
                            ) : (
                              <div className="flex items-center">
                                Get Started Free
                                <ArrowRight className="ml-2 h-4 w-4" />
                              </div>
                            )}
                          </Button>
                        </form>

                        <div className="mt-6 text-center">
                          <div className="flex items-center justify-center text-sm text-gray-500 mb-2">
                            <Lock className="h-3 w-3 mr-1" />
                            100% secure and private
                          </div>
                          <p className="text-xs text-gray-400">
                            By signing up, you agree to our Terms of Service and Privacy Policy
                          </p>
                        </div>
                      </motion.div>
                    ) : (
                      <motion.div
                        key="success"
                        initial={{ opacity: 0, scale: 0.8 }}
                        animate={{ opacity: 1, scale: 1 }}
                        exit={{ opacity: 0, scale: 0.8 }}
                        transition={{ duration: 0.5 }}
                        className="text-center py-8"
                      >
                        <motion.div
                          initial={{ scale: 0 }}
                          animate={{ scale: 1 }}
                          transition={{ duration: 0.5, delay: 0.2 }}
                          className="w-16 h-16 bg-green-500 rounded-full flex items-center justify-center mx-auto mb-4"
                        >
                          <CheckCircle className="h-8 w-8 text-white" />
                        </motion.div>
                        
                        <h3 className="text-2xl font-bold text-gray-900 mb-2">
                          Welcome to SkillShare Connect!
                        </h3>
                        <p className="text-gray-600 mb-6">
                          Check your email to complete your account setup and start connecting with skill exchange partners.
                        </p>
                        
                        <Button
                          className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white px-8"
                        >
                          Go to Dashboard
                        </Button>
                      </motion.div>
                    )}
                  </AnimatePresence>
                </CardContent>
              </Card>
            </motion.div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default CallToActionSection;