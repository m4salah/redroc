
import { ThemeToggle } from "@/components/ThemeToggle";
import { UploadDialog } from "./UploadDialog";

export function Nav() {
    return (
        <nav
            className="flex px-6 py-3 border-b border-border justify-between max-h-lg mb-12"
        >
            <div className="max-w-screen-lg flex w-full m-auto items-center">
                <div className="flex flex-col">
                    <a href="/" className="text-2xl font-bold text-primary"> Redroc</a>
                </div>

                <div className="flex flex-1 justify-end text-right h-fit items-center gap-2">
                    <UploadDialog />
                    <ThemeToggle />
                </div>
            </div>
        </nav>

    )
}