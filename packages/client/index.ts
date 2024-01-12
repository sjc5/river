export type { RouteData } from "../common/index.mjs";
export {
  customHistory,
  getIsInternalLink,
  getShouldPreventLinkDefault,
  initPreactClient,
  navigate,
  submit,
} from "./src/init-preact-client.js";
export { RootOutlet, getAdHocDataSignal } from "./src/recursive.js";
