"use client";

import { useState, useRef, useEffect } from "react";
import { Send, Bot, User, Sparkles, StopCircle } from "lucide-react";
import ReactMarkdown from "react-markdown";
import { sendChatMessage } from "../api/chatService";
import { Message } from "../types";

export default function ChatInterface() {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: "1",
      role: "bot",
      content: "Hello. I;m ready to help you analyze your documents.",
    },
  ]);
  const [inputValue, setInputValue] = useState("");
  const [isThinking, setIsThinking] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, isThinking]);

  const handleSendMessage = async (e?: React.FormEvent) => {
    e?.preventDefault();
    if (!inputValue.trim()) return;

    const userMsg: Message = {
      id: Date.now().toString(),
      role: "user",
      content: inputValue,
    };
    setMessages((prev) => [...prev, userMsg]);
    setInputValue("");
    setIsThinking(true);

    try {
      const data = await sendChatMessage(userMsg.content);
      const botMsg: Message = {
        id: (Date.now() + 1).toString(),
        role: "bot",
        content: data.response,
        sources: data.sources || [],
      };
      setMessages((prev) => [...prev, botMsg]);
    } catch {
      setMessages((prev) => [
        ...prev,
        {
          id: (Date.now() + 1).toString(),
          role: "bot",
          content:
            "I'm having trouble connecting to the knowledge base. Please try again.",
        },
      ]);
    } finally {
      setIsThinking(false);
    }
  };

  return (
    <div className="flex flex-col h-screen bg-white font-sans">
      {/* Minimal Header */}
      <header className="h-16 flex items-center justify-between px-8 border-b border-gray-100 bg-white/80 backdrop-blur-sm sticky top-0 z-10 transition-all">
        <div className="flex items-center gap-3">
          <div className="flex items-center justify-center w-10 h-10 rounded-full bg-[#0A2540] text-white shadow-sm">
            <Bot size={20} />
          </div>
          <div className="flex flex-col justify-center">
            <h1 className="text-sm font-bold text-gray-900 leading-tight">
              BPT AI Assistant
            </h1>
            <div className="flex items-center gap-1.5 mt-0.5">
              <div className="w-1.5 h-1.5 rounded-full bg-green-500"></div>
              <span className="text-xs font-medium text-gray-500">
                Online • Powered by Couchbase
              </span>
            </div>
          </div>
        </div>
      </header>

      {/* Chat Area */}
      <div className="flex-1 overflow-y-auto p-8 scroll-smooth bg-gray-50/30">
        <div className="max-w-full mx-auto space-y-8 px-4">
          {messages.map((msg) => (
            <div key={msg.id}>
              <div
                className={`flex gap-4 ${msg.role === "user" ? "flex-row-reverse" : "flex-row"}`}
              >
                {/* Avatar */}
                <div
                  className={`w-8 h-8 rounded-xl flex items-center justify-center shrink-0 shadow-sm mt-1 transition-all
                    ${
                      msg.role === "user"
                        ? "bg-white border border-[var(--border-secondary)] text-[var(--text-secondary)]"
                        : "bg-[var(--brand-blue)] text-white shadow-md shadow-blue-900/20"
                    }
                  `}
                >
                  {msg.role === "user" ? <User size={16} /> : <Bot size={16} />}
                </div>

                {/* Message Content */}
                <div
                  className={`flex flex-col max-w-[70%] ${msg.role === "user" ? "items-end" : "items-start"}`}
                >
                  <div
                    className={`px-6 py-4 text-[15px] leading-relaxed shadow-sm transition-all
                      ${
                        msg.role === "user"
                          ? "bg-[#0A2540] text-white rounded-2xl"
                          : "bg-white text-gray-800 rounded-2xl border border-gray-200"
                      }
                    `}
                  >
                    {msg.role === "bot" ? (
                      <div className="prose prose-sm prose-slate max-w-none">
                        <ReactMarkdown>{msg.content || ""}</ReactMarkdown>
                      </div>
                    ) : (
                      <p className="whitespace-pre-wrap">{msg.content}</p>
                    )}
                  </div>

                  {/* Sources Citation */}
                  {msg.role === "bot" &&
                    msg.sources &&
                    msg.sources.length > 0 && (
                      <div className="mt-2 w-full pt-2 border-t border-gray-100/50">
                        <div className="flex flex-wrap gap-x-3 gap-y-1">
                          <span className="text-[10px] font-medium text-gray-400 uppercase tracking-wide self-center">
                            Sources:
                          </span>
                          {msg.sources.map((source, idx) => (
                            <div
                              key={idx}
                              className="flex items-center gap-1.5 px-2 py-1 bg-gray-50/50 rounded hover:bg-gray-100 transition-colors cursor-pointer group"
                            >
                              <span className="text-[11px] text-gray-500 group-hover:text-blue-600 transition-colors truncate max-w-[150px]">
                                {source.filename}
                              </span>
                              {source.page > 0 && (
                                <span className="text-[10px] text-gray-400 group-hover:text-blue-400">
                                  (p. {source.page})
                                </span>
                              )}
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                </div>
              </div>
            </div>
          ))}

          {/* Thinking State */}
          {isThinking && (
            <div className="flex gap-4 animate-in fade-in slide-in-from-bottom-2">
              <div className="w-8 h-8 rounded-xl bg-[var(--brand-blue)] text-white flex items-center justify-center shrink-0 mt-1 shadow-md shadow-blue-900/20">
                <Bot size={16} />
              </div>
              <div className="px-5 py-4 bg-white border border-[var(--border-secondary)] rounded-2xl rounded-tl-sm shadow-sm">
                <div className="flex gap-1.5 items-center h-full">
                  <div className="w-1.5 h-1.5 bg-[var(--brand-accent)] rounded-full animate-bounce [animation-delay:-0.3s]"></div>
                  <div className="w-1.5 h-1.5 bg-[var(--brand-accent)] rounded-full animate-bounce [animation-delay:-0.15s]"></div>
                  <div className="w-1.5 h-1.5 bg-[var(--brand-accent)] rounded-full animate-bounce"></div>
                </div>
              </div>
            </div>
          )}
          <div ref={messagesEndRef} />
        </div>
      </div>

      {/* Input Area */}
      <div className="p-6 bg-white border-t border-gray-100 shrink-0 z-20">
        <div className="max-w-full mx-auto px-4">
          <form
            onSubmit={handleSendMessage}
            className="relative flex items-center gap-2 p-2 bg-gray-50 border border-gray-200 rounded-2xl shadow-sm focus-within:ring-2 focus-within:ring-blue-500/20 focus-within:border-blue-500/50 transition-all"
          >
            <input
              type="text"
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              placeholder="Ask anything..."
              className="flex-1 px-4 py-2 bg-transparent text-gray-900 placeholder:text-gray-400 focus:outline-none text-[15px]"
              disabled={isThinking}
            />
            <button
              type="submit"
              disabled={!inputValue.trim() || isThinking}
              className="p-2.5 bg-[#0A2540] text-white rounded-xl hover:bg-[#0A2540]/90 disabled:bg-gray-100 disabled:text-gray-300 transition-all shadow-md shadow-blue-900/10 hover:shadow-lg disabled:shadow-none"
            >
              {isThinking ? <StopCircle size={18} /> : <Send size={18} />}
            </button>
          </form>
          <p className="text-center text-[10px] text-[var(--text-tertiary)] mt-3 font-medium">
            Powered by BPT AI. Information may be inaccurate.
          </p>
        </div>
      </div>
    </div>
  );
}
