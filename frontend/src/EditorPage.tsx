import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";

export function EditorPage() {
	const location = useLocation();
	const [str, setStr] = useState("");

	useEffect(() => {
		fetch(
			`${import.meta.env.VITE_SERVER}/files${location.pathname.slice(
				"/edit".length
			)}`
		)
			.then((res) => {
				if (!res.ok) throw new Error(res.statusText);
				if (location.pathname.endsWith(".md")) return res.text();
				return res.text();
			})
			.then((body) => {
				console.log(body);
			});
	}, [location]);

	return <div>Editor goes here</div>;
}
