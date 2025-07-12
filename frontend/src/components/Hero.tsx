'use client';

import { useEffect, useRef } from 'react';
import { motion } from 'framer-motion';
import { gsap } from 'gsap';
import { Button } from '@/components/ui/button';
import { ArrowRight, Users, BookOpen, TrendingUp, Zap, Sparkles } from 'lucide-react';

const HeroSection = () => {
	const heroRef = useRef(null);
	const skillIconsRef = useRef<(HTMLDivElement | null)[]>([]);

	useEffect(() => {
		const ctx = gsap.context(() => {
			// Animate floating skill icons with more organic movement
			skillIconsRef.current.forEach((icon, index) => {
				if (icon) {
					gsap.set(icon, {
						x: Math.random() * 300 - 150,
						y: Math.random() * 300 - 150
					});

					gsap.to(icon, {
						y: '+=40',
						x: '+=20',
						rotation: '+=15',
						duration: 3 + index * 0.7,
						repeat: -1,
						yoyo: true,
						ease: 'power2.inOut'
					});
				}
			});

			// Parallax background elements
			gsap.to('.bg-shape', {
				y: -50,
				rotation: 10,
				ease: 'none',
				scrollTrigger: {
					trigger: heroRef.current,
					start: 'top bottom',
					end: 'bottom top',
					scrub: true
				}
			});
		}, heroRef);

		return () => ctx.revert();
	}, []);

	const skillIcons = [Users, BookOpen, TrendingUp, Zap];

	return (
		<section
			ref={heroRef}
			className="relative min-h-screen flex items-center justify-center overflow-hidden bg-gradient-to-br from-stone-50 via-amber-50 to-orange-50 dark:from-stone-900 dark:via-amber-900 dark:to-orange-900"
		>
			{/* Organic Background Shapes */}
			<div className="absolute inset-0 pointer-events-none">
				<div className="bg-shape absolute top-20 left-10 w-64 h-64 bg-gradient-to-br from-sage-200 to-sage-300 rounded-full opacity-20 blur-3xl"></div>
				<div className="bg-shape absolute top-40 right-20 w-48 h-48 bg-gradient-to-br from-terracotta-200 to-terracotta-300 rounded-full opacity-20 blur-2xl"></div>
				<div className="bg-shape absolute bottom-20 left-1/3 w-56 h-56 bg-gradient-to-br from-amber-200 to-amber-300 rounded-full opacity-20 blur-3xl"></div>

				{/* Floating Skill Icons */}
				{skillIcons.map((Icon, index) => (
					<div
						key={index}
						ref={(el) => {
							skillIconsRef.current[index] = el;
						}}
						className="absolute opacity-15 dark:opacity-25"
						style={{
							left: `${15 + index * 20}%`,
							top: `${15 + index * 20}%`
						}}
					>
						<div className="p-4 bg-white/30 backdrop-blur-sm rounded-2xl border border-white/20">
							<Icon size={32} className="text-sage-600" />
						</div>
					</div>
				))}
			</div>

			<div className="container mx-auto px-4 z-10">
				<div className="text-center max-w-5xl mx-auto">
					{/* Animated Headline */}
					<motion.div
						initial={{ opacity: 0, y: 50 }}
						animate={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.8, delay: 0.2 }}
					>
						<div className="inline-flex items-center gap-2 bg-sage-100 dark:bg-sage-800 px-4 py-2 rounded-full mb-6">
							<Sparkles className="h-4 w-4 text-sage-600" />
							<span className="text-sage-700 dark:text-sage-300 font-medium">
								Where Skills Meet Opportunity
							</span>
						</div>

						<h1 className="text-5xl md:text-7xl lg:text-8xl font-bold mb-8 leading-tight">
							<motion.span
								initial={{ opacity: 0, x: -50 }}
								animate={{ opacity: 1, x: 0 }}
								transition={{ duration: 0.8, delay: 0.5 }}
								className="text-stone-900 dark:text-stone-100 block"
							>
								Trade Skills,
							</motion.span>
							<motion.span
								initial={{ opacity: 0, x: 50 }}
								animate={{ opacity: 1, x: 0 }}
								transition={{ duration: 0.8, delay: 0.7 }}
								className="bg-gradient-to-r from-sage-600 via-terracotta-500 to-amber-600 bg-clip-text text-transparent block"
							>
								Build Futures
							</motion.span>
						</h1>
					</motion.div>

					<motion.p
						initial={{ opacity: 0, y: 30 }}
						animate={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.8, delay: 0.9 }}
						className="text-xl md:text-2xl text-stone-600 dark:text-stone-300 mb-10 max-w-3xl mx-auto leading-relaxed"
					>
						Join a community where knowledge flows freely. Exchange your expertise,
						learn new skills, and forge meaningful professional connections.
					</motion.p>

					{/* CTA Buttons */}
					<motion.div
						initial={{ opacity: 0, y: 30 }}
						animate={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.8, delay: 1.1 }}
						className="flex flex-col sm:flex-row gap-6 justify-center items-center mb-16"
					>
						<motion.div whileHover={{ scale: 1.05, y: -2 }} whileTap={{ scale: 0.95 }}>
							<Button
								size="lg"
								className="bg-gradient-to-r from-sage-600 to-terracotta-600 hover:from-sage-700 hover:to-terracotta-700 text-black hover:text-white px-10 py-6 text-lg rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 group"
							>
								Start Your Journey
								<ArrowRight className="ml-3 h-5 w-5 group-hover:translate-x-1 transition-transform" />
							</Button>
						</motion.div>

						<motion.div whileHover={{ scale: 1.05, y: -2 }} whileTap={{ scale: 0.95 }}>
							<Button
								variant="outline"
								size="lg"
								className="px-10 py-6 text-lg border-2 border-stone-300 hover:bg-stone-50 dark:border-stone-600 dark:hover:bg-stone-800 rounded-2xl transition-all duration-300 bg-transparent"
							>
								Explore Skills
							</Button>
						</motion.div>
					</motion.div>

					{/* Enhanced Stats */}
					<motion.div
						initial={{ opacity: 0, y: 30 }}
						animate={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.8, delay: 1.3 }}
						className="grid grid-cols-3 gap-8 max-w-2xl mx-auto"
					>
						{[
							{ number: '15K+', label: 'Skill Traders', icon: Users },
							{ number: '800+', label: 'Skills Available', icon: BookOpen },
							{ number: '98%', label: 'Success Rate', icon: TrendingUp }
						].map((stat, index) => (
							<motion.div
								key={index}
								initial={{ scale: 0 }}
								animate={{ scale: 1 }}
								transition={{ duration: 0.5, delay: 1.5 + index * 0.1 }}
								className="text-center group"
							>
								<div className="inline-flex p-3 bg-white/60 backdrop-blur-sm rounded-2xl mb-3 group-hover:scale-110 transition-transform duration-300">
									<stat.icon className="h-6 w-6 text-sage-600" />
								</div>
								<div className="text-3xl md:text-4xl font-bold text-stone-900 dark:text-stone-100 mb-1">
									{stat.number}
								</div>
								<div className="text-sm text-stone-600 dark:text-stone-400 font-medium">
									{stat.label}
								</div>
							</motion.div>
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
					animate={{ y: [0, 12, 0] }}
					transition={{ duration: 2, repeat: Number.POSITIVE_INFINITY }}
					className="w-6 h-10 border-2 border-stone-400 dark:border-stone-600 rounded-full p-1"
				>
					<div className="w-1 h-3 bg-sage-500 rounded-full mx-auto"></div>
				</motion.div>
			</motion.div>
		</section>
	);
};

export default HeroSection;
