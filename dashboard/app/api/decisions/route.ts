import { NextResponse } from 'next/server';
import fs from 'fs';
import path from 'path';

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url);
    const repo = searchParams.get('repo') || '';
    
    const decisions: any[] = [];
    
    // If repo is specified, get from data directory
    if (repo) {
      const isVercel = process.env.VERCEL === '1';
      const baseDir = isVercel ? '/tmp/costguard-data' : path.join(process.cwd(), 'data');
      const repoDir = path.join(baseDir, repo.replace('/', '_'));
      
      if (fs.existsSync(repoDir)) {
        const files = fs.readdirSync(repoDir)
          .filter(f => f.startsWith('decision-') && f.endsWith('.json'))
          .map(f => {
            const filePath = path.join(repoDir, f);
            const content = JSON.parse(fs.readFileSync(filePath, 'utf-8'));
            return {
              scan_id: f.replace('decision-', '').replace('.json', ''),
              timestamp: content.timestamp || fs.statSync(filePath).mtime.toISOString(),
              total_savings: content.total_savings_usd || 0,
              actions_applied: content.actions_to_apply || 0,
              pr_url: content.pr_url || '',
              pr_number: content.pr_number || 0,
              summary: content.summary || '',
              repo_full_name: content.repo_full_name || repo,
            };
          })
          .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
          .slice(0, 10);
        
        decisions.push(...files);
      }
    } else {
      // Get decisions from all repos
      const isVercel = process.env.VERCEL === '1';
      const dataDir = isVercel ? '/tmp/costguard-data' : path.join(process.cwd(), 'data');
      
      if (fs.existsSync(dataDir)) {
        const repos = fs.readdirSync(dataDir, { withFileTypes: true })
          .filter(dirent => dirent.isDirectory());
        
        for (const repoDir of repos) {
          const repoPath = path.join(dataDir, repoDir.name);
          const files = fs.readdirSync(repoPath)
            .filter(f => f.startsWith('decision-') && f.endsWith('.json'))
            .map(f => {
              const filePath = path.join(repoPath, f);
              const content = JSON.parse(fs.readFileSync(filePath, 'utf-8'));
              return {
                scan_id: f.replace('decision-', '').replace('.json', ''),
                timestamp: content.timestamp || fs.statSync(filePath).mtime.toISOString(),
                total_savings: content.total_savings_usd || 0,
                actions_applied: content.actions_to_apply || 0,
                pr_url: content.pr_url || '',
                pr_number: content.pr_number || 0,
                summary: content.summary || '',
                repo_full_name: content.repo_full_name || repoDir.name.replace('_', '/'),
              };
            });
          
          decisions.push(...files);
        }
      }
      
      // Fallback: try local .costguard directory (for local dev)
      const decisionsDir = path.join(process.cwd(), '../../.costguard');
      if (fs.existsSync(decisionsDir)) {
        const files = fs.readdirSync(decisionsDir)
          .filter(f => f.startsWith('decisions-') && f.endsWith('.json'))
          .map(f => {
            const filePath = path.join(decisionsDir, f);
            const content = JSON.parse(fs.readFileSync(filePath, 'utf-8'));
            return {
              scan_id: f.replace('decisions-', '').replace('.json', ''),
              timestamp: fs.statSync(filePath).mtime.toISOString(),
              total_savings: content.total_savings_usd || 0,
              actions_applied: content.actions_to_apply || 0,
              pr_url: '',
              pr_number: 0,
              summary: content.summary || '',
              repo_full_name: 'local',
            };
          });
        
        decisions.push(...files);
      }
    }
    
    // Sort by timestamp and limit
    decisions.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime());
    
    return NextResponse.json(decisions.slice(0, 50));
  } catch (error) {
    console.error('Error loading decisions:', error);
    return NextResponse.json(
      { error: 'Failed to load decisions' },
      { status: 500 }
    );
  }
}

