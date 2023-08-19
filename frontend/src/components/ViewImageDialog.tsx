import { Button } from "@/shadcn/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/shadcn/ui/dialog";
import { DownloadIcon } from "@radix-ui/react-icons";

interface Props {
    item: string;
}

function ViewImageDialog({ item }: Props) {
    return (
        <Dialog>
            <DialogTrigger asChild>
                <img
                    src={new URL(item, "https://api.redroc.xyz").toString()}
                    className="w-full h-full object-cover rounded hover:scale-110 transition-all duration-200 ease-in-out"
                    loading="lazy"
                />
            </DialogTrigger>
            <DialogContent className="sm:max-w-screen-xl">
                <DialogHeader>
                    <DialogTitle>
                        <Button onClick={() => {
                            window.open(new URL(item.replace("thumbnail_", ""), "https://api.redroc.xyz").toString(), "_blank")
                        }}>
                            <DownloadIcon />
                        </Button>
                    </DialogTitle>
                </DialogHeader>
                <img
                    src={new URL(item.replace("thumbnail_", ""), "https://api.redroc.xyz").toString()}
                    className="w-full h-full rounded"
                />
            </DialogContent>
        </Dialog>
    )
}

export default ViewImageDialog;
