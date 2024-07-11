// ui/src/components/HealthBadge.tsx
import type React from "react";
import { Badge } from "@/components/ui/badge";
import type { HealthStatus } from "@/types/instance";
import { CheckCircle, Loader2, XCircle, Clock } from "lucide-react";

interface HealthBadgeProps {
  health?: HealthStatus;
}

const HealthBadge: React.FC<HealthBadgeProps> = ({ health }) => {
  if (!health) {
    return null;
  }

  const getIcon = () => {
    switch (health.state) {
      case "ready":
        return <CheckCircle className="h-3 w-3" />;
      case "starting":
        return <Loader2 className="h-3 w-3 animate-spin" />;
      case "restarting":
        return <Loader2 className="h-3 w-3 animate-spin" />;
      case "shutting_down":
        return <Loader2 className="h-3 w-3 animate-spin" />;
      case "stopped":
        return <Clock className="h-3 w-3" />;
      case "failed":
        return <XCircle className="h-3 w-3" />;
    }
  };

  const getVariant = () => {
    switch (health.state) {
      case "ready":
        return "default";
      case "starting":
        return "outline";
      case "restarting":
        return "outline";
      case "shutting_down":
        return "outline";
