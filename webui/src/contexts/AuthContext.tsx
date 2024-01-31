import { type ReactNode, createContext, useCallback, useContext, useEffect, useState } from 'react'
import { authApi, type WhoamiResponse, type WhoamiUser } from '@/lib/api'

type AuthMethod = 'key' | 'session' | null

interface AuthContextState {
  isAuthenticated: boolean
  isLoading: boolean
  apiKey: string | null
  user: WhoamiUser | null
  authMethod: AuthMethod
  oidcEnabled: boolean
  error: string | null
  oidcError: string | null
}

interface AuthContextActions {
  login: (apiKey: string) => Promise<void>
  logout: () => void
  clearError: () => void
  clearOIDCError: () => void
  validateAuth: () => Promise<boolean>
}

type AuthContextType = AuthContextState & AuthContextActions

const AuthContext = createContext<AuthContextType | undefined>(undefined)

interface AuthProviderProps {
  children: ReactNode
}

const AUTH_STORAGE_KEY = 'llamadeck_management_key'

// Error codes the OIDC callback redirects back with (?auth_error=<code>);
// mapped to messages shown in the login dialog.
const AUTH_ERROR_MESSAGES: Record<string, string> = {
  denied: 'Access denied: your account is not a member of an allowed group.',
  state: 'The login request expired or was invalid. Please try again.',
  idp: 'The identity provider rejected the login.',
  exchange: 'Could not complete the login with the identity provider. Please try again.',
  token: 'The identity provider returned an invalid login token.',
  groups: 'The identity provider returned an unusable groups claim.',
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [isLoading, setIsLoading] = useState(true)
  const [apiKey, setApiKey] = useState<string | null>(null)
  const [user, setUser] = useState<WhoamiUser | null>(null)
  const [authMethod, setAuthMethod] = useState<AuthMethod>(null)
  const [oidcEnabled, setOidcEnabled] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [oidcError, setOidcError] = useState<string | null>(null)

  // A failed OIDC login redirects back to the app root with ?auth_error=
  // <code>. Map it to a message for the login dialog and strip it from the
  // URL so a refresh or SSO retry doesn't replay it.
  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const code = params.get('auth_error')
    if (!code) return

    params.delete('auth_error')
    const rest = params.toString()
    window.history.replaceState(null, '', window.location.pathname + (rest ? `?${rest}` : ''))
    setOidcError(AUTH_ERROR_MESSAGES[code] ?? 'Login failed. Please try again.')
  }, [])
