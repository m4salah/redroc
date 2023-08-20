import { Bricolage_Grotesque } from "next/font/google";
import Head from "next/head";
import { PropsWithChildren } from "react";
import { Nav } from "./Nav";
import { ThemeProvider } from "./ThemeProvider";

const font = Bricolage_Grotesque({ weight: ['200', "300", "400", "500", "600"], subsets: ['latin'], variable: "--font-bricolage" });

export default function RootLayout({ children }: PropsWithChildren) {
    return (
        <>
            <Head>
                <title>Redroc</title>
                <meta name="description" content="Welcome to Redroc the place that hold your images" />
                <link rel="icon" href="/favicon.svg" />
            </Head>
            <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
                <style jsx global>{`
                    html {
                        font-family: ${font.style.fontFamily};
                    }
                `}</style>
                <Nav />
                <main className={`flex min-h-screen flex-col items-center justify-center bg-gradient-to-b bg-background`}>
                    {children}
                </main>
            </ThemeProvider>
        </>
    )
}
