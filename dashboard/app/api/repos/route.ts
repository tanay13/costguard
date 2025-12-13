import { NextResponse } from 'next/server';
import fs from 'fs';
import path from 'path';

export async function GET() {
  try {
    // Use /tmp directory for Vercel
    const isVercel = process.env.VERCEL === '1';
    const dataDir = isVercel ? '/tmp/costguard-data' : path.join(process.cwd(), 'data');
    
    if (!fs.existsSync(dataDir)) {
      return NextResponse.json([]);
    }
    
    const repos = fs.readdirSync(dataDir, { withFileTypes: true })
      .filter(dirent => dirent.isDirectory())
      .map(dirent => {
        const repoDir = path.join(dataDir, dirent.name);
        const latestScanFile = path.join(repoDir, 'latest-scan.json');
        const latestDecisionFile = path.join(repoDir, 'latest-decision.json');
        
        let lastScanTime = '';
        let totalSavings = 0;
        
        if (fs.existsSync(latestScanFile)) {
          const scanData = JSON.parse(fs.readFileSync(latestScanFile, 'utf-8'));
          lastScanTime = scanData.timestamp || '';
        }
        
        if (fs.existsSync(latestDecisionFile)) {
          const decisionData = JSON.parse(fs.readFileSync(latestDecisionFile, 'utf-8'));
          totalSavings = decisionData.total_savings_usd || 0;
        }
        
        return {
          repo_full_name: dirent.name.replace('_', '/'),
          last_scan: lastScanTime,
          total_savings: totalSavings,
        };
      })
      .filter(repo => repo.last_scan) // Only return repos with scan data
      .sort((a, b) => new Date(b.last_scan).getTime() - new Date(a.last_scan).getTime());
    
    return NextResponse.json(repos);
  } catch (error) {
    console.error('Error loading repos:', error);
    return NextResponse.json(
      { error: 'Failed to load repos' },
      { status: 500 }
    );
  }
}

