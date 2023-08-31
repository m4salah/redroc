import axios from "axios";
import { env } from "~/env.mjs";

export const redrocClient = axios.create({ baseURL: env.NEXT_PUBLIC_BACKEND_URL })