export enum PermissionMode {
  AllowAll = "allow_all",
  PerInstance = "per_instance"
}

export interface ApiKey {
  id: number
  name: string
  user_id: string
  permission_mode: PermissionMode
  expires_at: number | null
  created_at: number
  updated_at: number
