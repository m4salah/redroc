/* eslint-disable @typescript-eslint/no-misused-promises */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
import { Button } from "@/shadcn/ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/shadcn/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/shadcn/ui/form";
import { Input } from "@/shadcn/ui/input";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/router";
import { useState } from "react";
import { useForm } from "react-hook-form";
import * as z from "zod";
import { redrocClient } from "~/apiClient/redrocClient";
import { refreshData } from "~/lib/utils";

const MAX_FILE_SIZE = 4 * 1024 * 1024; // 4 MB

const formSchema = z.object({
  username: z.string().min(2, {
    message: "Username must be at least 2 characters.",
  }),
  file: z
    .any()
    .refine((file) => file !== null, {
      message: "File is required.",
    })
    .refine((file: File) => file?.size <= MAX_FILE_SIZE, {
      message: `Max image size is 4MB.`,
    }),
  hashtags: z.string().min(2, {
    message: "Hashtags must be at least 2 characters.",
  }),
});

export function UploadDialog() {
  const [isLoading, setIsLoading] = useState(false);

  const [open, setOpen] = useState(false);

  const router = useRouter();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
      file: null,
      hashtags: "",
    },
  });

  function onSubmit(values: z.infer<typeof formSchema>) {
    setIsLoading(true);

    const formData = new FormData();

    const hashtags = JSON.stringify(
      values.hashtags.replaceAll(" ", "").split(",")
    );

    formData.append("username", values.username);
    // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
    formData.append("file", values.file);
    formData.append("hashtags", hashtags);
    console.log(values);

    redrocClient
      .post("upload", formData, {
        headers: {
          "Content-Type": "multipart/form-data; charset=utf-8",
        },
      })
      .then(({ data }) => {
        console.log(data);
        setIsLoading(false);
        refreshData(router);
        setOpen(false);
      })
      .catch((err) => {
        console.log(err.response);
        setIsLoading(false);
        refreshData(router);
        setOpen(false);
      });
  }

  const onOpenDialog = (open: boolean) => {
    setOpen(open);
    if (!open) form.reset();
  };

  return (
    <Dialog open={open} onOpenChange={onOpenDialog}>
      <DialogTrigger asChild>
        <Button variant={"secondary"}>Upload</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Upload Image</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="username"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Username</FormLabel>
                  <FormControl>
                    <Input placeholder="mohamed" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="file"
              render={({ field: { value, onChange, ...field } }) => (
                <FormItem>
                  <FormLabel>File</FormLabel>
                  <FormControl>
                    <Input
                      {...field}
                      value={value?.fileName}
                      type="file"
                      onChange={(e) => {
                        if (e.target.files?.length) {
                          onChange(e.target.files[0]);
                          return;
                        }
                        onChange(null);
                      }}
                      className="cursor-pointer"
                      accept="image/*"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="hashtags"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Hashtags</FormLabel>
                  <FormControl>
                    <Input {...field} placeholder="hashtag1, hashtag2, ...." />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? "Uploading..." : "Upload"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
