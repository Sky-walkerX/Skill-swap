import type { Metadata } from 'next';
import { Geist, Geist_Mono } from 'next/font/google';
import { QueryClientProvider } from '@/providers/QueryClientProvider';
import { Navigation } from '@/components/Navigation';
import NextAuthSessionProvider from '@/components/AuthProvider';
import './globals.css';
import { SmoothScroll } from '@/components/smooth-scroll';

const geistSans = Geist({
	variable: '--font-geist-sans',
	subsets: ['latin']
});

const geistMono = Geist_Mono({
	variable: '--font-geist-mono',
	subsets: ['latin']
});

export const metadata: Metadata = {
	title: 'SkillShare - Exchange Skills with Others',
	description: 'Connect with people in your area to exchange skills and knowledge'
};

export default function RootLayout({
	children
}: Readonly<{
	children: React.ReactNode;
}>) {
	return (
		<html lang="en">
			<body
				className={`${geistSans.variable} ${geistMono.variable} antialiased bg-background`}
			>
				<SmoothScroll />
        <NextAuthSessionProvider>
				  <QueryClientProvider>
					  <Navigation />
					  <main className="min-h-screen">{children}</main>
				  </QueryClientProvider>
				</NextAuthSessionProvider>
			</body>
		</html>
	);
}
