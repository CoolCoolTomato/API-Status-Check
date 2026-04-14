import { useEffect, useState, useMemo } from 'react';
import type { CheckRecord } from '@/types';
import { apiService } from '@/services/api';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Clock } from 'lucide-react';

interface APIChannel {
  id: string;
  name: string;
  tag: string;
  model: string;
  records: CheckRecord[];
  latest: CheckRecord;
}

function AvailabilityBar({ records }: { records: CheckRecord[] }) {
  const slots = 50;
  const last50 = records.slice(-slots);
  const padded = Array(slots - last50.length).fill(null).concat(last50);

  return (
    <div className="flex gap-px mt-3">
      {padded.map((r, i) =>
        r === null ? (
          <div key={i} className="flex-1 h-7 rounded-sm bg-muted" />
        ) : (
          <div
            key={i}
            title={`${new Date(r.checked_at).toLocaleString()}\n${r.available ? `${r.latency_ms}ms` : r.error_message}`}
            className={`flex-1 h-7 rounded-sm cursor-default transition-opacity hover:opacity-60 ${r.available ? 'bg-green-500' : 'bg-destructive'}`}
          />
        )
      )}
    </div>
  );
}

export default function Dashboard() {
  const [recent, setRecent] = useState<CheckRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [checking, setChecking] = useState(false);

  const token = localStorage.getItem('token');

  const loadData = async () => {
    const res = await apiService.getRecent();
    if (res.code === 0) setRecent(res.data || []);
    setLoading(false);
  };

  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, 30000);
    return () => clearInterval(interval);
  }, []);

  const channels = useMemo<APIChannel[]>(() => {
    const map = new Map<string, CheckRecord[]>();
    for (const r of recent) {
      if (!map.has(r.api_id)) map.set(r.api_id, []);
      map.get(r.api_id)!.push(r);
    }
    return Array.from(map.entries()).map(([, records]) => ({
      id: records[0].api_id,
      name: records[0].name,
      tag: records[0].tag,
      model: records[0].model,
      records,
      latest: records[records.length - 1],
    }));
  }, [recent]);

  const runCheck = async () => {
    setChecking(true);
    await apiService.runCheck();
    setTimeout(() => { loadData(); setChecking(false); }, 2000);
  };

  if (loading) return (
    <div className="flex items-center justify-center h-64 text-muted-foreground">Loading...</div>
  );

  const totalUp = channels.filter(c => c.latest.available).length;

  return (
    <div className="p-8 max-w-7xl mx-auto">
      <div className="flex justify-between items-center mb-2">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">API Status</h1>
          <p className="text-sm text-muted-foreground mt-0.5">
            {totalUp}/{channels.length} channels operational · auto-refresh 30s
          </p>
        </div>
        {token? <Button onClick={runCheck} disabled={checking} size="sm">
          {checking ? 'Checking...' : 'Run Check'}
        </Button>: <></> }
        
      </div>

      <div className="grid gap-4 mt-6 md:grid-cols-2 xl:grid-cols-3">
        {channels.map(ch => {
          const available = ch.latest.available;
          const upCount = ch.records.filter(r => r.available).length;
          const upRate = Math.round((upCount / ch.records.length) * 100);
          const avgLatency = Math.round(
            ch.records.filter(r => r.available).reduce((s, r) => s + r.latency_ms, 0) /
            (ch.records.filter(r => r.available).length || 1)
          );

          return (
            <Card key={ch.id} className="shadow-sm">
              <CardHeader className="pb-2 flex flex-row items-start justify-between space-y-0">
                <div className="min-w-0">
                  <p className="font-semibold text-sm leading-tight truncate">{ch.name}</p>
                  <p className="text-xs text-muted-foreground mt-0.5">{ch.model}</p>
                </div>
                <Badge variant={available ? 'default' : 'destructive'} className="ml-2 shrink-0 text-xs">
                  {available ? 'UP' : 'DOWN'}
                </Badge>
              </CardHeader>
              <CardContent>
                {ch.tag && (
                  <Badge variant="secondary" className="text-xs mb-2">{ch.tag}</Badge>
                )}
                <AvailabilityBar records={ch.records} />
                <div className="flex justify-between text-xs text-muted-foreground mt-2">
                  <span>{upRate}% uptime · {ch.records.length} checks</span>
                  <span>avg {avgLatency}ms</span>
                </div>
                <div className="flex items-center gap-1 text-xs text-muted-foreground mt-1">
                  <Clock size={11} />
                  <span>{new Date(ch.latest.checked_at).toLocaleString()}</span>
                </div>
              </CardContent>
            </Card>
          );
        })}

        {channels.length === 0 && (
          <div className="col-span-full flex flex-col items-center justify-center py-20 text-muted-foreground">
            <p className="text-sm">No check records yet.</p>
            <p className="text-xs mt-1">Add an API in Management and run a check.</p>
          </div>
        )}
      </div>
    </div>
  );
}
