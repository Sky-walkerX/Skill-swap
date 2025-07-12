import { useEffect, useRef, useState } from 'react';
import { motion } from 'framer-motion';
import { gsap } from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { 
  Code, 
  Palette, 
  BarChart, 
  Megaphone, 
  Camera, 
  Music, 
  Languages, 
  Wrench,
  Search,
  Filter,
  TrendingUp
} from 'lucide-react';

gsap.registerPlugin(ScrollTrigger);

const SkillsMarketplace = () => {
  const sectionRef = useRef(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('All');
  const [skillCount, setSkillCount] = useState(0);

  const skillCategories = [
    { icon: Code, name: 'Programming', count: 124, color: 'bg-blue-500' },
    { icon: Palette, name: 'Design', count: 89, color: 'bg-purple-500' },
    { icon: BarChart, name: 'Analytics', count: 67, color: 'bg-green-500' },
    { icon: Megaphone, name: 'Marketing', count: 78, color: 'bg-red-500' },
    { icon: Camera, name: 'Photography', count: 45, color: 'bg-yellow-500' },
    { icon: Music, name: 'Music', count: 34, color: 'bg-pink-500' },
    { icon: Languages, name: 'Languages', count: 56, color: 'bg-indigo-500' },
    { icon: Wrench, name: 'Engineering', count: 92, color: 'bg-orange-500' }
  ];

  const categories = ['All', ...skillCategories.map(cat => cat.name)];
  const totalSkills = skillCategories.reduce((sum, cat) => sum + cat.count, 0);

  const filteredCategories = selectedCategory === 'All' 
    ? skillCategories 
    : skillCategories.filter(cat => cat.name === selectedCategory);

  useEffect(() => {
    // Animate skill counter
    gsap.to({ count: 0 }, {
      count: totalSkills,
      duration: 2,
      ease: "power2.out",
      scrollTrigger: {
        trigger: sectionRef.current,
        start: "top 80%",
        onEnter: () => {
          gsap.to({ count: 0 }, {
            count: totalSkills,
            duration: 2,
            ease: "power2.out",
            onUpdate: function() {
              setSkillCount(Math.round(this.targets()[0].count));
            }
          });
        }
      }
    });

    const ctx = gsap.context(() => {
      // Animate category cards
      gsap.fromTo('.skill-category',
        { y: 50, opacity: 0, scale: 0.9 },
        {
          y: 0,
          opacity: 1,
          scale: 1,
          duration: 0.6,
          stagger: 0.1,
          scrollTrigger: {
            trigger: '.skills-grid',
            start: "top 80%",
            end: "bottom 20%",
            toggleActions: "play none none reverse"
          }
        }
      );
    }, sectionRef);

    return () => ctx.revert();
  }, [totalSkills]);

  return (
    <section ref={sectionRef} className="py-20 bg-white dark:bg-gray-900">
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
            Skills Marketplace
          </h2>
          <p className="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto mb-8">
            Explore our diverse range of skills available for exchange. Find exactly what you're looking for.
          </p>
          
          {/* Real-time Counter */}
          <motion.div
            initial={{ scale: 0 }}
            whileInView={{ scale: 1 }}
            transition={{ duration: 0.8, delay: 0.3 }}
            viewport={{ once: true }}
            className="inline-flex items-center gap-2 bg-gradient-to-r from-blue-600 to-purple-600 text-white px-6 py-3 rounded-full text-lg font-semibold"
          >
            <TrendingUp className="h-5 w-5" />
            <span className="font-mono">{skillCount}</span>
            <span>Skills Available</span>
          </motion.div>
        </motion.div>

        {/* Search and Filter */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.2 }}
          viewport={{ once: true }}
          className="max-w-2xl mx-auto mb-12"
        >
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
              <Input
                placeholder="Search skills..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10 h-12"
              />
            </div>
            <Button
              variant="outline"
              className="h-12 px-6 border-2 hover:bg-gray-50 dark:hover:bg-gray-800"
            >
              <Filter className="mr-2 h-4 w-4" />
              Advanced Filter
            </Button>
          </div>
        </motion.div>

        {/* Category Filter */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.3 }}
          viewport={{ once: true }}
          className="mb-12"
        >
          <div className="flex flex-wrap justify-center gap-2">
            {categories.map((category) => (
              <Button
                key={category}
                variant={selectedCategory === category ? "default" : "outline"}
                onClick={() => setSelectedCategory(category)}
                className={`transition-all duration-300 ${
                  selectedCategory === category
                    ? 'bg-gradient-to-r from-blue-600 to-purple-600 text-white'
                    : 'hover:bg-gray-50 dark:hover:bg-gray-800'
                }`}
              >
                {category}
              </Button>
            ))}
          </div>
        </motion.div>

        {/* Skills Grid */}
        <div className="skills-grid grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {filteredCategories.map((skill, index) => (
            <motion.div
              key={skill.name}
              className="skill-category"
              whileHover={{ y: -10, scale: 1.05 }}
              transition={{ duration: 0.3 }}
            >
              <Card className="h-full border-0 shadow-lg hover:shadow-2xl transition-all duration-300 cursor-pointer group">
                <CardContent className="p-6 text-center">
                  <div className={`inline-flex p-4 rounded-2xl ${skill.color} mb-4 group-hover:scale-110 transition-transform duration-300`}>
                    <skill.icon className="h-8 w-8 text-white" />
                  </div>
                  
                  <h3 className="text-xl font-bold mb-2 text-gray-900 dark:text-white">
                    {skill.name}
                  </h3>
                  
                  <p className="text-gray-600 dark:text-gray-300 mb-4">
                    {skill.count} available experts
                  </p>
                  
                  <Badge 
                    variant="secondary" 
                    className="bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-300"
                  >
                    High Demand
                  </Badge>
                </CardContent>
              </Card>
            </motion.div>
          ))}
        </div>

        {/* Popular Skills Preview */}
        <motion.div
          initial={{ opacity: 0, y: 50 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.5 }}
          viewport={{ once: true }}
          className="mt-20 text-center"
        >
          <h3 className="text-2xl md:text-3xl font-bold mb-8 text-gray-900 dark:text-white">
            Most Popular This Week
          </h3>
          
          <div className="flex flex-wrap justify-center gap-3 mb-8">
            {[
              'React.js', 'UI/UX Design', 'Python', 'Digital Marketing', 
              'Photography', 'Data Analysis', 'Spanish', 'Video Editing'
            ].map((skill) => (
              <motion.div
                key={skill}
                whileHover={{ scale: 1.1 }}
                whileTap={{ scale: 0.95 }}
              >
                <Badge 
                  variant="outline" 
                  className="px-4 py-2 text-lg border-2 hover:bg-blue-50 dark:hover:bg-blue-900 cursor-pointer transition-colors"
                >
                  {skill}
                </Badge>
              </motion.div>
            ))}
          </div>
          
          <Button 
            size="lg"
            className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white px-8 py-4 text-lg"
          >
            Browse All Skills
          </Button>
        </motion.div>
      </div>
    </section>
  );
};

export default SkillsMarketplace;