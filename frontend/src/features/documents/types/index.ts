export interface Document {
  id: string;
  filename: string;
  display_name?: string;
  uploaded_at: string;
  updated_at?: string;
  element_count: number;
  version?: number;
  category?: string;
  description?: string;
}
