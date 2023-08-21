import { ThemeToggle } from "@/components/ThemeToggle";
import Link from "next/link";
import { UploadDialog } from "./UploadDialog";

export function Nav() {
  return (
    <nav className="max-h-lg mb-12 flex justify-between border-b border-border px-6 py-3">
      <div className="m-auto flex w-full max-w-screen-lg items-center">
        <div className="flex flex-col">
          <Link href="/" className="text-2xl font-bold text-primary">
            {" "}
            Redroc
          </Link>
        </div>

        <div className="flex h-fit flex-1 items-center justify-end gap-2 text-right">
          <UploadDialog />
          <ThemeToggle />
        </div>
      </div>
    </nav>
  );
}
