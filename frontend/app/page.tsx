import FileUploader from "@/features/documents/components/FileUploader";

export default function Home() {
  return (
    <main className="flex-1 h-screen overflow-auto bg-white font-sans">
      {/* Header */}
      <header className="flex items-center justify-between px-10 py-6 bg-white/80 backdrop-blur-sm sticky top-0 z-10 border-b border-gray-100">
        <div>
          <h1 className="text-xl font-bold text-gray-900 tracking-tight flex items-center gap-2">
            <span className="text-[var(--brand-blue)]">Knowledge Base</span>{" "}
            Manager
          </h1>
          <p className="text-gray-500 text-xs font-medium mt-1">
            Configure and update your AI knowledge sources.
          </p>
        </div>
      </header>

      {/* Content */}
      <div className="w-full flex flex-col items-center px-10 mt-16 text-center animate-in fade-in slide-in-from-bottom-4 duration-700">
        <div className="mb-10">
          <h2 className="text-3xl font-bold text-gray-900 mb-3 tracking-tight">
            Ingest New Knowledge
          </h2>
          <p className="text-gray-600 max-w-lg mx-auto leading-relaxed text-[15px]">
            Upload your technical manuals, sales decks (PPT), or PDF guidelines
            to enhance the BPT Chatbot's intelligence.
          </p>
        </div>

        <FileUploader />
      </div>
    </main>
  );
}
