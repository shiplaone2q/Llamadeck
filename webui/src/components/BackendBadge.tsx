import type React from "react";
import { Badge } from "@/components/ui/badge";
import { BackendType, type BackendTypeValue } from "@/types/instance";
import { Server, Package } from "lucide-react";

interface BackendBadgeProps {
  backend?: BackendTypeValue;
  docker?: boolean;
  customName?: string;
}

const BackendBadge: React.FC<BackendBadgeProps> = ({ backend, docker, customName }) => {
  if (!backend) {
    return null;
  }

  const getText = () => {
    switch (backend) {
      case BackendType.LLAMA_CPP:
        return "llama.cpp";
      case BackendType.MLX_LM:
        return "MLX";
      case BackendType.VLLM:
        return "vLLM";
      case BackendType.CUSTOM:
        return customName || "custom";
      default:
        return backend;
