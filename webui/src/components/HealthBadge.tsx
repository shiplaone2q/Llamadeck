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
