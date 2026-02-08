"use client";

import {
  UploadCloud,
  MessageSquare,
  Layers,
  LogOut,
  ChevronRight,
  Database,
} from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import Image from "next/image";

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <div className="w-[280px] h-screen bg-[#0A2540] text-white flex flex-col shrink-0 font-sans shadow-xl z-50">
      {/* Brand */}
      <div className="h-24 flex items-center px-6 mb-2">
        <div className="flex items-center gap-3">
          <div className="relative w-10 h-10 shrink-0">
            <Image
              src="/logos/bpt_logos/bpt2.png"
              alt="BPT Logo"
              fill
              className="object-contain"
            />
          </div>
          <div className="flex flex-col justify-center">
            <h1 className="font-semibold text-lg leading-tight tracking-tight text-white">
              BluePower
            </h1>
            <span className="text-sm text-blue-200/80 font-light tracking-wide">
              Knowledge
            </span>
          </div>
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 px-4 space-y-1">
        <div className="mb-4 px-2">
          <p className="text-xs font-semibold text-blue-400/80 uppercase tracking-wider">
            Platform
          </p>
        </div>

        <NavLink
          href="/chat"
          icon={<MessageSquare size={18} />}
          label="AI Assistant"
          active={pathname === "/chat"}
        />
        <NavLink
          href="/documents"
          icon={<Layers size={18} />}
          label="Documents"
          active={pathname === "/documents"}
        />
        <NavLink
          href="/"
          icon={<UploadCloud size={18} />}
          label="Import Data"
          active={pathname === "/"}
        />
      </nav>

      {/* User Section */}
      <div className="p-4 border-t border-white/10">
        <div className="flex items-center justify-between p-3 rounded-xl bg-white/5 border border-white/10 hover:bg-white/10 transition-colors cursor-pointer group">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-full bg-blue-500/30 flex items-center justify-center text-xs font-medium text-white border border-white/10">
              JP
            </div>
            <div className="flex flex-col">
              <span className="text-sm font-medium text-white">Jessa P.</span>
              <span className="text-[10px] text-blue-300/80">
                Admin Workspace
              </span>
            </div>
          </div>
          <LogOut
            size={14}
            className="text-blue-300/60 group-hover:text-white transition-colors"
          />
        </div>
      </div>
    </div>
  );
}

function NavLink({
  href,
  icon,
  label,
  active,
}: {
  href: string;
  icon: React.ReactNode;
  label: string;
  active: boolean;
}) {
  const isActive = active || (href !== "/" && active);

  return (
    <Link href={href}>
      <div
        className={`group flex items-center justify-between px-3 py-2.5 rounded-lg transition-all duration-200 ${
          active
            ? "bg-white/10 text-white shadow-sm ring-1 ring-white/10"
            : "text-blue-100/70 hover:bg-white/5 hover:text-white"
        }`}
      >
        <div className="flex items-center gap-3">
          <span
            className={
              active ? "text-white" : "text-blue-300/70 group-hover:text-white"
            }
          >
            {icon}
          </span>
          <span className="text-sm font-medium">{label}</span>
        </div>
        {active && <ChevronRight size={14} className="text-white/50" />}
      </div>
    </Link>
  );
}
