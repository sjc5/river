import { Context, Next } from "hono";
import { renderToString } from "preact-render-to-string";
import { HeadBlock, RouteData } from "../../../common/index.mjs";
import { getRouteData } from "../utils/get-root-data.js";

type JSXElement = any;

async function renderRoot({
  c,
  next,
  defaultHeadBlocks,
  root: Root,
}: {
  c: Context;
  next: Next;
  defaultHeadBlocks: HeadBlock[];
  root: (props: RouteData) => JSXElement;
}) {
  const routeData = await getRouteData({ c, defaultHeadBlocks });

  if (routeData instanceof Response) {
    return routeData;
  }

  if (!routeData) {
    return await next();
  }

  return c.html(
    "<!doctype html>" +
      renderToString(
        // @ts-ignore
        <Root {...routeData} />,
      ),
  );
}

export { renderRoot };
