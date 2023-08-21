import { Button } from "@/shadcn/ui/button";
import { Dialog, DialogContent, DialogTrigger } from "@/shadcn/ui/dialog";
import { DownloadIcon } from "lucide-react";
import Image from "next/image";
import { useState } from "react";

interface Props {
  item: string;
}
const shimmer = (w: number, h: number) => `
<svg width="${w}" height="${h}" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
  <defs>
    <linearGradient id="g">
      <stop stop-color="#333" offset="20%" />
      <stop stop-color="#222" offset="50%" />
      <stop stop-color="#333" offset="70%" />
    </linearGradient>
  </defs>
  <rect width="${w}" height="${h}" fill="#333" />
  <rect id="r" width="${w}" height="${h}" fill="url(#g)" />
  <animate xlink:href="#r" attributeName="x" from="-${w}" to="${w}" dur="1s" repeatCount="indefinite"  />
</svg>`;

const toBase64 = (str: string) =>
  typeof window === "undefined"
    ? Buffer.from(str).toString("base64")
    : window.btoa(str);

export function ViewImageDialog({ item }: Props) {
  const [open, setOpen] = useState(false);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Image
          placeholder={`data:image/svg+xml;base64,${toBase64(
            shimmer(700, 475)
          )}`}
          src={new URL(item, "https://api.redroc.xyz").toString()}
          className="h-full w-full rounded object-cover transition-all duration-200 ease-in-out hover:scale-110"
          loading="lazy"
          height={192}
          width={192}
          alt={item}
        />
      </DialogTrigger>
      <DialogContent className="h-5/6 max-w-screen-2xl">
        <Image
          quality={100}
          placeholder={`data:image/svg+xml;base64,${toBase64(
            shimmer(700, 475)
          )}`}
          src={new URL(
            item.replace("thumbnail_", ""),
            "https://api.redroc.xyz"
          ).toString()}
          className="h-full w-full rounded object-contain"
          alt={item}
          priority
          fill
        />
        <Button
          className="absolute right-20 top-4"
          variant={"default"}
          onClick={() => {
            window.open(
              new URL(
                item.replace("thumbnail_", ""),
                "https://api.redroc.xyz"
              ).toString(),
              "_blank"
            );
          }}
        >
          <DownloadIcon className="h-4 w-4" />
        </Button>
      </DialogContent>
    </Dialog>
  );
}
