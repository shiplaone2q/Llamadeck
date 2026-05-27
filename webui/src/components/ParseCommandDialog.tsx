import type React from "react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { BackendType, type BackendTypeValue, type CreateInstanceOptions } from "@/types/instance";
import { backendsApi } from "@/lib/api";
import { toast } from "sonner";

interface ParseCommandDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onParsed: (options: CreateInstanceOptions) => void;
  backendType: BackendTypeValue;
}

const ParseCommandDialog: React.FC<ParseCommandDialogProps> = ({
  open,
  onOpenChange,
  onParsed,
  backendType,
}) => {
  const [command, setCommand] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleParse = async () => {
    if (!command.trim()) {
      setError("Command cannot be empty");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      let options: CreateInstanceOptions;

      // Parse based on selected backend type
      switch (backendType) {
        case BackendType.LLAMA_CPP:
          options = await backendsApi.llamaCpp.parseCommand(command);
          break;
        case BackendType.MLX_LM:
          options = await backendsApi.mlx.parseCommand(command);
          break;
        case BackendType.VLLM:
          options = await backendsApi.vllm.parseCommand(command);
          break;
        default:
          throw new Error(`Unsupported backend type: ${String(backendType)}`);
      }

      onParsed(options);
      onOpenChange(false);
      setCommand('');
      setError(null);
      toast.success('Command parsed successfully');
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to parse command';
