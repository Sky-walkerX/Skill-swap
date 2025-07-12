"use client"

import { motion } from "framer-motion"
import { Users, BookOpen, TrendingUp, Zap, Heart } from "lucide-react"

const LoadingAnimation = () => {
  const icons = [Users, BookOpen, TrendingUp, Zap]

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-gradient-to-br from-stone-50 via-sage-50 to-terracotta-50 dark:from-stone-900 dark:via-sage-900 dark:to-terracotta-900">
      <div className="text-center">
        {/* Enhanced Logo Animation */}
        <motion.div
          initial={{ scale: 0, rotate: 0 }}
          animate={{ scale: 1, rotate: 360 }}
          transition={{ duration: 1.2, ease: "easeOut" }}
          className="mb-10"
        >
          <div className="w-24 h-24 bg-gradient-to-r from-sage-600 to-terracotta-600 rounded-3xl flex items-center justify-center mx-auto shadow-2xl">
            <Heart className="h-12 w-12 text-white" />
          </div>
        </motion.div>

        {/* Enhanced Floating Icons */}
        <div className="relative w-40 h-40 mx-auto mb-10">
          {icons.map((Icon, index) => (
            <motion.div
              key={index}
              className="absolute"
              style={{
                left: `${50 + 45 * Math.cos((index * Math.PI) / 2)}%`,
                top: `${50 + 45 * Math.sin((index * Math.PI) / 2)}%`,
              }}
              initial={{ opacity: 0, scale: 0 }}
              animate={{
                opacity: 1,
                scale: 1,
                rotate: 360,
                x: [0, 15, 0, -15, 0],
                y: [0, -15, 0, 15, 0],
              }}
              transition={{
                duration: 3,
                delay: index * 0.3,
                repeat: Number.POSITIVE_INFINITY,
                repeatType: "reverse",
              }}
            >
              <div className="w-10 h-10 bg-white dark:bg-stone-800 rounded-2xl shadow-lg flex items-center justify-center border-2 border-sage-200 dark:border-sage-700">
                <Icon className="h-5 w-5 text-sage-600" />
              </div>
            </motion.div>
          ))}
        </div>

        {/* Enhanced Loading Text */}
        <motion.h1
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.5 }}
          className="text-3xl md:text-4xl font-bold text-stone-900 dark:text-stone-100 mb-4"
        >
          SkillShare Connect
        </motion.h1>

        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.8, delay: 0.7 }}
          className="text-stone-600 dark:text-stone-300 mb-10 text-lg"
        >
          Where skills meet opportunity...
        </motion.p>

        {/* Enhanced Progress Bar */}
        <div className="w-80 h-3 bg-stone-200 dark:bg-stone-700 rounded-full overflow-hidden mx-auto">
          <motion.div
            className="h-full bg-gradient-to-r from-sage-600 to-terracotta-600 rounded-full"
            initial={{ width: 0 }}
            animate={{ width: "100%" }}
            transition={{ duration: 2.5, ease: "easeInOut" }}
          />
        </div>
      </div>
    </div>
  )
}

export default LoadingAnimation
