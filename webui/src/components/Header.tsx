import { Button } from "@/components/ui/button";
import { HelpCircle, LogOut, Moon, Settings, Sun, UserCircle } from "lucide-react";
import { useAuth } from "@/contexts/AuthContext";
import { useTheme } from "@/contexts/ThemeContext";

interface HeaderProps {
  onCreateInstance: () => void;
  onShowSystemInfo: () => void;
  onShowSettings: () => void;
}

function Header({ onCreateInstance, onShowSystemInfo, onShowSettings }: HeaderProps) {
  const { logout, user } = useAuth();
  const { theme, toggleTheme } = useTheme();

  const handleLogout = () => {
    logout();
  };

  return (
    <header className="bg-card border-b border-border">
      <div className="container mx-auto max-w-4xl px-4 py-4">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold text-foreground">
            Llamadeck Dashboard
          </h1>

          <div className="flex items-center gap-2">
            <Button onClick={onCreateInstance} data-testid="create-instance-button">
              Create Instance
            </Button>

            <Button
              variant="outline"
              size="icon"
              onClick={toggleTheme}
              data-testid="theme-toggle-button"
              title={`Switch to ${theme === 'light' ? 'dark' : 'light'} mode`}
            >
              {theme === 'light' ? <Moon className="h-4 w-4" /> : <Sun className="h-4 w-4" />}
            </Button>

            <Button
              variant="outline"
