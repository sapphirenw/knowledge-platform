/** @type {import('next').NextConfig} */
const nextConfig = {
    output: "standalone",
    experimental: {
        instrumentationHook: true,
    },
    transpilePackages: ['lucide-react'],
};

export default nextConfig;
