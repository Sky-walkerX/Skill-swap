import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
	images: {
		remotePatterns: [
			{
				protocol: 'https',
				hostname: 'media.gettyimages.com',
				pathname: '/id/**'
			}
		]
	},
	// Allow configuration of development server port
	...(process.env.NODE_ENV === 'development' && {
		experimental: {
			serverComponentsExternalPackages: []
		}
	})
};

export default nextConfig;
