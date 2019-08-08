import { useEffect, useState, useCallback, Fragment } from "react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Trash2, Copy, Check, X, ChevronDown, ChevronRight } from "lucide-react";
import { apiKeysApi } from "@/lib/api";
import { type ApiKey, type KeyPermissionResponse, PermissionMode } from "@/types/apiKey";
import CreateApiKeyDialog from "@/components/apikeys/CreateApiKeyDialog";
import { format, formatDistanceToNow } from "date-fns";

function ApiKeysSection() {
  const [keys, setKeys] = useState<ApiKey[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [expandedRowId, setExpandedRowId] = useState<number | null>(null);
  const [newKeyPlainText, setNewKeyPlainText] = useState<string | null>(null);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [copiedKey, setCopiedKey] = useState(false);
  const [permissions, setPermissions] = useState<Record<number, KeyPermissionResponse[]>>({});
  const [loadingPermissions, setLoadingPermissions] = useState<Record<number, boolean>>({});

  const fetchKeys = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await apiKeysApi.list();
      setKeys(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load API keys");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void fetchKeys();
  }, [fetchKeys]);

  const fetchPermissions = async (keyId: number) => {
    if (permissions[keyId]) return;

    setLoadingPermissions({ ...loadingPermissions, [keyId]: true });
    try {
      const data = await apiKeysApi.getPermissions(keyId);
      setPermissions({ ...permissions, [keyId]: data });
    } catch (err) {
      console.error("Failed to load permissions:", err);
    } finally {
      setLoadingPermissions({ ...loadingPermissions, [keyId]: false });
    }
  };

  const handleKeyCreated = (plainTextKey: string) => {
    setNewKeyPlainText(plainTextKey);
    void fetchKeys();
    setCreateDialogOpen(false);
  };

  const dismissSuccessBanner = () => {
    setNewKeyPlainText(null);
  };

