import { addBuildIDListener, getRootEl, initClient } from "river.now/client";
import { render } from "solid-js/web";
import { App } from "./home.tsx";

await initClient(() => {
	render(() => <App />, getRootEl());
});

addBuildIDListener((e) => {
	window.location.reload();
});
