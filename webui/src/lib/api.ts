import type { CreateInstanceOptions, Instance } from "@/types/instance";
import type { AppConfig } from "@/types/config";
import type { ApiKey, CreateKeyRequest, CreateKeyResponse, KeyPermissionResponse } from "@/types/apiKey";
import type { DownloadJob, CachedModel, ModelFormat } from "@/types/model";
import { handleApiError } from "./errorUtils";

// Adding baseURI as a prefix to support being served behind a subpath
// e.g. when llmamctl's `/` is served behind a reverse proxy at `/proxy/...`
// the baseURI will be `/proxy/` and the API calls will be made to `/proxy/api/v1/<endpoint>`
export const API_BASE = `${document.baseURI}api/v1`;

// Base API call function with error handling
async function apiCall<T>(
  endpoint: string,
  options: RequestInit = {},
  responseType: "json" | "text" = "json"
): Promise<T> {
  const url = `${API_BASE}${endpoint}`;

  // Get auth token from sessionStorage (same as AuthContext)
  const storedKey = sessionStorage.getItem('llamadeck_management_key');

  // Prepare headers with auth
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };

  // Add auth header if available
  if (storedKey) {
    headers.Authorization = `Bearer ${storedKey}`;
  }

  try {
    const response = await fetch(url, {
      ...options,
      headers,
    });

    // Signal auth loss so the app falls back to the login dialog instead of
    // stranding contexts with error banners
    if (response.status === 401) {
      window.dispatchEvent(new CustomEvent('llamadeck:unauthorized'))
    }

    // Handle errors using centralized error handler
    await handleApiError(response);

    // Handle empty responses (like DELETE)
    if (response.status === 204) {
      return undefined as T;
    }

    // Parse response based on type
    if (responseType === "text") {
      const text = await response.text();
      return text as T;
    } else {
      // Handle empty responses for JSON endpoints
      const contentLength = response.headers.get('content-length');
      if (contentLength === '0') {
        return {} as T; // Return empty object for empty JSON responses
      }
      const data = await response.json() as T;
      return data;
    }
  } catch (error) {
    if (error instanceof Error) {
      throw error;
    }
    throw new Error("Network error occurred");
  }
}

// Server API functions (moved to llama-cpp backend)
export const serverApi = {
  // GET /backends/llama-cpp/help
  getHelp: () => apiCall<string>("/backends/llama-cpp/help", {}, "text"),

  // GET /backends/llama-cpp/version
  getVersion: () => apiCall<string>("/backends/llama-cpp/version", {}, "text"),

  // GET /backends/llama-cpp/devices
  getDevices: () => apiCall<string>("/backends/llama-cpp/devices", {}, "text"),

  // GET /config
  getConfig: () => apiCall<AppConfig>("/config"),
};

// Backend API functions
export const backendsApi = {
  llamaCpp: {
    // POST /backends/llama-cpp/parse-command
    parseCommand: (command: string) =>
      apiCall<CreateInstanceOptions>('/backends/llama-cpp/parse-command', {
        method: 'POST',
        body: JSON.stringify({ command }),
      }),
  },
