import axios from "axios";
import {
  type GetServerSideProps,
  type InferGetServerSidePropsType,
} from "next";
import { ViewImageDialog } from "~/components/ViewImageDialog";

type Repo = {
  data: string[];
};

export default function Home({
  repo,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
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
}> = async () => {
  const { data } = await axios.get<string[]>("https://api.redroc.xyz/search");
  return { props: { repo: { data } } };
};
