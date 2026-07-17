// ui/src/components/InstanceCard.tsx
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import type { Instance } from "@/types/instance";
import { Edit, ExternalLink, FileText, Play, Square, Trash2, MoreHorizontal, Download, Boxes, Layers } from "lucide-react";
import LogsDialog from "@/components/LogDialog";
import ModelsDialog from "@/components/ModelsDialog";
import HealthBadge from "@/components/HealthBadge";
import BackendBadge from "@/components/BackendBadge";
import { useState, useEffect } from "react";
import { useInstanceHealth } from "@/hooks/useInstanceHealth";
import { instancesApi, llamaCppApi, type Model } from "@/lib/api";

interface InstanceCardProps {
  instance: Instance;
  startInstance: (name: string) => void;
  stopInstance: (name: string) => void;
  deleteInstance: (name: string) => void;
  editInstance: (instance: Instance) => void;
}

function InstanceCard({
  instance,
  startInstance,
  stopInstance,
  deleteInstance,
  editInstance,
}: InstanceCardProps) {
  const [isLogsOpen, setIsLogsOpen] = useState(false);
  const [isModelsOpen, setIsModelsOpen] = useState(false);
  const [showAllActions, setShowAllActions] = useState(false);
  const [models, setModels] = useState<Model[]>([]);
  const health = useInstanceHealth(instance.name, instance.status);

  const running = instance.status === "running";
  const isLlamaCpp = instance.options?.backend_type === "llama_cpp";

  // Fetch models for llama.cpp instances
  useEffect(() => {
    if (!isLlamaCpp || !running) {
      setModels([]);
      return;
    }

    void (async () => {
      try {
        const fetchedModels = await llamaCppApi.getModels(instance.name);
        setModels(fetchedModels);
      } catch {
        setModels([]);
      }
    })();
  }, [instance.name, isLlamaCpp, running]);

  // Calculate model counts
  const totalModels = models.length;
  const loadedModels = models.filter(m => m.status?.value === "loaded").length;

  const handleStart = () => {
    startInstance(instance.name);
  };

  const handleStop = () => {
    stopInstance(instance.name);
  };

  const handleDelete = () => {
    if (
      confirm(`Are you sure you want to delete instance "${instance.name}"?`)
    ) {
      deleteInstance(instance.name);
    }
  };

  const handleEdit = () => {
    editInstance(instance);
  };
