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
	}
};

export default nextConfig;
