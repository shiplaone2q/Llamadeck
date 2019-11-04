import { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Checkbox } from "@/components/ui/checkbox";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Loader2 } from "lucide-react";
import { apiKeysApi } from "@/lib/api";
import { PermissionMode, type CreateKeyRequest, type InstancePermissionSpec } from "@/types/apiKey";
import { useInstances } from "@/contexts/InstancesContext";
import { format } from "date-fns";

interface CreateApiKeyDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onKeyCreated: (plainTextKey: string) => void;
}

interface InstancePermissionEntry {
  selected: boolean
  canStart: boolean
  canEvict: boolean
}

// There are three meaningful permission levels. Modeling them as tiers makes the
// impossible "evict without start" combo unrepresentable, so the state never holds it.
type InstanceTier = "running" | "start" | "evict";

const tierFromEntry = (e?: InstancePermissionEntry): InstanceTier => {
  if (!e?.canStart) return "running";
  return e.canEvict ? "evict" : "start";
};

function CreateApiKeyDialog({ open, onOpenChange, onKeyCreated }: CreateApiKeyDialogProps) {
  const { instances } = useInstances();
  const [name, setName] = useState("");
  const [permissionMode, setPermissionMode] = useState<PermissionMode>(PermissionMode.AllowAll);
  const [expiresAt, setExpiresAt] = useState<string>("");
  const [instancePermissions, setInstancePermissions] = useState<Record<number, InstancePermissionEntry>>({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const formatDisplayDate = (dateString: string) => {
    if (!dateString) return null;
    try {
      const date = new Date(dateString);
      return format(date, "d MMMM yyyy");
    } catch {
      return null;
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    // Validation
    if (!name.trim()) {
