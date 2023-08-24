import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import { env } from "~/env.mjs"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * @param thumbnailName string
 * @returns url string to the thumbnail
 */
export function getThumbnailURL(thumbnailName: string): string {
  return new URL(thumbnailName, env.NEXT_PUBLIC_BACKEND_URL).toString()
}

/**
 * @param thumbnailName string
 * @returns url string to the full size image
 */
export function getImageURL(thumbnailName: string): string {
  return new URL(
    thumbnailName.replace("thumbnail_", ""),
    env.NEXT_PUBLIC_BACKEND_URL
  ).toString()
}