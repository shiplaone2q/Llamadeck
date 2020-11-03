import { useState } from "react";
import Header from "@/components/Header";
import InstanceList from "@/components/InstanceList";
import ModelsList from "@/components/ModelsList";
import NodesList from "@/components/nodes/NodesList";
import InstanceDialog from "@/components/InstanceDialog";
import LoginDialog from "@/components/LoginDialog";
import SystemInfoDialog from "./components/SystemInfoDialog";
import SettingsDialog from "./components/settings/SettingsDialog";
import type { CreateInstanceOptions, Instance } from "@/types/instance";
import { useInstances } from "@/contexts/InstancesContext";
import { useAuth } from "@/contexts/AuthContext";
import { ModelsProvider } from "@/contexts/ModelsContext";
import { ThemeProvider } from "@/contexts/ThemeContext";
import { Toaster } from "sonner";
import { cn } from "@/lib/utils";

function App() {
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const [activeTab, setActiveTab] = useState<'instances' | 'nodes' | 'models'>('instances');
  const [isInstanceModalOpen, setIsInstanceModalOpen] = useState(false);
  const [isSystemInfoModalOpen, setIsSystemInfoModalOpen] = useState(false);
  const [isSettingsModalOpen, setIsSettingsModalOpen] = useState(false);
  const [editingInstance, setEditingInstance] = useState<Instance | undefined>(
    undefined
  );
  const { createInstance, updateInstance } = useInstances();

  const handleCreateInstance = () => {
    setEditingInstance(undefined);
    setIsInstanceModalOpen(true);
  };

  const handleEditInstance = (instance: Instance) => {
    setEditingInstance(instance);
    setIsInstanceModalOpen(true);
  };

  const handleSaveInstance = (name: string, options: CreateInstanceOptions) => {
    if (editingInstance) {
      void updateInstance(editingInstance.name, options);
    } else {
      void createInstance(name, options);
    }
  };

  const handleShowSystemInfo = () => {
    setIsSystemInfoModalOpen(true);
  };

