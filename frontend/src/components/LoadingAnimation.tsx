import { motion } from 'framer-motion';
import { Users, BookOpen, TrendingUp, Zap } from 'lucide-react';

const LoadingAnimation = () => {
  const icons = [Users, BookOpen, TrendingUp, Zap];

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 dark:from-gray-900 dark:via-blue-900 dark:to-purple-900">
      <div className="text-center">
        {/* Logo Animation */}
        <motion.div
          initial={{ scale: 0, rotate: 0 }}
          animate={{ scale: 1, rotate: 360 }}
          transition={{ duration: 1, ease: "easeOut" }}
          className="mb-8"
        >
          <div className="w-20 h-20 bg-gradient-to-r from-blue-600 to-purple-600 rounded-2xl flex items-center justify-center mx-auto">
            <Users className="h-10 w-10 text-white" />
          </div>
        </motion.div>

        {/* Floating Icons */}
        <div className="relative w-32 h-32 mx-auto mb-8">
          {icons.map((Icon, index) => (
            <motion.div
              key={index}
              className="absolute"
              style={{
                left: `${50 + 40 * Math.cos(index * Math.PI / 2)}%`,
                top: `${50 + 40 * Math.sin(index * Math.PI / 2)}%`,
              }}
              initial={{ opacity: 0, scale: 0 }}
              animate={{ 
                opacity: 1, 
                scale: 1,
                rotate: 360,
                x: [0, 10, 0, -10, 0],
                y: [0, -10, 0, 10, 0]
              }}
              transition={{ 
                duration: 2,
                delay: index * 0.2,
                repeat: Infinity,
                repeatType: "reverse"
              }}
            >
              <div className="w-8 h-8 bg-white dark:bg-gray-800 rounded-lg shadow-lg flex items-center justify-center">
                <Icon className="h-4 w-4 text-blue-600" />
              </div>
            </motion.div>
          ))}
        </div>

        {/* Loading Text */}
        <motion.h1
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.5 }}
          className="text-2xl md:text-3xl font-bold text-gray-900 dark:text-white mb-4"
        >
          SkillShare Connect
        </motion.h1>

        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.8, delay: 0.7 }}
          className="text-gray-600 dark:text-gray-300 mb-8"
        >
          Connecting skills, creating opportunities...
        </motion.p>

        {/* Progress Bar */}
        <div className="w-64 h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden mx-auto">
          <motion.div
            className="h-full bg-gradient-to-r from-blue-600 to-purple-600"
            initial={{ width: 0 }}
            animate={{ width: "100%" }}
            transition={{ duration: 2, ease: "easeInOut" }}
          />
        </div>
      </div>
    </div>
  );
};

export default LoadingAnimation;