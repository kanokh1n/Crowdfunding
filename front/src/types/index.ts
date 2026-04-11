// User
export interface User {
  id: number
  email: string
  fio: string
  description?: string
  phone?: string
  profile_img?: string
  role: 'user' | 'admin'
  is_verified: boolean
  created_at: string
  updated_at: string
}

// Category
export interface Category {
  id: number
  title: string
}

// Project
export interface Project {
  id: number
  user_id: number
  title: string
  description: string
  goal_amount: number
  current_amount: number
  project_img?: string
  status: 'pending_ai' | 'pending_human' | 'active' | 'rejected_ai' | 'rejected' | 'completed' | 'cancelled'
  end_date?: string
  created_at: string
  updated_at: string
  // Relations (included in API responses)
  user?: User
  categories?: Category[]
  likes_count?: number
  is_liked?: boolean
}

// Project moderation
export interface ProjectModeration {
  id: number
  project_id: number
  ai_status: 'pending' | 'passed' | 'failed'
  ai_score?: number
  ai_flags?: string[]
  ai_checked_at?: string
  human_status: 'pending' | 'approved' | 'rejected'
  moderator_id?: number
  moderator_note?: string
  human_moderated_at?: string
  created_at: string
  updated_at: string
}

// Pledge
export interface Pledge {
  id: number
  user_id: number
  project_id: number
  amount: number
  created_at: string
  user?: User
}

// Comment
export interface Comment {
  id: number
  user_id: number
  project_id: number
  content: string
  created_at: string
  user?: User
}

// Like
export interface Like {
  id: number
  user_id: number
  project_id: number
  created_at: string
}

// Message
export interface Message {
  id: number
  sender_id: number
  recipient_id: number
  project_id: number
  title: string
  content: string
  is_read: boolean
  created_at: string
  sender?: User
  recipient?: User
  project?: Project
}

// Auth tokens
export interface AuthTokens {
  access_token: string
  refresh_token: string
  user: User
}

// API pagination
export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  limit: number
}

// Project creation/update
export interface CreateProjectInput {
  title: string
  description?: string
  goal_amount: number
  end_date?: string
  project_img?: string
  category_ids?: number[]
}

export interface UpdateProjectInput extends Partial<CreateProjectInput> {}

// Auth
export interface RegisterInput {
  email: string
  password: string
  fio: string
}

export interface LoginInput {
  email: string
  password: string
}

// Pagination params
export interface PaginationParams {
  page?: number
  limit?: number
}

// Project list params
export interface ProjectListParams extends PaginationParams {
  category_id?: number
  search?: string
  status?: string
  sort?: 'current_amount' | 'likes_count'
}
