import type React from "react";
import { useEffect, useState } from "react";
import type { CreateInstanceOptions } from "@/types/instance";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import AutoRestartConfiguration from "@/components/instance/AutoRestartConfiguration";
import NumberInput from "@/components/form/NumberInput";
import CheckboxInput from "@/components/form/CheckboxInput";
import SelectInput from "@/components/form/SelectInput";
import { nodesApi, type NodesMap } from "@/lib/api";

interface GeneralTabProps {
  instanceName: string;
  nameError: string;
  isEditing: boolean;
  formData: CreateInstanceOptions;
  onNameChange: (name: string) => void;
  onChange: (key: keyof CreateInstanceOptions, value: unknown) => void;
}

const GeneralTab: React.FC<GeneralTabProps> = ({
  instanceName,
  nameError,
  isEditing,
  formData,
  onNameChange,
  onChange,
}) => {
  const [nodes, setNodes] = useState<NodesMap>({});
  const [loadingNodes, setLoadingNodes] = useState(true);

  useEffect(() => {
    const fetchNodes = async () => {
      try {
        const fetchedNodes = await nodesApi.list();
        setNodes(fetchedNodes);

        const nodeNames = Object.keys(fetchedNodes);
        if (
