import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen, waitFor, act } from "@testing-library/react";
import type { ReactNode } from "react";
import { InstancesProvider, useInstances } from "@/contexts/InstancesContext";
import { instancesApi } from "@/lib/api";
import { BackendType, type Instance } from "@/types/instance";
import { AuthProvider } from "../AuthContext";

// Mock the API module
vi.mock("@/lib/api", () => ({
  instancesApi: {
    list: vi.fn(),
    get: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    start: vi.fn(),
    stop: vi.fn(),
    restart: vi.fn(),
    delete: vi.fn(),
    getHealth: vi.fn(),
  },
}));

// Mock health service
vi.mock("@/lib/healthService", () => ({
  healthService: {
    subscribe: vi.fn(() => () => {}),
    refreshHealth: vi.fn(() => Promise.resolve()),
    checkHealthAfterOperation: vi.fn(),
    performHealthCheck: vi.fn(() => Promise.resolve({
      state: 'ready',
      instanceStatus: 'running',
      lastChecked: new Date(),
      source: 'http'
    })),
  },
  checkHealth: vi.fn(() => Promise.resolve({
    state: 'ready',
    instanceStatus: 'running',
    lastChecked: new Date(),
    source: 'http'
  })),
}));

// Test component to access context
function TestComponent() {
  const {
    instances,
    loading,
    error,
    createInstance,
    updateInstance,
    startInstance,
    stopInstance,
    restartInstance,
    deleteInstance,
    clearError,
  } = useInstances();

  return (
    <div>
      <div data-testid="loading">{loading.toString()}</div>
      <div data-testid="error">{error || "no-error"}</div>
      <div data-testid="instances-count">{instances.length}</div>
      {instances.map((instance) => (
        <div key={instance.name} data-testid={`instance-${instance.name}`}>
          {instance.name}:{instance.status}
        </div>
      ))}

      {/* Action buttons for testing with specific instances */}
      <button
        type="button"
        onClick={() => void createInstance("new-instance", { backend_type: BackendType.LLAMA_CPP, backend_options: { model: "test.gguf" } })}
        data-testid="create-instance"
      >
        Create Instance
      </button>
      <button
        type="button"
        onClick={() => void updateInstance("instance1", { backend_type: BackendType.LLAMA_CPP, backend_options: { model: "updated.gguf" } })}
        data-testid="update-instance"
      >
        Update Instance
      </button>
      <button
        type="button"
        onClick={() => void startInstance("instance2")}
        data-testid="start-instance"
      >
        Start Instance2
      </button>
      <button
        type="button"
        onClick={() => void stopInstance("instance1")}
        data-testid="stop-instance"
      >
        Stop Instance1
      </button>
      <button
        type="button"
        onClick={() => void restartInstance("instance1")}
        data-testid="restart-instance"
      >
        Restart Instance1
      </button>
      <button
        type="button"
        onClick={() => void deleteInstance("instance2")}
        data-testid="delete-instance"
      >
        Delete Instance2
      </button>
      <button type="button" onClick={clearError} data-testid="clear-error">
        Clear Error
      </button>
    </div>
  );
