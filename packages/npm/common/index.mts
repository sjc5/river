export const HWY_PREFIX = "__hwy_internal__";
export const HWY_PREFIX_JSON = `${HWY_PREFIX}json`;
export const HWY_SYMBOL = Symbol.for(HWY_PREFIX);
export const HWY_ROUTE_CHANGE_EVENT_KEY = "hwy:route-change";

export interface AdHocData extends Record<string, any> {}

type HwyClientGlobal = {
  loadersData: Array<any>;
  importURLs: Array<string>;
  outermostErrorBoundaryIndex: number;
  splatSegments: Array<string>;
  params: Record<string, string>;
  actionData: any;
  activeComponents: Array<any>;
  activeErrorBoundaries: Array<any>;
  adHocData: AdHocData;
  buildID: string;
};

export type HwyClientGlobalKey = keyof HwyClientGlobal;

export function getHwyClientGlobal() {
  const dangerousGlobalThis = globalThis as any;

  function get<K extends HwyClientGlobalKey>(key: K) {
    return dangerousGlobalThis[HWY_SYMBOL][key] as HwyClientGlobal[K];
  }

  function set<K extends HwyClientGlobalKey, V extends HwyClientGlobal[K]>(
    key: K,
    value: V,
  ) {
    dangerousGlobalThis[HWY_SYMBOL][key] = value;
  }

  return { get, set };
}

export type TitleHeadBlock = { title: string };
export type OtherHeadBlock = {
  tag: "meta" | "base" | "link" | "style" | "script" | "noscript" | string;
  attributes: Record<string, string | undefined>;
};
export type HeadBlock = TitleHeadBlock | OtherHeadBlock;

export type GetRouteDataOutput = {
  title: string;
  metaHeadBlocks: Array<OtherHeadBlock>;
  restHeadBlocks: Array<OtherHeadBlock>;
  loadersData: Array<any>;
  importURLs: Array<string>;
  outermostErrorBoundaryIndex: number;
  splatSegments: Array<string>;
  params: Record<string, string>;
  actionData: Array<any>;
  adHocData: Record<string, any>;
  buildID: string;
  deps: Array<string>;

  // SSR Only
  activeErrorBoundaries: Array<any> | null;
  activeComponents: Array<any> | null;
};

export type RouteData = {
  response: Response | null;
  data: GetRouteDataOutput | null;
  mergedResponseInit: ResponseInit | null;
  ssrData?: {
    ssrInnerHTML: string;
    clientEntryURL: string;
    devRefreshScript: string;
    criticalCSSElementID: string;
    criticalCSS: string;
    bundledCSSURL: string;
  };
};
