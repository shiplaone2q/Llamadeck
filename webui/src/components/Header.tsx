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
