import { useEffect, useState } from 'react';
import type { APIConfig } from '@/types';
import { apiService } from '@/services/api';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Eye, EyeOff, Pencil, X } from 'lucide-react';

type FormState = { name: string; tag: string; api_url: string; token: string; model: string; enabled: boolean };
const emptyForm: FormState = { name: '', tag: '', api_url: '', token: '', model: '', enabled: true };

const inputCls = "w-full rounded-md border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring";

function TokenField({ value, onChange }: { value: string; onChange: (v: string) => void }) {
  const [show, setShow] = useState(false);
  return (
    <div className="relative">
      <input
        className={inputCls + " pr-9"}
        placeholder="Token *"
        type={show ? 'text' : 'password'}
        value={value}
        onChange={e => onChange(e.target.value)}
        required
      />
      <button type="button" onClick={() => setShow(!show)} className="absolute right-2 top-2 text-muted-foreground hover:text-foreground">
        {show ? <EyeOff size={15} /> : <Eye size={15} />}
      </button>
    </div>
  );
}

function APIForm({ initial, onSubmit, onCancel }: { initial: FormState; onSubmit: (f: FormState) => void; onCancel: () => void }) {
  const [form, setForm] = useState(initial);
  const set = (k: keyof FormState) => (e: React.ChangeEvent<HTMLInputElement>) => setForm({ ...form, [k]: e.target.value });

  return (
    <form onSubmit={e => { e.preventDefault(); onSubmit(form); }} className="space-y-3">
      <div className="grid grid-cols-2 gap-3">
        <input className={inputCls} placeholder="Name *" value={form.name} onChange={set('name')} required />
        <input className={inputCls} placeholder="Tag" value={form.tag} onChange={set('tag')} />
      </div>
      <input className={inputCls} placeholder="API URL *" value={form.api_url} onChange={set('api_url')} required />
      <TokenField value={form.token} onChange={v => setForm({ ...form, token: v })} />
      <input className={inputCls} placeholder="Model *" value={form.model} onChange={set('model')} required />
      <div className="flex justify-end gap-2">
        <Button type="button" size="sm" variant="outline" onClick={onCancel}>Cancel</Button>
        <Button type="submit" size="sm">Save</Button>
      </div>
    </form>
  );
}

export default function Management() {
  const [apis, setApis] = useState<APIConfig[]>([]);
  const [showCreate, setShowCreate] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);

  const loadAPIs = async () => {
    const res = await apiService.getAPIs();
    if (res.code === 0) setApis(res.data || []);
  };

  useEffect(() => { loadAPIs(); }, []);

  const handleCreate = async (form: FormState) => {
    await apiService.createAPI(form);
    setShowCreate(false);
    loadAPIs();
  };

  const handleUpdate = async (id: string, form: FormState) => {
    await apiService.updateAPI(id, form);
    setEditingId(null);
    loadAPIs();
  };

  const handleDelete = async (id: string) => {
    if (confirm('Delete this API and all its history?')) {
      await apiService.deleteAPI(id);
      loadAPIs();
    }
  };

  const toggleEnabled = async (api: APIConfig) => {
    await apiService.updateAPI(api.id, { enabled: !api.enabled });
    loadAPIs();
  };

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">API Management</h1>
          <p className="text-sm text-muted-foreground mt-0.5">{apis.length} configured</p>
        </div>
        <Button size="sm" onClick={() => { setShowCreate(!showCreate); setEditingId(null); }} variant={showCreate ? 'outline' : 'default'}>
          {showCreate ? <><X size={14} className="mr-1" />Cancel</> : '+ Add API'}
        </Button>
      </div>

      {showCreate && (
        <Card className="mb-6 shadow-sm">
          <CardContent className="pt-5">
            <APIForm initial={emptyForm} onSubmit={handleCreate} onCancel={() => setShowCreate(false)} />
          </CardContent>
        </Card>
      )}

      <div className="space-y-3">
        {apis.map(api => (
          <Card key={api.id} className="shadow-sm">
            <CardHeader className="pb-2 flex flex-row items-start justify-between space-y-0">
              <div className="min-w-0">
                <div className="flex items-center gap-2">
                  <p className="font-semibold text-sm">{api.name}</p>
                  {api.tag && <Badge variant="secondary" className="text-xs">{api.tag}</Badge>}
                  <Badge variant={api.enabled ? 'default' : 'outline'} className="text-xs">
                    {api.enabled ? 'Enabled' : 'Disabled'}
                  </Badge>
                </div>
                <p className="text-xs text-muted-foreground mt-1 truncate">{api.api_url}</p>
              </div>
              <div className="flex gap-2 ml-4 shrink-0">
                <Button size="sm" variant="outline" onClick={() => setEditingId(editingId === api.id ? null : api.id)}>
                  <Pencil size={13} />
                </Button>
                <Button size="sm" variant="outline" onClick={() => toggleEnabled(api)}>
                  {api.enabled ? 'Disable' : 'Enable'}
                </Button>
                <Button size="sm" variant="destructive" onClick={() => handleDelete(api.id)}>Delete</Button>
              </div>
            </CardHeader>
            <CardContent>
              {editingId === api.id ? (
                <APIForm
                  initial={{ name: api.name, tag: api.tag, api_url: api.api_url, token: api.token, model: api.model, enabled: api.enabled }}
                  onSubmit={form => handleUpdate(api.id, form)}
                  onCancel={() => setEditingId(null)}
                />
              ) : (
                <p className="text-xs text-muted-foreground">Model: {api.model}</p>
              )}
            </CardContent>
          </Card>
        ))}

        {apis.length === 0 && (
          <div className="flex flex-col items-center justify-center py-20 text-muted-foreground">
            <p className="text-sm">No APIs configured yet.</p>
          </div>
        )}
      </div>
    </div>
  );
}
