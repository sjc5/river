/**********************************************************************
/ Generated by tsgen. DO NOT EDIT.
/*********************************************************************/

/**********************************************************************
/ Collection:
/*********************************************************************/

const routes = [
	{
		_type: "loader",
		isRootData: true,
		pattern: "",
		phantomOutputType: null as unknown as RootData,
	},
	{
		_type: "loader",
		pattern: "/",
		phantomOutputType: null as unknown as string,
	},
	{
		_type: "loader",
		pattern: "/start",
		phantomOutputType: undefined,
	},
] as const;

/**********************************************************************
/ Ad Hoc Types:
/*********************************************************************/

export type RootData = {
	SiteTitle: string;
	LatestVersion: string;
};

/**********************************************************************
/ Extra TS Code:
/*********************************************************************/

export type RiverLoader = Extract<(typeof routes)[number], { _type: "loader" }>;
export type RiverLoaders = { [K in RiverLoaderPattern]: Extract<RiverLoader, { pattern: K }>; };
export type RiverLoaderPattern = RiverLoader["pattern"];
export type RiverLoaderOutput<T extends RiverLoaderPattern> = Extract<RiverLoader, { pattern: T }>["phantomOutputType"];
export type RiverRootData = Extract<(typeof routes)[number], { isRootData: true }>["phantomOutputType"];

export const ACTIONS_ROUTER_MOUNT_ROOT = "/river-api/";

/**********************************************************************
/ River Vite Plugin:
/*********************************************************************/

import type { Plugin } from "vite";

const rollupOptions = {
	input: [
		"./ts/entry.tsx",
		"ts/home.tsx",
		"ts/root.tsx",
		"ts/start.tsx",
	] as string[],
	preserveEntrySignatures: "exports-only",
	output: {
		assetFileNames: "river_out_[name]-[hash][extname]",
		chunkFileNames: "river_out_[name]-[hash].js",
		entryFileNames: "river_out_[name]-[hash].js",
	},
} as const;

export const staticPublicAssetMap = {
	"desktop.svg": "desktop_eebc981612eb.svg",
	"favicon.svg": "favicon_ed2aaf004a0d.svg",
	"full-logo.svg": "full-logo_e0ea7a3d3cf2.svg",
	"logo.svg": "logo_d9b0e1618362.svg",
	"moon.svg": "moon_7e0c08985ebe.svg",
	"river-banner.webp": "river-banner_6dfc0fa16366.webp",
	"sun.svg": "sun_338b26f6045d.svg"
} as const;

export type StaticPublicAsset = keyof typeof staticPublicAssetMap;

declare global {
	function hashedURL(staticPublicAsset: StaticPublicAsset): string;
}

export function riverVitePlugin(): Plugin {
	return {
		name: "river-vite-plugin",
		config(c, { command }) {
			const mp = c.build?.modulePreload;
			const roi = c.build?.rollupOptions?.input;
			const ign = c.server?.watch?.ignored;
			const dedupe = c.resolve?.dedupe;

			const isDev = command === "serve";

			return {
				...c,
				base: isDev ? "/" : "/public/",
				build: {
					target: "es2022",
					...c.build,
					modulePreload: { 
						polyfill: false,
						...(typeof mp === "object" ? mp : {}),
					},
					rollupOptions: {
						...c.build?.rollupOptions,
						...rollupOptions,
						input: [
							...rollupOptions.input,
							...(Array.isArray(roi) ? roi : []),
						],
					},
				},
				server: {
					...c.server,
					headers: {
						...c.server?.headers,
						// ensure versions of dynamic imports without the latest
						// hmr updates are not cached by the browser during dev
						"cache-control": "no-store",
					},
					watch: {
						...c.server?.watch,
						ignored: [
							...(Array.isArray(ign) ? ign : []),
							...[
								"**/*.go",
								"**/static/private",
								"**/go/app/kiruna.config.json",
								"**/ts/river.gen.ts",
								"**/ts/routes.ts"
							],
						],
					},
				},
				resolve: {
					...c.resolve,
					dedupe: [
						...(Array.isArray(dedupe) ? dedupe : []),
						...["solid-js","solid-js/web"]
					],
				},
			};
		},
		transform(code, id) {
			const isNodeModules = /node_modules/.test(id);
			if (isNodeModules) return null;
			const assetRegex = /hashedURL\s*\(\s*(["'`])(.*?)\1\s*\)/g;
			const needsReplacement = assetRegex.test(code);
			if (!needsReplacement) return null;
			const replacedCode = code.replace(
				assetRegex,
				(original, _, assetPath) => {
					const hashed = (staticPublicAssetMap as Record<string, string>)[assetPath];
					if (!hashed) return '\"' + assetPath + '\"';
					return `"/public/${hashed}"`;
				},
			);
			if (replacedCode === code) return null;
			return replacedCode;
		},
	};
}
