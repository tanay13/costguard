import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'CostGuard - Autonomous Cloud Cost Optimization',
  description: 'AI-driven cloud cost optimization with automated PRs',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}

