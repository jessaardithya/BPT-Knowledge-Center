"use client";

import { useEffect, useState, useRef } from "react";
import {
  fetchDocuments,
  deleteDocument,
  updateDocument,
} from "../api/documentService";
import { Document } from "../types";
import {
  Folder,
  X,
  Save,
  Edit2,
  FileText,
  MoreHorizontal,
  Calendar,
  Database,
} from "lucide-react";
import { apiClient } from "@/utils/api";

export default function DocumentList() {
  const [documents, setDocuments] = useState<Document[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeMenu, setActiveMenu] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [reuploadingId, setReuploadingId] = useState<string | null>(null);
  const [editForm, setEditForm] = useState<{
    id: string;
    display_name: string;
    category: string;
    description: string;
  } | null>(null);

  useEffect(() => {
    fetchDocs();
  }, []);

  const fetchDocs = async () => {
    try {
      const data = await fetchDocuments();
      setDocuments(data);
    } catch {
      console.error("Fetch failed");
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Are you sure?")) return;
    try {
      await deleteDocument(id);
      setDocuments(documents.filter((doc) => doc.id !== id));
    } catch {
      alert("Delete failed");
    }
  };

  const handleReupload = (docId: string) => {
    setReuploadingId(docId);
    fileInputRef.current?.click();
    setActiveMenu(null);
  };

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !reuploadingId) return;

    const doc = documents.find((d) => d.id === reuploadingId);
    if (!doc) return;

    const formData = new FormData();
    formData.append("file", file);
    formData.append("document_id", reuploadingId);
    formData.append("display_name", doc.display_name || doc.filename);

    try {
      await apiClient.post("/documents/upload", formData);
      await fetchDocs();
    } catch {
      alert("Re-upload failed");
    } finally {
      setReuploadingId(null);
      if (fileInputRef.current) fileInputRef.current.value = "";
    }
  };

  const handleEditClick = (doc: Document) => {
    setEditForm({
      id: doc.id,
      display_name: doc.display_name || doc.filename,
      category: doc.category || "",
      description: doc.description || "",
    });
    setActiveMenu(null);
  };

  const handleSaveEdit = async () => {
    if (!editForm) return;
    try {
      await updateDocument(editForm.id, {
        display_name: editForm.display_name,
        category: editForm.category,
        description: editForm.description,
      });
      await fetchDocs();
      setEditForm(null);
    } catch {
      alert("Update failed");
    }
  };

  if (loading) return null;

  return (
    <div className="space-y-6 max-w-6xl mx-auto">
      <input
        type="file"
        ref={fileInputRef}
        className="hidden"
        accept=".pdf,.pptx"
        onChange={handleFileSelect}
      />

      {documents.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 bg-white border border-gray-200 border-dashed rounded-xl">
          <div className="p-4 bg-gray-50 rounded-full mb-4">
            <Folder className="text-gray-400" size={32} />
          </div>
          <h3 className="text-lg font-medium text-gray-900">
            No documents yet
          </h3>
          <p className="text-gray-500 mt-1">Upload files to get started</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
          {documents.map((doc) => (
            <div
              key={doc.id}
              className="group relative bg-white border border-gray-200 rounded-xl p-5 hover:border-gray-300 hover:shadow-sm transition-all duration-200"
            >
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-lg bg-gray-50 border border-gray-100 flex items-center justify-center text-gray-500">
                    <FileText size={20} className="text-gray-600" />
                  </div>
                  <div className="flex flex-col">
                    <h3 className="text-sm font-semibold text-gray-900 leading-tight mb-0.5 line-clamp-1">
                      {doc.display_name || doc.filename}
                    </h3>
                    <div className="flex items-center gap-2">
                      <span className="text-[10px] uppercase font-bold tracking-wider text-gray-400">
                        {doc.category || "Uncategorized"}
                      </span>
                      {doc.version && (
                        <span className="px-1.5 py-0.5 bg-gray-100 text-gray-600 text-[10px] font-medium rounded-md">
                          v{doc.version}
                        </span>
                      )}
                    </div>
                  </div>
                </div>

                <div className="relative">
                  <button
                    onClick={() =>
                      setActiveMenu(activeMenu === doc.id ? null : doc.id)
                    }
                    className="p-1.5 text-gray-300 hover:text-gray-600 hover:bg-gray-50 rounded-md transition-colors"
                  >
                    <MoreHorizontal size={18} />
                  </button>

                  {activeMenu === doc.id && (
                    <div className="absolute right-0 mt-2 w-32 bg-white border border-gray-200 rounded-lg shadow-lg z-10 py-1">
                      <button
                        onClick={() => handleEditClick(doc)}
                        className="w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-2"
                      >
                        <Edit2 size={14} className="text-gray-400" />
                        Edit Details
                      </button>
                      <button
                        onClick={() => handleReupload(doc.id)}
                        className="w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-50"
                      >
                        Update File
                      </button>
                      <button
                        onClick={() => handleDelete(doc.id)}
                        className="w-full text-left px-3 py-2 text-sm text-red-600 hover:bg-red-50"
                      >
                        Delete
                      </button>
                    </div>
                  )}
                </div>
              </div>

              <div className="text-xs text-gray-500 border-t border-gray-100 pt-3 flex items-center justify-between">
                <div className="flex items-center gap-1.5">
                  <Calendar size={12} />
                  {new Date(doc.uploaded_at).toLocaleDateString()}
                </div>
                <div className="flex items-center gap-1.5">
                  <Database size={12} />
                  {doc.element_count} chunks
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Edit Modal */}
      {editForm && (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl w-full max-w-lg shadow-2xl animate-in zoom-in-95 duration-200">
            <div className="flex items-center justify-between p-6 border-b border-gray-100">
              <h3 className="text-lg font-bold text-gray-900">Edit Document</h3>
              <button
                onClick={() => setEditForm(null)}
                className="p-2 hover:bg-gray-100 rounded-full text-gray-500 transition-colors"
              >
                <X size={20} />
              </button>
            </div>

            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Display Name
                </label>
                <input
                  type="text"
                  value={editForm.display_name}
                  onChange={(e) =>
                    editForm &&
                    setEditForm({ ...editForm, display_name: e.target.value })
                  }
                  className="w-full px-4 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none"
                  placeholder="e.g. Q1 Financial Report"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Category
                </label>
                <input
                  type="text"
                  value={editForm.category}
                  onChange={(e) =>
                    editForm &&
                    setEditForm({ ...editForm, category: e.target.value })
                  }
                  className="w-full px-4 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none"
                  placeholder="e.g. Finance, Technical, Sales"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Description
                </label>
                <textarea
                  value={editForm.description}
                  onChange={(e) =>
                    editForm &&
                    setEditForm({ ...editForm, description: e.target.value })
                  }
                  className="w-full px-4 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none min-h-[100px] resize-none"
                  placeholder="Add a brief description..."
                />
              </div>
            </div>

            <div className="p-6 border-t border-gray-100 flex justify-end gap-3 bg-gray-50/50 rounded-b-2xl">
              <button
                onClick={() => setEditForm(null)}
                className="px-4 py-2 text-gray-600 font-medium hover:bg-gray-200/50 rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleSaveEdit}
                className="flex items-center gap-2 px-4 py-2 bg-[#0A2540] text-white font-medium rounded-lg hover:bg-[#0A2540]/90 transition-colors shadow-lg shadow-blue-900/10"
              >
                <Save size={16} />
                Save Changes
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
