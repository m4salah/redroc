import { ThemeToggle } from "@/components/ThemeToggle";
import { SearchIcon } from "lucide-react";
import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import { Button } from "~/shadcn/ui/button";
import { Input } from "~/shadcn/ui/input";
import { UploadDialog } from "./UploadDialog";

export function Nav() {
  const [q, setQ] = useState("");

  const router = useRouter();

  useEffect(() => {
    const q = (router.query.q ?? "") as string;
    setQ(q);
  }, [router.query.q]);

  const handleSearch = (q: string) => {
    void router.push({ query: { q }, pathname: "/" });
  };

  return (
    <nav className="max-h-lg mb-12 flex justify-between border-b border-border px-6 py-3">
      <div className="m-auto flex w-full max-w-screen-lg items-center">
        <div className="flex flex-col">
          <Link href="/" className="text-2xl font-bold text-primary">
            <Image
              src={"/favicon.svg"}
              alt="Redroc"
              width={40}
              height={40}
              className="mr-4 -scale-x-100"
            />
          </Link>
        </div>

        <div className="flex h-fit flex-1 items-center justify-end gap-2 text-right">
          <div className="relative flex flex-row items-center gap-1">
            <Input
              id="search"
              name="search"
              type="text"
              placeholder="Search..."
              className="max-w-xs pr-10"
              onChange={(e) => setQ(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  handleSearch(q);
                }
              }}
              value={q}
            />
            <Button
              onClick={() => handleSearch(q)}
              className="absolute right-0 rounded-l-none"
              variant={"outline"}
              size={"icon"}
            >
              <SearchIcon className="h-4 w-4" />
            </Button>
          </div>
          <UploadDialog />
          <ThemeToggle />
        </div>
      </div>
    </nav>
  );
}
