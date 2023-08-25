import {
  type GetServerSideProps,
  type InferGetServerSidePropsType,
} from "next";
import { useEffect } from "react";
import { redrocClient } from "~/apiClient/redrocClient";
import { ViewImageDialog } from "~/components/ViewImageDialog";
import { getSocketURL } from "~/lib/utils";

type Repo = {
  data: string[];
};

export default function Home({
  repo,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  useEffect(() => {
    const io = new WebSocket(getSocketURL());
    io.onopen = () => {
      console.log("connected");
      io.send("hello");
    };
    io.onclose = () => {
      console.log("closed");
    };
    io.onmessage = (e) => {
      console.log("message", e.data);
    };
  }, []);

  return !repo.data ? (
    <div className="flex h-full w-full items-center">
      <p className="text text-center text-xl">
        No images with that keyword, Try different keyword !!!
      </p>
    </div>
  ) : (
    <ul role="list" className="flex flex-wrap justify-center gap-4">
      {repo.data.map((item) => (
        <li className="h-48 w-48 cursor-pointer rounded" key={item}>
          <ViewImageDialog item={item} />
        </li>
      ))}
    </ul>
  );
}

export const getServerSideProps: GetServerSideProps<{
  repo: Repo;
}> = async (context) => {
  const q = context.query.q ?? "";
  const { data } = await redrocClient.get<string[]>("search", {
    params: { q },
  });
  return { props: { repo: { data } } };
};
