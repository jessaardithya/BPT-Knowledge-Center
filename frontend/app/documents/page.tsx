"use client";
import DocumentList from "@/features/documents/components/DocumentList";
import { Plus } from "lucide-react";

export default function DocumentsPage() {
  return (
    <main className="flex-1 h-screen overflow-auto bg-white font-sans">
      <div className="max-w-7xl mx-auto px-8 py-10">
        <div className="flex items-end justify-between mb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 tracking-tight">
              Documents
            </h1>
            <p className="text-gray-500 mt-2">
              Manage your knowledge base and training data.
            </p>
          </div>
          <button className="flex items-center gap-2 px-4 py-2.5 bg-black text-white text-sm font-medium rounded-lg hover:bg-gray-800 transition-colors shadow-sm">
            <Plus size={16} />
            <span>Upload New</span>
          </button>
        </div>
        <DocumentList />
      </div>
    </main>
  );
}
