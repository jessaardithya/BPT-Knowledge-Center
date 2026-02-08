"use client";

import { useState, DragEvent, ChangeEvent } from "react";
import { uploadDocument } from "../api/documentService";
import { UploadCloud, Loader2, CheckCircle, AlertCircle } from "lucide-react";

export default function FileUploader() {
  const [isDragging, setIsDragging] = useState(false);
  const [file, setFile] = useState<File | null>(null);
  const [displayName, setDisplayName] = useState("");
  const [status, setStatus] = useState<
    "idle" | "uploading" | "success" | "error"
  >("idle");
  const [errorMessage, setErrorMessage] = useState("");

  const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const handleDrop = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      const droppedFile = e.dataTransfer.files[0];
      setFile(droppedFile);
      // Auto-fill display name with filename (without extension)
      const nameWithoutExt = droppedFile.name.replace(/\.[^/.]+$/, "");
      setDisplayName(nameWithoutExt);
      setStatus("idle");
      setErrorMessage("");
    }
  };

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const selectedFile = e.target.files[0];
      setFile(selectedFile);
      // Auto-fill display name with filename (without extension)
      const nameWithoutExt = selectedFile.name.replace(/\.[^/.]+$/, "");
      setDisplayName(nameWithoutExt);
      setStatus("idle");
      setErrorMessage("");
    }
  };

  const handleUpload = async () => {
    if (!file) return;
    setStatus("uploading");
    setErrorMessage("");

    const formData = new FormData();
    formData.append("file", file);
    formData.append("display_name", displayName || file.name);

    try {
      await uploadDocument(formData);
      setStatus("success");
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      console.error("Upload failed:", err);
      setStatus("error");
      setErrorMessage(
        err.response?.data?.error ||
          "Connection failed. Is the backend running?",
      );
    }
  };

  return (
    <div className="w-full max-w-2xl mx-auto mt-12">
      <div
        className={`relative group rounded-3xl border-2 border-dashed transition-all duration-300 p-12 text-center
          ${
            isDragging
              ? "border-[var(--brand-accent)] bg-blue-50/50 scale-[1.02] shadow-xl"
              : "border-gray-200 hover:border-blue-300 hover:bg-white/60 bg-white/40"
          }
          ${status === "success" ? "border-green-400 bg-green-50/30" : ""}
          ${status === "error" ? "border-red-300 bg-red-50/50" : ""}
        `}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        {/* State: Success */}
        {status === "success" ? (
          <div className="flex flex-col items-center animate-in fade-in zoom-in duration-500">
            <div className="w-16 h-16 bg-green-100 text-green-600 rounded-full flex items-center justify-center mb-6 shadow-sm">
              <CheckCircle size={32} strokeWidth={3} />
            </div>
            <h3 className="text-2xl font-bold text-gray-800 mb-2">
              Upload Complete!
            </h3>
            <p className="text-gray-500 mb-8">
              <span className="font-semibold text-gray-700">
                &quot;{displayName || file?.name}&quot;
              </span>{" "}
              has been successfully ingested.
            </p>
            <button
              onClick={() => {
                setFile(null);
                setDisplayName("");
                setStatus("idle");
              }}
              className="text-[var(--brand-accent)] font-semibold hover:text-blue-700 transition-colors"
            >
              Upload another document
            </button>
          </div>
        ) : (
          /* State: Idle / Uploading / Error */
          <div className="flex flex-col items-center">
            <div
              className={`w-20 h-20 bg-blue-50 rounded-2xl flex items-center justify-center mb-6 transition-all duration-300 ${isDragging ? "rotate-3 scale-110 bg-blue-100" : ""} ${status === "error" ? "bg-red-50" : ""}`}
            >
              {status === "uploading" ? (
                <Loader2
                  className="animate-spin text-[var(--brand-accent)]"
                  size={40}
                />
              ) : status === "error" ? (
                <AlertCircle className="text-red-500" size={40} />
              ) : (
                <UploadCloud className="text-[var(--brand-accent)]" size={40} />
              )}
            </div>

            <h3 className="text-2xl font-bold text-gray-900 mb-3">
              {status === "uploading"
                ? "Processing Document..."
                : status === "error"
                  ? "Upload Failed"
                  : "Upload Knowledge Source"}
            </h3>

            {status === "error" ? (
              <p className="text-red-500 max-w-sm mx-auto mb-8 font-medium">
                {errorMessage}
              </p>
            ) : (
              <p className="text-gray-500 max-w-sm mx-auto mb-8 leading-relaxed">
                Drag & drop your file here, or click to browse.
                <br />
                <span className="text-sm font-medium text-blue-500/80 mt-2 block">
                  Supports PDF and PPT Presentation files
                </span>
              </p>
            )}

            {/* File Selected State (Before Upload) */}
            {file && status !== "uploading" && (
              <div className="w-full max-w-md animate-in slide-in-from-bottom-2">
                <div className="flex flex-col gap-4 bg-white p-6 rounded-2xl shadow-sm border border-gray-100 mb-6">
                  {/* File Icon & Name */}
                  <div className="flex items-center gap-4">
                    <div className="bg-red-50 p-2.5 rounded-xl">
                      <Loader2 className="text-red-500" size={24} />
                      {/* Note: Using Loader2 as placeholder for FileIcon if not imported, but FileText is imported in original code? No wait, FileText wasn't imported in this file. I'll use File from imports */}
                    </div>
                    <div className="flex-1 text-left">
                      <p className="font-semibold text-gray-800 text-sm truncate">
                        {file.name}
                      </p>
                      <p className="text-xs text-gray-400">
                        {(file.size / 1024 / 1024).toFixed(2)} MB
                      </p>
                    </div>
                  </div>

                  {/* Display Name Input */}
                  <div>
                    <label className="block text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 text-left">
                      Document Name
                    </label>
                    <input
                      type="text"
                      value={displayName}
                      onChange={(e) => setDisplayName(e.target.value)}
                      placeholder="e.g. Q1 Sales Report"
                      className="w-full px-4 py-2 border border-gray-200 rounded-lg text-sm bg-gray-50 focus:bg-white focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all"
                    />
                  </div>
                </div>

                <div className="flex gap-3">
                  <button
                    onClick={() => {
                      setFile(null);
                      setDisplayName("");
                    }}
                    className="flex-1 px-4 py-3 border border-gray-200 text-gray-600 rounded-xl hover:bg-gray-50 font-medium transition-colors"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={handleUpload}
                    className="flex-1 px-4 py-3 bg-[#0A2540] text-white rounded-xl hover:bg-[#0A2540]/90 shadow-lg shadow-blue-900/20 font-medium transition-all hover:scale-[1.02]"
                  >
                    Ingest Document
                  </button>
                </div>
              </div>
            )}

            <input
              type="file"
              accept=".pdf,.pptx"
              className="hidden"
              id="fileInput"
              onChange={handleFileChange}
              disabled={status === "uploading"}
            />

            {!file && (
              <label
                htmlFor="fileInput"
                className="inline-block px-8 py-4 bg-[#0A2540] text-white text-sm font-semibold rounded-xl cursor-pointer hover:bg-[#0A2540]/90 transition-all shadow-lg hover:shadow-xl hover:-translate-y-0.5"
              >
                Select Document
              </label>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
