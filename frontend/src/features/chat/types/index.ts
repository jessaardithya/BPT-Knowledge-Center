export interface Source {
  filename: string;
  page: number;
}

export interface Message {
  id: string;
  role: "user" | "bot";
  content: string;
  sources?: Source[];
}
