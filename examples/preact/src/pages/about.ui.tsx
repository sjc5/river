import { UIProps } from "@hwy-js/react";
import Clicker from "./_index.ui.js";

export default function ({ Outlet, loaderData }: UIProps) {
  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "1.5rem" }}>
      <Clicker />
      Random number from loader: {loaderData}
      <p>
        This is the about page. This is coming from{" "}
        <code>src/pages/about.ui.tsx</code> This is functioning as a "layout
        route" because it is rendering an <code>{`<Outlet />`}</code> (see
        below).
      </p>
      <div style={{ textTransform: "uppercase", fontWeight: "bold" }}>
        Outlet:
      </div>
      <div className="outlet-wrapper">
        <Outlet />
      </div>
    </div>
  );
}
