import { BrowserRouter, Routes, Route, NavLink, useNavigate } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import Management from './pages/Management';
import Login from './pages/Login';
import ProtectedRoute from './components/ProtectedRoute';
import { useTheme } from '@/components/theme-provider';
import { Sun, Moon } from 'lucide-react';

function ThemeToggle() {
  const { theme, setTheme } = useTheme();
  return (
    <button
      onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
      className="text-muted-foreground hover:text-foreground transition-colors p-1.5 rounded-md hover:bg-accent/50"
      title="Toggle theme"
    >
      {theme === 'dark' ? <Sun size={15} /> : <Moon size={15} />}
    </button>
  );
}

function Header() {
  const navigate = useNavigate();
  const token = localStorage.getItem('token');
  const navCls = ({ isActive }: { isActive: boolean }) =>
    `text-sm font-medium transition-colors px-3 py-1.5 rounded-md ${isActive ? 'bg-accent text-accent-foreground' : 'text-muted-foreground hover:text-foreground hover:bg-accent/50'}`;

  const logout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  return (
    <header className="border-b sticky top-0 z-10 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="max-w-7xl mx-auto flex items-center gap-1 px-8 h-14">
        <span className="font-bold text-sm mr-4">API status check</span>
        <NavLink to="/" end className={navCls}>Dashboard</NavLink>
        {token && <NavLink to="/management" className={navCls}>Management</NavLink>}
        <div className="ml-auto flex items-center gap-1">
          <ThemeToggle />
          {token
            ? <button onClick={logout} className="text-xs text-muted-foreground hover:text-foreground transition-colors">Logout</button>
            : <></>}
        </div>
      </div>
    </header>
  );
}

export function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-background">
        <Header />
        <main>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/login" element={<Login />} />
            <Route path="/management" element={
              <ProtectedRoute><Management /></ProtectedRoute>
            } />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}

export default App;
